package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/brotherlogic/goserver"
	"github.com/brotherlogic/goserver/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"

	pbgh "github.com/brotherlogic/githubcard/proto"
	pbg "github.com/brotherlogic/goserver/proto"
	rapb "github.com/brotherlogic/recordadder/proto"
	rbpb "github.com/brotherlogic/recordbudget/proto"
	pbrc "github.com/brotherlogic/recordcollection/proto"
	rcpb "github.com/brotherlogic/recordcollection/proto"
	pb "github.com/brotherlogic/recordwants/proto"
)

func init() {
	resolver.Register(&utils.DiscoveryServerResolverBuilder{})
}

var (
	//wants - the print queue
	wants = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "recordwants_total_wants",
		Help: "The size of the print queue",
	})
)

const (
	// KEY - where the wants are stored
	KEY = "/github.com/brotherlogic/recordwants/config"
)

type alerter interface {
	alert(ctx context.Context, want *pb.MasterWant, c, total int)
}

type prodAlerter struct {
	dial func(ctx context.Context, server string) (*grpc.ClientConn, error)
}

func (p *prodAlerter) alert(ctx context.Context, want *pb.MasterWant, c, total int) {
	conn, err := p.dial(ctx, "githubcard")
	if err == nil {
		defer conn.Close()
		client := pbgh.NewGithubClient(conn)
		client.AddIssue(ctx, &pbgh.Issue{Service: "recordwants", Title: fmt.Sprintf("Want Processing Needed!"), Body: fmt.Sprintf("%v/%v - %v", c, total, want)}, grpc.FailFast(false))
	}
}

type recordAdder interface {
	getAdds(ctx context.Context) ([]int32, error)
}

type prodRecordAdder struct {
	dial func(ctx context.Context, server string) (*grpc.ClientConn, error)
}

func (p *prodRecordAdder) getAdds(ctx context.Context) ([]int32, error) {
	conn, err := p.dial(ctx, "recordadder")
	if err != nil {
		return []int32{}, err
	}
	defer conn.Close()

	client := rapb.NewAddRecordServiceClient(conn)
	resp, err := client.ListQueue(ctx, &rapb.ListQueueRequest{})
	if err != nil {
		return []int32{}, err
	}

	nums := []int32{}
	for _, add := range resp.GetRequests() {
		nums = append(nums, add.GetId())
	}
	return nums, err
}

type recordGetter interface {
	getRecordsSince(ctx context.Context, since int64) ([]int32, error)
	getRecords(ctx context.Context, id int32) ([]int32, error)
	getRecord(ctx context.Context, id int32) (*pbrc.Record, error)
	getWants(ctx context.Context) ([]*pbrc.Want, error)
	unwant(ctx context.Context, want *pb.MasterWant) error
	want(ctx context.Context, want *pb.MasterWant) error
}

type prodGetter struct {
	dial func(ctx context.Context, server string) (*grpc.ClientConn, error)
}

func (p *prodGetter) getRecordsSince(ctx context.Context, since int64) ([]int32, error) {
	conn, err := p.dial(ctx, "recordcollection")
	if err != nil {
		return []int32{}, err
	}
	defer conn.Close()

	client := pbrc.NewRecordCollectionServiceClient(conn)
	resp, err := client.QueryRecords(ctx, &pbrc.QueryRecordsRequest{Query: &pbrc.QueryRecordsRequest_UpdateTime{since}})

	if err != nil {
		return []int32{}, err
	}

	return resp.GetInstanceIds(), err
}
func (p *prodGetter) getRecords(ctx context.Context, id int32) ([]int32, error) {
	conn, err := p.dial(ctx, "recordcollection")
	if err != nil {
		return []int32{}, err
	}
	defer conn.Close()

	client := pbrc.NewRecordCollectionServiceClient(conn)
	resp, err := client.QueryRecords(ctx, &pbrc.QueryRecordsRequest{Query: &pbrc.QueryRecordsRequest_ReleaseId{id}})

	if err != nil {
		return []int32{}, err
	}

	return resp.GetInstanceIds(), err
}

func (p *prodGetter) getRecord(ctx context.Context, instanceID int32) (*pbrc.Record, error) {
	conn, err := p.dial(ctx, "recordcollection")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pbrc.NewRecordCollectionServiceClient(conn)
	resp, err := client.GetRecord(ctx, &pbrc.GetRecordRequest{InstanceId: instanceID})

	if err != nil {
		return nil, err
	}

	return resp.GetRecord(), err
}

func (p *prodGetter) getWants(ctx context.Context) ([]*pbrc.Want, error) {
	conn, err := p.dial(ctx, "recordcollection")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pbrc.NewRecordCollectionServiceClient(conn)
	resp, err := client.GetWants(ctx, &pbrc.GetWantsRequest{})
	if err != nil {
		return nil, err
	}

	return resp.GetWants(), nil

}

func (p *prodGetter) unwant(ctx context.Context, want *pb.MasterWant) error {
	conn, err := p.dial(ctx, "recordcollection")
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pbrc.NewRecordCollectionServiceClient(conn)
	_, err = client.UpdateWant(ctx, &pbrc.UpdateWantRequest{Update: &pbrc.Want{Release: want.GetRelease()}, Remove: true})
	return err
}

func (p *prodGetter) want(ctx context.Context, want *pb.MasterWant) error {
	conn, err := p.dial(ctx, "recordcollection")
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pbrc.NewRecordCollectionServiceClient(conn)
	_, err = client.UpdateWant(ctx, &pbrc.UpdateWantRequest{Update: &pbrc.Want{Release: want.GetRelease(), EnableWant: true}})
	return err
}

//Server main server type
type Server struct {
	*goserver.GoServer
	recordGetter recordGetter
	alerter      alerter
	lastRun      time.Time
	lastProc     int32
	lastPull     int32
	pull         string
	mmonth       int32
	lastUnwant   string
	budgetPull   time.Duration
	recordAdder  recordAdder
	testing      bool
}

// Init builds the server
func Init() *Server {
	s := &Server{
		&goserver.GoServer{},
		&prodGetter{},
		&prodAlerter{},
		time.Now(),
		-1,
		-1,
		"",
		0,
		"",
		0,
		&prodRecordAdder{},
		true,
	}
	s.recordGetter = &prodGetter{dial: s.FDialServer}
	s.alerter = &prodAlerter{dial: s.FDialServer}
	s.recordAdder = &prodRecordAdder{dial: s.FDialServer}
	return s
}

func (s *Server) save(ctx context.Context, config *pb.Config) error {
	return s.KSclient.Save(ctx, KEY, config)
}

func (s *Server) load(ctx context.Context) (*pb.Config, error) {
	config := &pb.Config{}
	data, _, err := s.KSclient.Read(ctx, KEY, config)

	if err != nil {
		return nil, err
	}

	//Clean out the wants here
	wmap := make(map[int32]*pb.MasterWant)
	config = data.(*pb.Config)

	if config.Spends == nil {
		config.Spends = make(map[int32]*pb.RecordSpend)
	}

	for _, want := range config.Wants {
		if val, ok := wmap[want.GetRelease().GetId()]; ok {
			val.Superwant = val.Superwant || want.Superwant
		} else {
			wmap[want.GetRelease().GetId()] = want
		}
	}

	config = data.(*pb.Config)
	config.Wants = []*pb.MasterWant{}
	for _, want := range wmap {
		config.Wants = append(config.Wants, want)
	}

	return config, err
}

// DoRegister does RPC registration
func (s *Server) DoRegister(server *grpc.Server) {
	pb.RegisterWantServiceServer(server, s)
	rcpb.RegisterClientUpdateServiceServer(server, s)
}

// ReportHealth alerts if we're not healthy
func (s *Server) ReportHealth() bool {
	return true
}

// Shutdown close the server
func (s *Server) Shutdown(ctx context.Context) error {
	return nil
}

// Mote promotes/demotes this server
func (s *Server) Mote(ctx context.Context, master bool) error {
	return nil
}

func (s *Server) getBudget(ctx context.Context) (int32, error) {
	if s.testing {
		return 0, nil
	}
	conn, err := s.FDialServer(ctx, "recordbudget")
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	client := rbpb.NewRecordBudgetServiceClient(conn)
	budg, err := client.GetBudget(ctx, &rbpb.GetBudgetRequest{Year: int32(time.Now().Year())})

	if err != nil {
		return 0, err
	}

	return budg.GetBudget() - budg.GetSpends() - budg.GetPreSpends(), nil
}

// GetState gets the state of the server
func (s *Server) GetState() []*pbg.State {

	return []*pbg.State{}
}

func main() {
	var quiet = flag.Bool("quiet", false, "Show all output")
	var clear = flag.Bool("clear_suuper", false, "Clears all super wants")
	flag.Parse()

	//Turn off logging
	if *quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	server := Init()
	server.PrepServer()
	server.Register = server

	err := server.RegisterServerV2("recordwants", false, true)
	if err != nil {
		return
	}

	if *clear {
		ctx, cancel := utils.BuildContext("recordwants", "recordwants")
		defer cancel()

		config, err := server.load(ctx)
		if err != nil {
			log.Fatalf("Unable to read wants: %v", err)
		}

		for _, w := range config.Wants {
			w.Superwant = false
		}

		server.save(ctx, config)
		log.Fatalf("Saved out cleared wants")
	}

	server.Serve()
}
