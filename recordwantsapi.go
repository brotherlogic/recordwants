package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	pbgd "github.com/brotherlogic/godiscogs/proto"
	qpb "github.com/brotherlogic/queue/proto"
	rcpb "github.com/brotherlogic/recordcollection/proto"
	pb "github.com/brotherlogic/recordwants/proto"
	google_protobuf "github.com/golang/protobuf/ptypes/any"
)

// AddWant adds a want into the system
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
			Release:      &pbgd.Release{Id: req.ReleaseId},
			Level:        req.Level,
			RetireTime:   req.RetireTime,
			Budget:       req.GetBudget(),
			DesiredState: pb.MasterWant_WANTED,
		})

	return &pb.AddWantResponse{}, s.save(ctx, config)

}

// GetWants gets a want
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

// Update updates a given want
func (s *Server) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	if req.GetReason() == "" {
		return nil, fmt.Errorf("you must provide a reason for this update")
	}

	config, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	for _, want := range config.Wants {
		if want.GetRelease().Id == req.GetWant().Id {
			want.DesiredState = req.GetNewState()
			want.RetireTime = req.GetRetireTime()
			want.Budget = req.GetBudget()
			return &pb.UpdateResponse{}, s.save(ctx, config)
		}
	}
	return nil, fmt.Errorf("not found: %v", config.Wants)
}

// ClientUpdate on an updated record
func (s *Server) ClientUpdate(ctx context.Context, req *rcpb.ClientUpdateRequest) (*rcpb.ClientUpdateResponse, error) {
	err := s.updateWants(ctx, req.GetInstanceId())
	if err != nil {
		return nil, err
	}

	_, err = s.Sync(ctx, &pb.SyncRequest{Soft: true})
	return &rcpb.ClientUpdateResponse{}, err
}

func (s *Server) Sync(ctx context.Context, req *pb.SyncRequest) (*pb.SyncResponse, error) {
	config, err := s.load(ctx)
	if err != nil {
		return nil, err
	}

	if !req.GetSoft() {
		// Pull in existing wants
		wants, err := s.recordGetter.getWants(ctx)
		if err != nil {
			return nil, err
		}

		processed := make(map[int32]bool)
		for _, want := range wants {
			for _, in := range config.GetWants() {
				if in.GetRelease().GetId() == want.GetReleaseId() {
					in.CurrentState = pb.MasterWant_WANTED
					processed[want.GetReleaseId()] = true
				}
			}

			if !processed[want.GetReleaseId()] {
				// This is a new release
				config.Wants = append(config.Wants, &pb.MasterWant{
					DateAdded:    time.Now().Unix(),
					CurrentState: pb.MasterWant_WANTED,
					DesiredState: pb.MasterWant_UNWANTED,
					Release:      &pbgd.Release{Id: want.GetReleaseId()}})
			}
		}

		// Process anything we've missed
		for _, want := range config.GetWants() {
			if !processed[want.GetRelease().GetId()] && want.GetRelease().GetId() != 0 {
				want.CurrentState = pb.MasterWant_UNWANTED
			}
		}
	}

	err = s.updateWantState(ctx, config)
	if err != nil {
		return nil, err
	}

	err = s.save(ctx, config)
	if err != nil {
		return nil, err
	}

	if !req.GetSoft() {
		conn2, err2 := s.FDialServer(ctx, "queue")
		if err2 != nil {
			return nil, err2
		}
		defer conn2.Close()
		qclient := qpb.NewQueueServiceClient(conn2)
		syncreq := &pb.SyncRequest{}
		data, _ := proto.Marshal(syncreq)
		_, err3 := qclient.AddQueueItem(ctx, &qpb.AddQueueItemRequest{
			QueueName:     "wants_sync",
			RunTime:       time.Now().Add(time.Hour).Unix(),
			Payload:       &google_protobuf.Any{Value: data},
			Key:           "syncer",
			RequireUnique: true,
		})
		return &pb.SyncResponse{}, err3
	}

	return &pb.SyncResponse{}, err
}
