package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/recordwants/proto"
	pbt "github.com/brotherlogic/tracer/proto"
)

//GetSpending gets the spending over the course of months
func (s *Server) GetSpending(ctx context.Context, req *pb.SpendingRequest) (*pb.SpendingResponse, error) {
	ctx = s.LogTrace(ctx, "GetSpending", time.Now(), pbt.Milestone_START_FUNCTION)
	r, err := s.recordGetter.getRecords(ctx)
	if err != nil {
		return nil, err
	}

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
	for _, want := range s.config.Wants {
		if want.GetRelease().Id == req.GetWant().Id {
			if req.KeepWant {
				err := s.recordGetter.want(ctx, &pb.MasterWant{Release: req.GetWant()})
				if err != nil {
					return nil, err
				}
			}
			want.Staged = false
			want.Enable = req.KeepWant
		}
	}
	return nil, fmt.Errorf("Not found: %v", s.config.Wants)
}
