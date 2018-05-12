package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/brotherlogic/goserver"
	"github.com/brotherlogic/goserver/utils"
	"github.com/brotherlogic/keystore/client"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pbg "github.com/brotherlogic/goserver/proto"
	pbrc "github.com/brotherlogic/recordcollection/proto"
	pb "github.com/brotherlogic/recordwants/proto"
	pbt "github.com/brotherlogic/tracer/proto"
)

const (
	// KEY - where the wants are stored
	KEY = "/github.com/brotherlogic/recordwants/config"
)

type recordGetter interface {
	getRecords(ctx context.Context) ([]*pbrc.Record, error)
	getWants(ctx context.Context) ([]*pbrc.Want, error)
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

//Server main server type
type Server struct {
	*goserver.GoServer
	recordGetter recordGetter
	config       *pb.Config
}

// Init builds the server
func Init() *Server {
	s := &Server{GoServer: &goserver.GoServer{}}
	s.recordGetter = &prodGetter{}
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

// GetState gets the state of the server
func (s *Server) GetState() []*pbg.State {
	return []*pbg.State{}
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
	server.Log("Starting!")
	server.Serve()
}
