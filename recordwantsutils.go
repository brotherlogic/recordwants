package main

import (
	"log"
	"time"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/recordwants/proto"
)

func (s *Server) alertNoStaging(ctx context.Context, overBudget bool) {
	for _, want := range s.config.Wants {
		if !want.Staged {
			err := s.recordGetter.unwant(ctx, want)
			if err != nil {
				log.Fatalf("Bah: %v", err)
			}
			s.alerter.alert(ctx, want)
		} else {
			if overBudget {
				s.recordGetter.unwant(ctx, want)
			}
		}
	}
}

func (s *Server) updateWants(ctx context.Context) {
	s.lastRun = time.Now()
	wants, err := s.recordGetter.getWants(ctx)
	if err == nil {
		for _, w := range wants {
			found := false
			for _, mw := range s.config.Wants {
				if mw.Release.Id == w.Release.Id {
					found = true
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
