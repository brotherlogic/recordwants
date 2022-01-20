package main

import (
	"fmt"
	"time"

	pbgd "github.com/brotherlogic/godiscogs"
	pb "github.com/brotherlogic/recordwants/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) updateWantState(ctx context.Context) error {
	config, err := s.load(ctx)
	if err != nil {
		return err
	}

	for _, want := range config.Wants {
		budget, err := s.getBudget(ctx, want.GetBudget())
		if err != nil {
			return err
		}
		err = s.updateWant(ctx, want, budget, time.Now())
		if err != nil {
			return err
		}

	}

	return s.save(ctx, config)
}

func (s *Server) updateWant(ctx context.Context, want *pb.MasterWant, budget int32, ti time.Time) error {
	s.CtxLog(ctx, fmt.Sprintf("Updating %v with %v", want.GetRelease().GetId(), budget))
	if want.GetDirty() {
		return nil
	}

	if want.GetRetireTime() > 0 {
		if ti.After(time.Unix(want.GetRetireTime(), 0)) {
			want.Dirty = true
			want.Level = want.GetRetireLevel()
		}
	}

	switch want.GetLevel() {
	case pb.MasterWant_UNKNOWN, pb.MasterWant_BOUGHT, pb.MasterWant_NEVER:
		if want.GetActive() {
			err := s.recordGetter.unwant(ctx, want)
			if err != nil && status.Convert(err).Code() != codes.NotFound {
				return err
			}
		}
	case pb.MasterWant_ALWAYS:
		if !want.GetActive() {
			err := s.recordGetter.want(ctx, want)
			if err != nil {
				return err
			}
		}
	case pb.MasterWant_WANT_DIGITAL:
		if !want.GetActive() && budget > 50 {
			recs, err := s.recordGetter.getAllRecords(ctx, want.GetRelease().GetId())
			if err != nil {
				return err
			}
			for _, r := range recs {
				err := s.recordGetter.want(ctx, &pb.MasterWant{Release: &pbgd.Release{Id: r}})
				if err != nil {
					return err
				}
			}
		} else if want.GetActive() && budget <= 50 {
			recs, err := s.recordGetter.getAllRecords(ctx, want.GetRelease().GetId())
			if err != nil {
				return err
			}
			for _, r := range recs {
				err := s.recordGetter.unwant(ctx, &pb.MasterWant{Release: &pbgd.Release{Id: r}})
				if err != nil && status.Convert(err).Code() != codes.NotFound {
					return err
				}
			}
		}
	case pb.MasterWant_WANT_OG:
		if !want.GetActive() && budget > 0 {
			err := s.recordGetter.want(ctx, want)
			if err != nil {
				return err
			}
			want.Dirty = true
		} else if want.GetActive() && budget <= 0 {
			err := s.recordGetter.unwant(ctx, want)
			if err != nil && status.Convert(err).Code() != codes.NotFound {
				return err
			}
			want.Dirty = true
		}
	case pb.MasterWant_ANYTIME:
		if !want.GetActive() && budget >= 0 {
			err := s.recordGetter.want(ctx, want)
			if err != nil {
				return err
			}
			want.Dirty = true
		} else if want.GetActive() && budget <= 0 {
			err := s.recordGetter.unwant(ctx, want)
			if err != nil && status.Convert(err).Code() != codes.NotFound {
				return err
			}
			want.Dirty = true
		}
	case pb.MasterWant_LIST, pb.MasterWant_ANYTIME_LIST:
		baseline := int32(0)
		if want.GetLevel() == pb.MasterWant_ANYTIME_LIST {
			baseline = 0
		}
		if !want.GetActive() && budget > baseline {
			err := s.recordGetter.want(ctx, want)
			if err != nil {
				return err
			}
			want.Dirty = true
		} else if want.GetActive() && budget <= baseline {
			err := s.recordGetter.unwant(ctx, want)
			if err != nil && status.Convert(err).Code() != codes.NotFound {
				return err
			}
			want.Dirty = true
		}
	case pb.MasterWant_STAGED_TO_BE_ADDED:
		if want.GetActive() {
			err := s.recordGetter.unwant(ctx, want)
			if err != nil && status.Convert(err).Code() != codes.NotFound {
				return err
			}
			want.Dirty = true
		}
	default:
		s.RaiseIssue("Cannot handle want", fmt.Sprintf("Have no means of processing %v level want", want.GetLevel()))
	}

	return nil
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
			return fmt.Errorf("this purchase (%v) has no budget that we can find", record.GetRelease().GetInstanceId())
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

func (s *Server) dealWithAddedRecords(ctx context.Context) error {
	nums, err := s.recordAdder.getAdds(ctx)
	if err != nil {
		return err
	}

	config, err := s.load(ctx)
	if err != nil {
		return err
	}

	for _, num := range nums {
		for _, w := range config.Wants {
			if w.GetRelease().Id == num {
				w.Level = pb.MasterWant_STAGED_TO_BE_ADDED
				w.Demoted = true
				w.Staged = true
			}
		}
	}

	return s.save(ctx, config)
}
