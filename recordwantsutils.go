package main

import "context"

func (s *Server) alertNoStaging(ctx context.Context) {
	for _, want := range s.config.Wants {
		if !want.Staged {
			s.alerter.alert(ctx, want)
		}
	}
}
