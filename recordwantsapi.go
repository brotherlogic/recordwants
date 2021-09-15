package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbgd "github.com/brotherlogic/godiscogs"
	rcpb "github.com/brotherlogic/recordcollection/proto"
	pb "github.com/brotherlogic/recordwants/proto"
)

//AddWant adds a want into the system
func (s *Server) AddWant(ctx context.Context, req *pb.AddWantRequest) (*pb.AddWantResponse, error) {
	config, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	for _, w := range config.Wants {
		if req.ReleaseId == w.GetRelease().GetId() {
			return &pb.AddWantResponse{}, status.Errorf(codes.FailedPrecondition, "Want already exists")
		}
	}

	config.Wants = append(config.Wants,
		&pb.MasterWant{
			Release:     &pbgd.Release{Id: req.ReleaseId},
			Level:       req.Level,
			RetireTime:  req.RetireTime,
			RetireLevel: req.RetireLevel,
		})

	return &pb.AddWantResponse{}, s.save(ctx, config)

}

//GetWants gets a want
func (s *Server) GetWants(ctx context.Context, req *pb.GetWantsRequest) (*pb.GetWantsResponse, error) {
	config, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	wants := []*pb.MasterWant{}
	for _, w := range config.Wants {
		found := len(req.GetReleaseId()) == 0
		for _, wr := range req.GetReleaseId() {
			if w.GetRelease().GetId() == wr {
				found = true
			}
		}

		if found {
			wants = append(wants, w)
		}
	}

	return &pb.GetWantsResponse{Want: wants}, nil
}

//GetSpending gets the spending over the course of months
func (s *Server) GetSpending(ctx context.Context, req *pb.SpendingRequest) (*pb.SpendingResponse, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

//Update updates a given want
func (s *Server) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	config, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	for _, want := range config.Wants {
		if want.GetRelease().Id == req.GetWant().Id {
			want.Staged = true
			want.Demoted = !req.KeepWant
			want.Superwant = req.Super
			if req.GetLevel() != pb.MasterWant_UNKNOWN {
				want.Level = req.GetLevel()
			}
			want.RetireTime = req.GetRetireTime()
			s.Log(fmt.Sprintf("Updated want: %v", want))
			return &pb.UpdateResponse{}, s.save(ctx, config)
		}
	}
	return nil, fmt.Errorf("Not found: %v", config.Wants)
}

//ClientUpdate on an updated record
func (s *Server) ClientUpdate(ctx context.Context, req *rcpb.ClientUpdateRequest) (*rcpb.ClientUpdateResponse, error) {
	t := time.Now()
	defer func() {
		s.Log(fmt.Sprintf("Client Update in %v", time.Now().Sub(t)))
	}()
	err := s.updateWants(ctx, req.GetInstanceId())
	if err != nil {
		return nil, err
	}
	return &rcpb.ClientUpdateResponse{}, s.updateWantState(ctx)
}

func (s *Server) Sync(ctx context.Context, req *pb.SyncRequest) (*pb.SyncResponse, error) {
	config, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	// Pull in existing wants
	wants, err := s.recordGetter.getWants(ctx)
	if err != nil {
		return nil, err
	}

	processed := make(map[int32]bool)
	for _, want := range wants {
		for _, in := range config.GetWants() {
			if in.GetRelease().GetId() == want.GetReleaseId() {
				if in.GetRelease().GetId() == 0 {
					in.Level = pb.MasterWant_NEVER
				} else {
					processed[want.GetReleaseId()] = true
					if in.GetLevel() != pb.MasterWant_ANYTIME_LIST {
						err := s.recordGetter.unwant(ctx, in)
						if err != nil && status.Convert(err).Code() != codes.NotFound {
							return nil, err
						}
					}
				}
			}
		}

		if !processed[want.GetReleaseId()] {
			// This is a new release
			config.Wants = append(config.Wants, &pb.MasterWant{Level: pb.MasterWant_UNKNOWN, Release: &pbgd.Release{Id: want.GetReleaseId()}})
		}
	}

	// Process anything we've missed
	for _, want := range config.GetWants() {
		if !processed[want.GetRelease().GetId()] && want.GetRelease().GetId() != 0 {
			if want.GetLevel() == pb.MasterWant_ANYTIME_LIST {
				err := s.recordGetter.want(ctx, want)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return &pb.SyncResponse{}, s.save(ctx, config)
}
