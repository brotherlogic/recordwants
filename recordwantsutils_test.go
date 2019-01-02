package main

import (
	"testing"

	"golang.org/x/net/context"

	pbgd "github.com/brotherlogic/godiscogs"
	pb "github.com/brotherlogic/recordwants/proto"
)

type testAlerter struct {
	count int
}

func (t *testAlerter) alert(ctx context.Context, want *pb.MasterWant, c, total int) {
	t.count++
}

func TestMainTest(t *testing.T) {
	s := InitTestServer()
	s.config.Wants = append(s.config.Wants, &pb.MasterWant{Staged: false, Release: &pbgd.Release{Id: 123}, Active: true})
	s.config.Wants = append(s.config.Wants, &pb.MasterWant{Staged: true, Release: &pbgd.Release{Id: 123}, Active: true})
	ta := &testAlerter{}
	s.alerter = ta
	s.alertNoStaging(context.Background(), false)

	if ta.count != 1 {
		t.Errorf("Not enough alerts!")
	}
}

func TestMainTestOverBudget(t *testing.T) {
	s := InitTestServer()
	s.config.Wants = append(s.config.Wants, &pb.MasterWant{Staged: true, Active: true, Release: &pbgd.Release{Id: 123}})
	ta := &testAlerter{}
	s.alerter = ta
	s.alertNoStaging(context.Background(), true)

	if ta.count != 0 {
		t.Errorf("Not enough alerts!")
	}
}

func TestUpdateWants(t *testing.T) {
	s := InitTestServer()
	s.config.Wants = append(s.config.Wants, &pb.MasterWant{Release: &pbgd.Release{Id: 123}, Active: true})
	s.updateWants(context.Background())

	if len(s.config.Wants) != 2 {
		t.Errorf("No wants added!")
	}
}

func TestUnwantActiveWhenDemoted(t *testing.T) {
	s := InitTestServer()
	s.config.Wants = append(s.config.Wants, &pb.MasterWant{Release: &pbgd.Release{Id: 123}, Active: true, Demoted: true, Staged: true})
	s.alerter = &testAlerter{}
	s.alertNoStaging(context.Background(), false)

	if s.recordGetter.(*testRecordGetter).lastUnwant != 123 {
		t.Errorf("Want has not been unwanted: %v", s.config)
	}
}

func TestKeepWantActiveWhenDemoted(t *testing.T) {
	s := InitTestServer()
	s.config.Wants = append(s.config.Wants, &pb.MasterWant{Release: &pbgd.Release{Id: 123}, Active: true, Staged: true, Superwant: true})
	s.alerter = &testAlerter{}
	s.alertNoStaging(context.Background(), false)

	if s.recordGetter.(*testRecordGetter).lastUnwant == 123 {
		t.Errorf("Superwant has been unwanted: %v", s.config)
	}
}

func TestUnwantActiveWhenNotDemoted(t *testing.T) {
	s := InitTestServer()
	s.config.Wants = append(s.config.Wants, &pb.MasterWant{Release: &pbgd.Release{Id: 123}, Active: false, Demoted: false, Staged: true})
	s.alerter = &testAlerter{}
	s.alertNoStaging(context.Background(), false)

	if s.recordGetter.(*testRecordGetter).lastWant != 123 {
		t.Errorf("Want has not been unwanted: %v", s.config)
	}
}

func TestSuperWantPromoted(t *testing.T) {
	s := InitTestServer()
	s.config.Wants = append(s.config.Wants, &pb.MasterWant{Release: &pbgd.Release{Id: 123}, Superwant: true})
	s.alerter = &testAlerter{}

	//It takes two passes to promote a super want
	s.alertNoStaging(context.Background(), false)
	s.alertNoStaging(context.Background(), false)

	if s.recordGetter.(*testRecordGetter).lastWant != 123 {
		t.Errorf("Want has not been unwanted: %v", s.config)
	}
}
