package main

import (
	"fmt"

	"golang.org/x/net/context"

	pbgd "github.com/brotherlogic/godiscogs"
	pb "github.com/brotherlogic/recordwants/proto"
)

//AddWant adds a want into the system
func (s *Server) AddWant(ctx context.Context, req *pb.AddWantRequest) (*pb.AddWantResponse, error) {
	s.config.Wants = append(s.config.Wants,
		&pb.MasterWant{
			Superwant: req.Superwant,
			Release:   &pbgd.Release{Id: req.ReleaseId},
		})

	return &pb.AddWantResponse{}, nil

}

//GetSpending gets the spending over the course of months
func (s *Server) GetSpending(ctx context.Context, req *pb.SpendingRequest) (*pb.SpendingResponse, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

//Update updates a given want
func (s *Server) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	for _, want := range s.config.Wants {
		if want.GetRelease().Id == req.GetWant().Id {
			want.Staged = true
			want.Demoted = !req.KeepWant
			want.Superwant = req.Super
			s.Log(fmt.Sprintf("Updated want: %v", want))
			return &pb.UpdateResponse{}, nil
		}
	}
	return nil, fmt.Errorf("Not found: %v", s.config.Wants)
}
