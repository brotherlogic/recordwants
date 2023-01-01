package main

import (
	"testing"
	"time"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/recordwants/proto"
)

type testAlerter struct {
	count int
}

func (t *testAlerter) alert(ctx context.Context, want *pb.MasterWant, c, total int) {
	t.count++
}

func TestUpdateWants(t *testing.T) {
	s := InitTestServer()
	//s.config.Wants = append(s.config.Wants, &pb.MasterWant{Release: &pbgd.Release{Id: 123}, Active: true})
	s.updateWants(context.Background(), int32(1234))

}

func TestUnwantActiveWhenDemoted(t *testing.T) {
	s := InitTestServer()
	//s.config.Wants = append(s.config.Wants, &pb.MasterWant{Release: &pbgd.Release{Id: 123}, Active: true, Demoted: true, Staged: true})
	s.alerter = &testAlerter{}
	//s.alertNoStaging(context.Background(), false)

	if s.recordGetter.(*testRecordGetter).lastUnwant == 123 {
		t.Errorf("Want has not been unwanted: %v", 123)
	}
}

func TestKeepWantActiveWhenDemoted(t *testing.T) {
	s := InitTestServer()
	//s.config.Wants = append(s.config.Wants, &pb.MasterWant{Release: &pbgd.Release{Id: 123}, Active: true, Staged: true, Superwant: true})
	s.alerter = &testAlerter{}
	//s.alertNoStaging(context.Background(), false)

	if s.recordGetter.(*testRecordGetter).lastUnwant == 123 {
		t.Errorf("Superwant has been unwanted: %v", 123)
	}
}

func TestUnwantActiveWhenNotDemoted(t *testing.T) {
	s := InitTestServer()
	//s.config.Wants = append(s.config.Wants, &pb.MasterWant{Release: &pbgd.Release{Id: 123}, Active: false, Demoted: false, Staged: true})
	s.alerter = &testAlerter{}
	//s.alertNoStaging(context.Background(), false)

	if s.recordGetter.(*testRecordGetter).lastWant == 123 {
		t.Errorf("Want has not been unwanted: %v", 123)
	}
}

func TestSuperWantPromoted(t *testing.T) {
	s := InitTestServer()
	//s.config.Wants = append(s.config.Wants, &pb.MasterWant{Release: &pbgd.Release{Id: 123}, Superwant: true})
	s.alerter = &testAlerter{}

	//It takes two passes to promote a super want
	//s.alertNoStaging(context.Background(), false)
	//s.alertNoStaging(context.Background(), false)

	if s.recordGetter.(*testRecordGetter).lastWant == 123 {
		t.Errorf("Want has not been unwanted: %v", 123)
	}
}

func TestUpdateWantState(t *testing.T) {
	s := InitTestServer()
	//s.config.Wants = append(s.config.Wants, &pb.MasterWant{})

	err := s.updateWantState(context.Background())

	if err != nil {
		t.Errorf("Bad update: %v", err)
	}
}

func TestUpdateWantBasic(t *testing.T) {
	s := InitTestServer()
	err := s.updateWant(context.Background(), &pb.MasterWant{Level: pb.MasterWant_ALWAYS}, time.Now())

	if err != nil {
		t.Errorf("bad update: %v", err)
	}
}

func TestUpdateWantBasicList(t *testing.T) {
	s := InitTestServer()
	err := s.updateWant(context.Background(), &pb.MasterWant{Level: pb.MasterWant_LIST}, time.Now())

	if err != nil {
		t.Errorf("bad update: %v", err)
	}
}

func TestUpdateWantBasicUNKNOWN(t *testing.T) {
	s := InitTestServer()
	err := s.updateWant(context.Background(), &pb.MasterWant{Level: pb.MasterWant_UNKNOWN, Active: true}, time.Now())

	if err != nil {
		t.Errorf("bad update: %v", err)
	}
}

func TestUpdateWantBasicLISTED(t *testing.T) {
	s := InitTestServer()
	err := s.updateWant(context.Background(), &pb.MasterWant{Level: pb.MasterWant_STAGED_TO_BE_ADDED, Active: true}, time.Now())

	if err != nil {
		t.Errorf("bad update: %v", err)
	}
}

func TestUpdateWantQuick(t *testing.T) {
	s := InitTestServer()
	err := s.updateWant(context.Background(), &pb.MasterWant{Level: pb.MasterWant_ALWAYS, Dirty: true}, time.Now())

	if err != nil {
		t.Errorf("bad update: %v", err)
	}
}

func TestUpdateFailRead(t *testing.T) {
	s := InitTestServer()
	s.GoServer.KSclient.Fail = true

	err := s.updateWantState(context.Background())

	if err == nil {
		t.Errorf("Bad load no fail")
	}
}

func TestUpdateWantStateANYTIMEUp(t *testing.T) {
	s := InitTestServer()
	//s.config.Budget = 10
	err := s.updateWant(context.Background(), &pb.MasterWant{Level: pb.MasterWant_ANYTIME, Active: false}, time.Now())

	if err != nil {
		t.Errorf("Bad update: %v", err)
	}
}

func TestUpdateWantBasicANYTIMEDown(t *testing.T) {
	s := InitTestServer()
	//s.config.Budget = -10
	err := s.updateWant(context.Background(), &pb.MasterWant{Level: pb.MasterWant_ANYTIME, Active: true}, time.Now())

	if err != nil {
		t.Errorf("bad update: %v", err)
	}
}
