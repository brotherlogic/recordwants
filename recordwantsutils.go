package main

import (
	"fmt"
	"time"

	pb "github.com/brotherlogic/recordwants/proto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	processed = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "recordwants_processed",
		Help: "The size of the wants queue",
	})
)

func (s *Server) updateWantState(ctx context.Context) error {
	config, err := s.load(ctx)
	if err != nil {
		return err
	}

	for _, want := range config.Wants {
		err, done := s.updateWant(ctx, want, time.Now())
		if err != nil {
			return err
		}
		if done {
			break
		}

	}

	return s.save(ctx, config)
}

func (s *Server) updateWant(ctx context.Context, want *pb.MasterWant, ti time.Time) (error, bool) {
	s.CtxLog(ctx, fmt.Sprintf("Updating %v", want.GetRelease().GetId()))
	if want.GetDirty() {
		return nil, false
	}

	if want.GetRetireTime() > 0 {
		if ti.After(time.Unix(want.GetRetireTime(), 0)) {
			want.DesiredState = pb.MasterWant_UNWANTED
		}
	}

	if want.GetDesiredState() != want.GetCurrentState() {
		if want.GetDesiredState() == pb.MasterWant_WANTED {
			err := s.recordGetter.want(ctx, want)
			if err != nil {
				return err, false
			}
			want.CurrentState = want.GetDesiredState()
		} else {
			err := s.recordGetter.unwant(ctx, want)

			if err != nil && status.Convert(err).Code() != codes.NotFound {
				return err, false
			}
			want.CurrentState = pb.MasterWant_UNWANTED
		}

		return nil, true
	}

	return nil, false
}

func (s *Server) updateWants(ctx context.Context, iid int32) error {
	config, err := s.load(ctx)
	if err != nil {
		return err
	}

	// Demote any wants we already own
	record, err := s.recordGetter.getRecord(ctx, iid)
	if err != nil {
		if status.Convert(err).Code() == codes.OutOfRange {
			return nil
		}
		return err
	}

	// Let's validate the budget for this want (we only started doing strict budgeting in 2022)
	if time.Unix(record.GetMetadata().GetDateAdded(), 0).Year() > 2021 && record.GetMetadata().GetPurchaseBudget() == "" {
		// Find the budget
		for _, want := range config.GetWants() {
			if want.Release.Id == record.GetRelease().GetId() {
				if want.GetBudget() != "" {
					err := s.recordGetter.updateBudget(ctx, iid, want.GetBudget())
					if err != nil {
						return err
					}
				}
			}
		}

		if record.GetMetadata().GetPurchaseBudget() == "" {
			return status.Errorf(codes.FailedPrecondition, "this purchase (%v) has no budget that we can find", record.GetRelease().GetInstanceId())
		}
	}

	for _, w := range config.Wants {
		if w.Level != pb.MasterWant_BOUGHT && record.GetRelease().GetId() == w.GetRelease().GetId() {
			w.Demoted = true
			w.Staged = true
			w.Level = pb.MasterWant_BOUGHT
		}
	}

	return s.save(ctx, config)
}
