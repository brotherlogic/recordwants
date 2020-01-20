package main

import (
	"fmt"
	"sort"
	"time"

	pb "github.com/brotherlogic/recordwants/proto"
	"golang.org/x/net/context"
)

func (s *Server) updateSpending(ctx context.Context) error {
	recs, err := s.recordGetter.getRecordsSince(ctx, s.config.LastSpendUpdate)
	if err != nil {
		return err
	}

	for _, id := range recs {
		rec, err := s.recordGetter.getRecord(ctx, id)
		if err != nil {
			return err
		}

		s.config.Spends[id] = &pb.RecordSpend{Cost: rec.GetMetadata().Cost, DateAdded: rec.GetMetadata().DateAdded}
	}

	return nil
}

func (s *Server) alertNoStaging(ctx context.Context, overBudget bool) {

	//Order the wants in time order
	sort.SliceStable(s.config.Wants, func(i, j int) bool {
		return s.config.Wants[i].DateAdded > s.config.Wants[j].DateAdded
	})

	for _, want := range s.config.Wants {
		if !want.Staged && !want.Superwant {
			if want.Active {
				s.Log(fmt.Sprintf("Unwanting %v as it's not staged but active", want.Release.Id))
				s.config.LastPush = time.Now().Unix()
				s.recordGetter.unwant(ctx, want)
			}
			if !want.Demoted {
				c := 0
				for _, w := range s.config.Wants {
					if w.Staged {
						c++
					}
				}

				s.alerter.alert(ctx, want, c, len(s.config.Wants))
			}

		} else {
			if overBudget && want.Active && !want.Superwant {
				s.lastProc = want.Release.Id
				s.Log(fmt.Sprintf("Unwanting %v because we're over budget", want.Release.Id))
				s.config.LastPush = time.Now().Unix()
				s.recordGetter.unwant(ctx, want)
			} else {
				if want.Active && want.Demoted {
					s.Log(fmt.Sprintf("Unwanting %v because it's demoted", want.Release.Id))
					s.config.LastPush = time.Now().Unix()
					s.recordGetter.unwant(ctx, want)
				}
				if !want.Active && !want.Demoted && (!overBudget || want.Superwant) {
					s.Log(fmt.Sprintf("WANTING %v -> %v, %v, %v", want.Release.Id, want.Active, want.Demoted, overBudget))
					s.config.LastPush = time.Now().Unix()
					s.recordGetter.want(ctx, want)
				}
			}
		}
	}
}

func (s *Server) updateWants(ctx context.Context) error {
	s.lastRun = time.Now()
	wants, err := s.recordGetter.getWants(ctx)
	s.lastPull = int32(len(wants))
	if err == nil {
		for _, w := range wants {
			found := false
			for _, mw := range s.config.Wants {
				if mw.Release.Id == w.Release.Id {
					found = true
					mw.Active = w.GetMetadata().Active
				}
			}
			if !found {
				s.config.Wants = append(s.config.Wants,
					&pb.MasterWant{Release: w.Release, DateAdded: time.Now().Unix()})
			}
		}
	}

	// Demote any wants we already own
	for _, w := range s.config.Wants {
		if !w.Demoted {
			records, err := s.recordGetter.getRecords(ctx, w.GetRelease().Id)
			if err == nil && len(records) > 0 {
				w.Demoted = true
				w.Staged = true
			}
		}
	}

	s.save(ctx)
	return nil
}

func (s *Server) dealWithAddedRecords(ctx context.Context) error {
	nums, err := s.recordAdder.getAdds(ctx)
	if err != nil {
		return err
	}

	for _, num := range nums {
		for _, w := range s.config.Wants {
			if w.GetRelease().Id == num {
				w.Demoted = true
				w.Staged = true
			}
		}
	}

	return s.save(ctx)
}
