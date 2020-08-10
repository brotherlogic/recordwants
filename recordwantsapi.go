package main

import (
	"fmt"

	"golang.org/x/net/context"

	pbgd "github.com/brotherlogic/godiscogs"
	rcpb "github.com/brotherlogic/recordcollection/proto"
	pb "github.com/brotherlogic/recordwants/proto"
)

//AddWant adds a want into the system
func (s *Server) AddWant(ctx context.Context, req *pb.AddWantRequest) (*pb.AddWantResponse, error) {
	for _, w := range s.config.Wants {
		if req.ReleaseId == w.GetRelease().GetId() {
			return &pb.AddWantResponse{}, nil
		}
	}

	s.config.Wants = append(s.config.Wants,
		&pb.MasterWant{
			Superwant: req.Superwant,
			Release:   &pbgd.Release{Id: req.ReleaseId},
		})

	return &pb.AddWantResponse{}, nil

}

//GetWant gets a want
func (s *Server) GetWant(ctx context.Context, req *pb.GetWantRequest) (*pb.GetWantResponse, error) {
	for _, w := range s.config.Wants {
		if w.GetRelease().GetId() == req.GetReleaseId() {
			return &pb.GetWantResponse{Want: w}, nil
		}
	}

	return nil, fmt.Errorf("Could not locate want")
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
			if req.GetLevel() != pb.MasterWant_UNKNOWN {
				want.Level = req.GetLevel()
			}
			s.Log(fmt.Sprintf("Updated want: %v", want))
			return &pb.UpdateResponse{}, s.save(ctx)
		}
	}
	return nil, fmt.Errorf("Not found: %v", s.config.Wants)
}

//ClientUpdate on an updated record
func (s *Server) ClientUpdate(ctx context.Context, req *rcpb.ClientUpdateRequest) (*rcpb.ClientUpdateResponse, error) {
	return &rcpb.ClientUpdateResponse{}, s.updateWants(ctx)
}
