package main

import (
	"fmt"
	"time"

	pb "github.com/brotherlogic/recordwants/proto"
	"golang.org/x/net/context"
)

func (s *Server) updateWantState(ctx context.Context) error {
	config, err := s.load(ctx)
	if err != nil {
		return err
	}

	for _, want := range config.Wants {
		err := s.updateWant(ctx, want)
		if err != nil {
			return err
		}
	}

	return s.save(ctx, config)
}

func (s *Server) updateWant(ctx context.Context, want *pb.MasterWant) error {

	if want.GetDirty() {
		return nil
	}

	switch want.GetLevel() {
	case pb.MasterWant_UNKNOWN, pb.MasterWant_BOUGHT:
		if want.GetActive() {
			err := s.recordGetter.unwant(ctx, want)
			if err != nil {
				return err
			}
			want.Dirty = true
		}
	case pb.MasterWant_ALWAYS:
		if !want.GetActive() {
			err := s.recordGetter.want(ctx, want)
			if err != nil {
				return err
			}
			want.Dirty = true
		}
	case pb.MasterWant_ANYTIME:
		budget, err := s.getBudget(ctx)
		s.Log(fmt.Sprintf("%v and budget %v, %v", want, budget, err))
		if err != nil {
			return err
		}
		if !want.GetActive() && budget > 0 {
			err := s.recordGetter.want(ctx, want)
			if err != nil {
				return err
			}
			want.Dirty = true
		} else if want.GetActive() && budget < 0 {
			err := s.recordGetter.unwant(ctx, want)
			if err != nil {
				return err
			}
			want.Dirty = true
		}
	case pb.MasterWant_LIST:
		if !want.GetActive() {
			err := s.recordGetter.want(ctx, want)
			if err != nil {
				return err
			}
			want.Dirty = true
		}

	case pb.MasterWant_STAGED_TO_BE_ADDED:
		if want.GetActive() {
			err := s.recordGetter.unwant(ctx, want)
			if err != nil {
				return err
			}
			want.Dirty = true
		}
	}

	return nil
}

func (s *Server) updateWants(ctx context.Context) error {
	config, err := s.load(ctx)
	if err != nil {
		return err
	}
	s.lastRun = time.Now()
	wants, err := s.recordGetter.getWants(ctx)
	s.lastPull = int32(len(wants))
	if err == nil {
		for _, w := range wants {
			found := false
			for _, mw := range config.Wants {
				if mw.Release.Id == w.Release.Id {
					found = true
					mw.Active = w.GetMetadata().GetActive()
					mw.Dirty = false
				}
			}
			if !found {
				config.Wants = append(config.Wants,
					&pb.MasterWant{Release: w.Release, DateAdded: time.Now().Unix()})
			}
		}
	}

	// Demote any wants we already own
	for _, w := range config.Wants {
		if w.Level != pb.MasterWant_BOUGHT {
			records, err := s.recordGetter.getRecords(ctx, w.GetRelease().Id)
			if err == nil && len(records) > 0 {
				w.Demoted = true
				w.Staged = true
				w.Level = pb.MasterWant_BOUGHT
			}
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
