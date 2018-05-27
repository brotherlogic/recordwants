package main

import (
	"testing"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/recordwants/proto"
)

type testAlerter struct {
	count int
}

func (t *testAlerter) alert(ctx context.Context, want *pb.MasterWant) {
	t.count++
}

func TestMainTest(t *testing.T) {
	s := InitTestServer()
	s.config.Wants = append(s.config.Wants, &pb.MasterWant{Staged: false})
	ta := &testAlerter{}
	s.alerter = ta
	s.alertNoStaging(context.Background())

	if ta.count != 1 {
		t.Errorf("Not enough alerts!")
	}
}
