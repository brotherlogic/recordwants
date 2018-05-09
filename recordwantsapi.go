package main

import (
	"time"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/recordwants/proto"
)

//GetSpending gets the spending over the course of months
func (s *Server) GetSpending(ctx context.Context, req *pb.SpendingRequest) (*pb.SpendingResponse, error) {
	r, err := s.recordGetter.getRecords()
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

	return resp, nil
}
