package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	pbgd "github.com/brotherlogic/godiscogs"
	pb "github.com/brotherlogic/recordwants/proto"
	pbt "github.com/brotherlogic/tracer/proto"
)

func (s *Server) AddWant(ctx context.Context, req *pb.AddWantRequest) (*pb.AddWantResponse, error) {
	ctx = s.LogTrace(ctx, "AddWant", time.Now(), pbt.Milestone_START_FUNCTION)

	s.config.Wants = append(s.config.Wants,
		&pb.MasterWant{
			Superwant: req.Superwant,
			Release:   &pbgd.Release{Id: req.ReleaseId},
		})

	s.LogTrace(ctx, "AddWant", time.Now(), pbt.Milestone_END_FUNCTION)
	return &pb.AddWantResponse{}, nil

}

//GetSpending gets the spending over the course of months
func (s *Server) GetSpending(ctx context.Context, req *pb.SpendingRequest) (*pb.SpendingResponse, error) {
	ctx = s.LogTrace(ctx, "GetSpending", time.Now(), pbt.Milestone_START_FUNCTION)
	r, err := s.recordGetter.getRecords(ctx)
	if err != nil {
		return nil, err
	}

	s.Log(fmt.Sprintf("SPENDS %v", len(r)))

	resp := &pb.SpendingResponse{Spends: make([]*pb.Spend, 0)}
	for i := 0; i < 12; i++ {
		resp.Spends = append(resp.Spends, &pb.Spend{Month: int32(i)})
	}
	for _, r := range r {
		d := time.Unix(r.Metadata.DateAdded, 0)
		if d.Year() == 2018 {
			resp.Spends[int(d.Month()-1)].Spend += r.Metadata.Cost
		}
	}

	s.LogTrace(ctx, "GetSpending", time.Now(), pbt.Milestone_END_FUNCTION)
	return resp, nil
}

//Update updates a given want
func (s *Server) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	ctx = s.LogTrace(ctx, "Update", time.Now(), pbt.Milestone_START_FUNCTION)
	for _, want := range s.config.Wants {
		if want.GetRelease().Id == req.GetWant().Id {
			want.Staged = true
			want.Demoted = !req.KeepWant
			want.Superwant = req.Super
			return &pb.UpdateResponse{}, nil
		}
	}
	s.LogTrace(ctx, "Update", time.Now(), pbt.Milestone_END_FUNCTION)
	return nil, fmt.Errorf("Not found: %v", s.config.Wants)
}
