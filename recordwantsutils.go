package main

import (
	"fmt"
	"time"

	pb "github.com/brotherlogic/recordwants/proto"
	"golang.org/x/net/context"
)

func (s *Server) alertNoStaging(ctx context.Context, overBudget bool) {
	for _, want := range s.config.Wants {
		if !want.Staged {
			if want.Active {
				s.Log(fmt.Sprintf("Unwanting %v as it's not staged but active", want.Release.Id))
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

			//Superwants don't have to be staged
			if want.Superwant {
				want.Staged = true
			}
		} else {
			if overBudget && want.Active && !want.Superwant {
				s.lastProc = want.Release.Id
				s.Log(fmt.Sprintf("Unwanting %v because we're over budget", want.Release.Id))
				s.recordGetter.unwant(ctx, want)
			} else {
				if want.Active && want.Demoted {
					s.Log(fmt.Sprintf("Unwanting %v because it's demoted", want.Release.Id))
					s.recordGetter.unwant(ctx, want)
				}
				if !want.Active && !want.Demoted && (!overBudget || want.Superwant) {
					s.Log(fmt.Sprintf("WANTING %v -> %v, %v, %v", want.Release.Id, want.Active, want.Demoted, overBudget))
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

	s.save(ctx)
}
