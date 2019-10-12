package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	pbgh "github.com/brotherlogic/githubcard/proto"
	"github.com/brotherlogic/goserver"
	pbg "github.com/brotherlogic/goserver/proto"
	"github.com/brotherlogic/goserver/utils"
	pbrc "github.com/brotherlogic/recordcollection/proto"
	pb "github.com/brotherlogic/recordwants/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	// KEY - where the wants are stored
	KEY = "/github.com/brotherlogic/recordwants/config"
)

type alerter interface {
	alert(ctx context.Context, want *pb.MasterWant, c, total int)
}

type prodAlerter struct {
	dial func(server string) (*grpc.ClientConn, error)
}

func (p *prodAlerter) alert(ctx context.Context, want *pb.MasterWant, c, total int) {
	conn, err := p.dial("githubcard")
	if err == nil {
		defer conn.Close()
		client := pbgh.NewGithubClient(conn)
		client.AddIssue(ctx, &pbgh.Issue{Service: "recordwants", Title: fmt.Sprintf("Want Processing Needed!"), Body: fmt.Sprintf("%v/%v - %v", c, total, want)}, grpc.FailFast(false))
	}
}

type recordGetter interface {
	getRecordsSince(ctx context.Context, since int64) ([]int32, error)
	getRecord(ctx context.Context, id int32) (*pbrc.Record, error)
	getWants(ctx context.Context) ([]*pbrc.Want, error)
	unwant(ctx context.Context, want *pb.MasterWant) error
	want(ctx context.Context, want *pb.MasterWant) error
}

type prodGetter struct {
	dial func(server string) (*grpc.ClientConn, error)
}

func (p *prodGetter) getRecordsSince(ctx context.Context, since int64) ([]int32, error) {
	conn, err := p.dial("recordcollection")
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
func (p *prodGetter) getRecord(ctx context.Context, instanceID int32) (*pbrc.Record, error) {
	conn, err := p.dial("recordcollection")
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
	conn, err := p.dial("recordcollection")
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
	conn, err := p.dial("recordcollection")
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pbrc.NewRecordCollectionServiceClient(conn)
	_, err = client.UpdateWant(ctx, &pbrc.UpdateWantRequest{Update: &pbrc.Want{Release: want.GetRelease()}, Remove: true})
	return err
}

func (p *prodGetter) want(ctx context.Context, want *pb.MasterWant) error {
	conn, err := p.dial("recordcollection")
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pbrc.NewRecordCollectionServiceClient(conn)
	_, err = client.UpdateWant(ctx, &pbrc.UpdateWantRequest{Update: &pbrc.Want{Release: want.GetRelease()}})
	return err
}

//Server main server type
type Server struct {
	*goserver.GoServer
	recordGetter recordGetter
	config       *pb.Config
	alerter      alerter
	lastRun      time.Time
	lastProc     int32
	lastPull     int32
	pull         string
	mmonth       int32
	lastUnwant   string
	budgetPull   time.Duration
}

// Init builds the server
func Init() *Server {
	s := &Server{
		&goserver.GoServer{},
		&prodGetter{},
		&pb.Config{},
		&prodAlerter{},
		time.Now(),
		-1,
		-1,
		"",
		0,
		"",
		0,
	}
	s.recordGetter = &prodGetter{dial: s.DialMaster}
	s.alerter = &prodAlerter{dial: s.DialMaster}
	return s
}

func (s *Server) save(ctx context.Context) {
	s.KSclient.Save(ctx, KEY, s.config)
}

func (s *Server) load(ctx context.Context) error {
	config := &pb.Config{}
	data, _, err := s.KSclient.Read(ctx, KEY, config)

	if err != nil {
		return err
	}

	//Clean out the wants here
	wmap := make(map[int32]*pb.MasterWant)
	config = data.(*pb.Config)

	if config.Spends == nil {
		config.Spends = make(map[int32]*pb.RecordSpend)
	}

	for _, want := range config.Wants {
		if val, ok := wmap[want.Release.Id]; ok {
			val.Superwant = val.Superwant || want.Superwant
		} else {
			wmap[want.Release.Id] = want
		}
	}

	s.config = data.(*pb.Config)
	s.config.Wants = []*pb.MasterWant{}
	for _, want := range wmap {
		s.config.Wants = append(s.config.Wants, want)
	}

	if len(config.Wants) > len(s.config.Wants) {
		log.Fatalf("Fatal error loading wants")
	}

	return nil
}

// DoRegister does RPC registration
func (s *Server) DoRegister(server *grpc.Server) {
	pb.RegisterWantServiceServer(server, s)
}

// ReportHealth alerts if we're not healthy
func (s *Server) ReportHealth() bool {
	return true
}

// Shutdown close the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.save(ctx)
	return nil
}

// Mote promotes/demotes this server
func (s *Server) Mote(ctx context.Context, master bool) error {
	if master {
		err := s.load(ctx)
		return err
	}

	return nil
}

func (s *Server) runUpdate(ctx context.Context) error {
	s.alertNoStaging(ctx, s.config.Budget <= 0)
	return nil
}

func (s *Server) getBudget(ctx context.Context) error {
	t := time.Now()
	spends, err := s.GetSpending(ctx, &pb.SpendingRequest{})

	if err == nil {
		s.mmonth = int32(0)
		for _, sp := range spends.Spends {
			if sp.Spend > 0 && sp.Month > s.mmonth {
				s.mmonth = sp.Month
			}
		}

		spendSum := int32(0)
		for _, sp := range spends.Spends {
			if sp.Month > s.mmonth-3 {
				spendSum += sp.Spend
			}
		}

		s.config.Budget = 30000*3 - spendSum
		s.save(ctx)
	} else {
		s.Log(fmt.Sprintf("Error getting spending: %v", err))
	}

	s.budgetPull = time.Now().Sub(t)
	return nil
}

// GetState gets the state of the server
func (s *Server) GetState() []*pbg.State {
	c := 0
	super := int64(0)
	testString := ""
	counts := make(map[int32]int)
	for _, w := range s.config.Wants {
		counts[w.Release.Id]++
		if w.Staged {
			c++
		}

		if w.Superwant {
			super++
		}

		if w.Release.Id == 3216287 {
			testString = fmt.Sprintf("st %v, ac %v, dem %v, super %v", w.Staged, w.Active, w.Demoted, w.Superwant)
		}
	}

	doubleCounts := 0
	for _, val := range counts {
		if val > 1 {
			doubleCounts++
		}
	}

	return []*pbg.State{
		&pbg.State{Key: "spends", Value: int64(len(s.config.Spends))},
		&pbg.State{Key: "wantcount", Value: int64(len(s.config.Wants))},
		&pbg.State{Key: "stagedcount", Value: int64(c)},
		&pbg.State{Key: "supercount", Value: super},
		&pbg.State{Key: "superstring", Text: testString},
		&pbg.State{Key: "lastwantrun", TimeValue: s.lastRun.Unix()},
		&pbg.State{Key: "lastproc", Value: int64(s.lastProc)},
		&pbg.State{Key: "lastpull", Value: int64(s.lastPull)},
		&pbg.State{Key: "pull", Text: s.pull},
		&pbg.State{Key: "budget", Value: int64(s.config.Budget)},
		&pbg.State{Key: "month", Value: int64(s.mmonth)},
		&pbg.State{Key: "last_want", Text: s.lastUnwant},
		&pbg.State{Key: "budget_pull_time", Text: fmt.Sprintf("%v", s.budgetPull)},
		&pbg.State{Key: "double_counts", Value: int64(doubleCounts)},
	}
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

	server.RegisterServer("recordwants", false)

	if *clear {
		ctx, cancel := utils.BuildContext("recordwants", "recordwants")
		defer cancel()

		err := server.load(ctx)
		if err != nil {
			log.Fatalf("Unable to read wants: %v", err)
		}

		for _, w := range server.config.Wants {
			w.Superwant = false
		}

		server.save(ctx)
		log.Fatalf("Saved out cleared wants")
	}

	server.RegisterRepeatingTask(server.updateWants, "update_wants", time.Hour)
	server.RegisterRepeatingTask(server.runUpdate, "run_update", time.Hour)
	server.RegisterRepeatingTask(server.getBudget, "get_budget", time.Hour)

	server.Serve()
}
