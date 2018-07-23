package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/recordwants/proto"
)

func (s *Server) alertNoStaging(ctx context.Context, overBudget bool) {
	for _, want := range s.config.Wants {
		if !want.Staged {
			if want.Active {
				s.lastUnwant = fmt.Sprintf("Unwanting %v", want.Release.Id)
				s.recordGetter.unwant(ctx, want)
			}
			if !want.Demoted {
				s.alerter.alert(ctx, want)
			}
		} else {
			if overBudget && want.Active {
				s.lastProc = want.Release.Id
				s.recordGetter.unwant(ctx, want)
			} else {
				if want.Active && want.Demoted {
					s.recordGetter.unwant(ctx, want)
				}
				if !want.Active && !want.Demoted {
					s.recordGetter.want(ctx, want)
				}
			}
		}
	}
}

func (s *Server) updateWants(ctx context.Context) {
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
	records, err := s.recordGetter.getRecords(ctx)
	if err == nil {
		for _, w := range s.config.Wants {
			if !w.Demoted {
				for _, r := range records {
					if r.GetRelease().Id == w.GetRelease().Id {
						w.Demoted = true
						w.Staged = true
						break
					}
				}
			}
		}
	}

	s.save()
}
