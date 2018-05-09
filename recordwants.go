package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/brotherlogic/goserver"
	"google.golang.org/grpc"

	pbg "github.com/brotherlogic/goserver/proto"
	pbrc "github.com/brotherlogic/recordcollection/proto"
)

type recordGetter interface {
	getRecords() ([]*pbrc.Record, error)
}

//Server main server type
type Server struct {
	*goserver.GoServer
	recordGetter recordGetter
}

// Init builds the server
func Init() *Server {
	s := &Server{GoServer: &goserver.GoServer{}}
	return s
}

// DoRegister does RPC registration
func (s *Server) DoRegister(server *grpc.Server) {
	// Do nothing
}

// ReportHealth alerts if we're not healthy
func (s *Server) ReportHealth() bool {
	return true
}

// Mote promotes/demotes this server
func (s *Server) Mote(master bool) error {
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
	server.PrepServer()
	server.Register = server

	server.RegisterServer("recordwants", false)
	server.Log("Starting!")
	server.Serve()
}
