package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/brotherlogic/goserver"
	"github.com/brotherlogic/goserver/utils"
	"github.com/brotherlogic/keystore/client"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pbgh "github.com/brotherlogic/githubcard/proto"
	pbg "github.com/brotherlogic/goserver/proto"
	pbrc "github.com/brotherlogic/recordcollection/proto"
	pb "github.com/brotherlogic/recordwants/proto"
	pbt "github.com/brotherlogic/tracer/proto"
)

const (
	// KEY - where the wants are stored
	KEY = "/github.com/brotherlogic/recordwants/config"
)

type alerter interface {
	alert(ctx context.Context, want *pb.MasterWant)
}

type prodAlerter struct{}

func (p *prodAlerter) alert(ctx context.Context, want *pb.MasterWant) {
	ip, port, _ := utils.Resolve("githubcard")
	if port > 0 {
		conn, err := grpc.Dial(ip+":"+strconv.Itoa(int(port)), grpc.WithInsecure())
		if err == nil {
			defer conn.Close()
			client := pbgh.NewGithubClient(conn)
			client.AddIssue(ctx, &pbgh.Issue{Service: "recordwants", Title: fmt.Sprintf("Want Processing Needed!"), Body: fmt.Sprintf("%v", want)}, grpc.FailFast(false))
		}
	}
}

type recordGetter interface {
	getRecords(ctx context.Context) ([]*pbrc.Record, error)
	getWants(ctx context.Context) ([]*pbrc.Want, error)
	unwant(ctx context.Context, want *pb.MasterWant) error
	want(ctx context.Context, want *pb.MasterWant) error
}

type prodGetter struct{}

func (p *prodGetter) getRecords(ctx context.Context) ([]*pbrc.Record, error) {
	ip, port, err := utils.Resolve("recordcollection")
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(ip+":"+strconv.Itoa(int(port)), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pbrc.NewRecordCollectionServiceClient(conn)
	utils.SendTrace(ctx, "Calling Get Records", time.Now(), pbt.Milestone_MARKER, "recordwants")
	resp, err := client.GetRecords(ctx, &pbrc.GetRecordsRequest{Strip: true, Filter: &pbrc.Record{Metadata: &pbrc.ReleaseMetadata{}}}, grpc.MaxCallRecvMsgSize(1024*1024*1024))
	if err != nil {
		return nil, err
	}
	return resp.GetRecords(), nil

}

func (p *prodGetter) getWants(ctx context.Context) ([]*pbrc.Want, error) {
	ip, port, err := utils.Resolve("recordcollection")
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(ip+":"+strconv.Itoa(int(port)), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pbrc.NewRecordCollectionServiceClient(conn)
	utils.SendTrace(ctx, "Calling Get Wants", time.Now(), pbt.Milestone_MARKER, "recordwants")
	resp, err := client.GetWants(ctx, &pbrc.GetWantsRequest{})
	if err != nil {
		return nil, err
	}

	return resp.GetWants(), nil

}

func (p *prodGetter) unwant(ctx context.Context, want *pb.MasterWant) error {
	ip, port, err := utils.Resolve("recordcollection")
	if err != nil {
		return err
	}

	conn, err := grpc.Dial(ip+":"+strconv.Itoa(int(port)), grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pbrc.NewRecordCollectionServiceClient(conn)
	_, err = client.UpdateWant(ctx, &pbrc.UpdateWantRequest{Update: &pbrc.Want{Release: want.GetRelease()}, Remove: true})
	return err
}

func (p *prodGetter) want(ctx context.Context, want *pb.MasterWant) error {
	ip, port, err := utils.Resolve("recordcollection")
	if err != nil {
		return err
	}

	conn, err := grpc.Dial(ip+":"+strconv.Itoa(int(port)), grpc.WithInsecure())
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
	}
	return s
}

func (s *Server) save() {
	s.KSclient.Save(KEY, s.config)
}

func (s *Server) load() error {
	config := &pb.Config{}
	data, _, err := s.KSclient.Read(KEY, config)

	if err != nil {
		return err
	}

	s.config = data.(*pb.Config)
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

// Mote promotes/demotes this server
func (s *Server) Mote(master bool) error {
	if master {
		err := s.load()
		return err
	}

	return nil
}

func (s *Server) runUpdate(ctx context.Context) {
	s.alertNoStaging(ctx, s.config.Budget > 0)
}

func (s *Server) getBudget(ctx context.Context) {
	spends, err := s.GetSpending(ctx, &pb.SpendingRequest{})

	if err == nil {
		mmonth := int32(6)

		spendSum := int32(0)
		for _, sp := range spends.Spends {
			if sp.Month > mmonth-3 {
				spendSum += sp.Spend
			}
		}

		s.config.Budget = 40000*3 - spendSum
	}
}

// GetState gets the state of the server
func (s *Server) GetState() []*pbg.State {
	c := 0
	found := 0
	stat := ""
	for _, w := range s.config.Wants {
		if w.Staged {
			c++
		}
		if w.GetRelease().Id == 3379533 {
			found++
			stat = fmt.Sprintf("%v", w)
		}
	}
	return []*pbg.State{
		&pbg.State{Key: "wantcount", Value: int64(len(s.config.Wants))},
		&pbg.State{Key: "stagedcount", Value: int64(c)},
		&pbg.State{Key: "lastwantrun", TimeValue: s.lastRun.Unix()},
		&pbg.State{Key: "lastproc", Value: int64(s.lastProc)},
		&pbg.State{Key: "lastpull", Value: int64(s.lastPull)},
		&pbg.State{Key: "found", Value: int64(found)},
		&pbg.State{Key: "stat", Text: stat},
		&pbg.State{Key: "pull", Text: s.pull},
		&pbg.State{Key: "budget", Value: int64(s.config.Budget)},
	}
}

func main() {
	var quiet = flag.Bool("quiet", false, "Show all output")
	flag.Parse()

	//Turn off logging
	if *quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	server := Init()
	server.GoServer.KSclient = *keystoreclient.GetClient(server.GetIP)
	server.PrepServer()
	server.Register = server

	server.RegisterServer("recordwants", false)
	server.RegisterRepeatingTask(server.updateWants, time.Minute*5)
	server.RegisterRepeatingTask(server.runUpdate, time.Hour)
	server.RegisterRepeatingTask(server.getBudget, time.Hour)
	server.Log("Starting!")
	server.Serve()
}
