package main

import "context"

func (s *Server) alertNoStaging(ctx context.Context, overBudget bool) {
	for _, want := range s.config.Wants {
		if !want.Staged {
			s.recordGetter.unwant(ctx, want)
			s.alerter.alert(ctx, want)
		} else {
			if overBudget {
				s.recordGetter.unwant(ctx, want)
			}
		}
	}
}
