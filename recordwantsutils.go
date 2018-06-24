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
				s.lastProc = want.Release.Id
				s.recordGetter.unwant(ctx, want)
				s.alerter.alert(ctx, want)
			}
		} else {
			if overBudget {
				s.recordGetter.unwant(ctx, want)
			} else {
				if !want.Active {
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
			if w.Release.Id == 3379533 {
				s.pull = fmt.Sprintf("%v", w)
			}
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

	s.save()
}
