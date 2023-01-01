package main

import (
	"testing"
	"time"

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

func TestUpdate(t *testing.T) {
	s := InitTestServer()

	//s.config.Wants = append(s.config.Wants, &pb.MasterWant{Release: &pbgd.Release{Id: 123}, Superwant: true})

	err := s.dealWithAddedRecords(context.Background())

	if err != nil {
		t.Errorf("Bad add: %v", err)
	}
}

func TestBadUpdate(t *testing.T) {
	s := InitTestServer()
	s.recordAdder = &testRecordAdder{fail: true}

	//s.config.Wants = append(s.config.Wants, &pb.MasterWant{Release: &pbgd.Release{Id: 123}, Superwant: true})

	err := s.dealWithAddedRecords(context.Background())

	if err == nil {
		t.Errorf("Bad add: %v", err)
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

func TestUpdateWantStateFail(t *testing.T) {
	s := InitTestServer()
	s.recordGetter = &testRecordGetter{fail: true}
	config := &pb.Config{}
	config.Wants = append(config.Wants, &pb.MasterWant{Level: pb.MasterWant_ALWAYS})
	s.save(context.Background(), config)

	err := s.updateWantState(context.Background())

	if err == nil {
		t.Errorf("Bad update: %v", err)
	}
}

func TestUpdateWantBasic(t *testing.T) {
	s := InitTestServer()
	err := s.updateWant(context.Background(), &pb.MasterWant{Level: pb.MasterWant_ALWAYS}, 0, time.Now())

	if err != nil {
		t.Errorf("bad update: %v", err)
	}
}

func TestUpdateWantBasicList(t *testing.T) {
	s := InitTestServer()
	err := s.updateWant(context.Background(), &pb.MasterWant{Level: pb.MasterWant_LIST}, 0, time.Now())

	if err != nil {
		t.Errorf("bad update: %v", err)
	}
}

func TestUpdateWantStateFailLISTTED(t *testing.T) {
	s := InitTestServer()
	s.recordGetter = &testRecordGetter{fail: true}
	err := s.updateWant(context.Background(), &pb.MasterWant{Level: pb.MasterWant_UNKNOWN, Active: true}, 0, time.Now())

	if err == nil {
		t.Errorf("Bad update: %v", err)
	}
}

func TestUpdateWantBasicUNKNOWN(t *testing.T) {
	s := InitTestServer()
	err := s.updateWant(context.Background(), &pb.MasterWant{Level: pb.MasterWant_UNKNOWN, Active: true}, 0, time.Now())

	if err != nil {
		t.Errorf("bad update: %v", err)
	}
}

func TestUpdateWantStateFailUNKNOWN(t *testing.T) {
	s := InitTestServer()
	s.recordGetter = &testRecordGetter{fail: true}
	err := s.updateWant(context.Background(), &pb.MasterWant{Level: pb.MasterWant_STAGED_TO_BE_ADDED, Active: true}, 0, time.Now())

	if err == nil {
		t.Errorf("Bad update: %v", err)
	}
}

func TestUpdateWantBasicLISTED(t *testing.T) {
	s := InitTestServer()
	err := s.updateWant(context.Background(), &pb.MasterWant{Level: pb.MasterWant_STAGED_TO_BE_ADDED, Active: true}, 0, time.Now())

	if err != nil {
		t.Errorf("bad update: %v", err)
	}
}

func TestUpdateWantQuick(t *testing.T) {
	s := InitTestServer()
	err := s.updateWant(context.Background(), &pb.MasterWant{Level: pb.MasterWant_ALWAYS, Dirty: true}, 0, time.Now())

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
	err := s.updateWant(context.Background(), &pb.MasterWant{Level: pb.MasterWant_ANYTIME, Active: false}, 10, time.Now())

	if err != nil {
		t.Errorf("Bad update: %v", err)
	}
}

func TestUpdateWantBasicANYTIMEDown(t *testing.T) {
	s := InitTestServer()
	//s.config.Budget = -10
	err := s.updateWant(context.Background(), &pb.MasterWant{Level: pb.MasterWant_ANYTIME, Active: true}, -10, time.Now())

	if err != nil {
		t.Errorf("bad update: %v", err)
	}
}

func TestUpdateWantStateANYTIMEUpFail(t *testing.T) {
	s := InitTestServer()
	s.recordGetter = &testRecordGetter{fail: true}
	//s.config.Budget = 10
	err := s.updateWant(context.Background(), &pb.MasterWant{Level: pb.MasterWant_ANYTIME, Active: false}, 10, time.Now())

	if err == nil {
		t.Errorf("Bad update: %v", err)
	}
}

func TestUpdateWantBasicANYTIMEDownFail(t *testing.T) {
	s := InitTestServer()
	s.recordGetter = &testRecordGetter{fail: true}
	//s.config.Budget = -10
	err := s.updateWant(context.Background(), &pb.MasterWant{Level: pb.MasterWant_ANYTIME, Active: true}, -10, time.Now())

	if err == nil {
		t.Errorf("bad update: %v", err)
	}
}

func TestUpdateWantCategory(t *testing.T) {
	s := InitTestServer()

	want := &pb.MasterWant{Release: &pbgd.Release{Id: 123}, Level: pb.MasterWant_ANYTIME, RetireTime: time.Now().Add(-time.Minute).Unix(), RetireLevel: pb.MasterWant_NEVER}
	err := s.updateWant(context.Background(), want, int32(10), time.Now())

	if err != nil {
		t.Errorf("Bad update: %v", err)
	}

	if want.GetLevel() == pb.MasterWant_ANYTIME {
		t.Errorf("The level hasn't change: %v", want)
	}
}
