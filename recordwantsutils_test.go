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

	if ta.count != 3 {
		t.Errorf("Not enough alerts!: %v", ta.count)
	}
}

func TestMainTestOverBudget(t *testing.T) {
	s := InitTestServer()
	s.config.Wants = append(s.config.Wants, &pb.MasterWant{Staged: true, Active: true, Release: &pbgd.Release{Id: 123}})
	ta := &testAlerter{}
	s.alerter = ta
	s.alertNoStaging(context.Background(), true)

	if ta.count != 1 {
		t.Errorf("Not enough alerts!: %v", ta.count)
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

	if s.recordGetter.(*testRecordGetter).lastUnwant == 123 {
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

	if s.recordGetter.(*testRecordGetter).lastWant == 123 {
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

	if s.recordGetter.(*testRecordGetter).lastWant == 123 {
		t.Errorf("Want has not been unwanted: %v", s.config)
	}
}

func TestUpdateSpending(t *testing.T) {
	s := InitTestServer()

	err := s.updateSpending(context.Background())

	if err != nil {
		t.Errorf("Bad update: %v", err)
	}
}

func TestUpdateSpendingFailSince(t *testing.T) {
	s := InitTestServer()
	s.recordGetter = &testRecordGetter{fail: true}

	err := s.updateSpending(context.Background())

	if err == nil {
		t.Errorf("Bad update: %v", err)
	}
}

func TestUpdateSpendingFailGet(t *testing.T) {
	s := InitTestServer()
	s.recordGetter = &testRecordGetter{failGet: true}

	err := s.updateSpending(context.Background())

	if err == nil {
		t.Errorf("Bad update: %v", err)
	}
}

func TestUpdate(t *testing.T) {
	s := InitTestServer()

	s.config.Wants = append(s.config.Wants, &pb.MasterWant{Release: &pbgd.Release{Id: 123}, Superwant: true})

	err := s.dealWithAddedRecords(context.Background())

	if err != nil {
		t.Errorf("Bad add: %v", err)
	}
}

func TestBadUpdate(t *testing.T) {
	s := InitTestServer()
	s.recordAdder = &testRecordAdder{fail: true}

	s.config.Wants = append(s.config.Wants, &pb.MasterWant{Release: &pbgd.Release{Id: 123}, Superwant: true})

	err := s.dealWithAddedRecords(context.Background())

	if err == nil {
		t.Errorf("Bad add: %v", err)
	}
}

func TestUpdateWantState(t *testing.T) {
	s := InitTestServer()
	s.config.Wants = append(s.config.Wants, &pb.MasterWant{})

	_, err := s.updateWantState(context.Background())

	if err != nil {
		t.Errorf("Bad update: %v", err)
	}
}

func TestUpdateWantStateFail(t *testing.T) {
	s := InitTestServer()
	s.recordGetter = &testRecordGetter{fail: true}
	s.config.Wants = append(s.config.Wants, &pb.MasterWant{Level: pb.MasterWant_ALWAYS})
	s.save(context.Background())

	_, err := s.updateWantState(context.Background())

	if err == nil {
		t.Errorf("Bad update: %v", err)
	}
}

func TestUpdateWantBasic(t *testing.T) {
	s := InitTestServer()
	err := s.updateWant(context.Background(), &pb.MasterWant{Level: pb.MasterWant_ALWAYS})

	if err != nil {
		t.Errorf("bad update: %v", err)
	}
}

func TestUpdateWantStateFailList(t *testing.T) {
	s := InitTestServer()
	s.recordGetter = &testRecordGetter{fail: true}
	s.config.Wants = append(s.config.Wants, &pb.MasterWant{Level: pb.MasterWant_LIST})
	s.save(context.Background())

	_, err := s.updateWantState(context.Background())

	if err == nil {
		t.Errorf("Bad update: %v", err)
	}
}

func TestUpdateWantBasicList(t *testing.T) {
	s := InitTestServer()
	err := s.updateWant(context.Background(), &pb.MasterWant{Level: pb.MasterWant_LIST})

	if err != nil {
		t.Errorf("bad update: %v", err)
	}
}

func TestUpdateWantStateFailLISTTED(t *testing.T) {
	s := InitTestServer()
	s.recordGetter = &testRecordGetter{fail: true}
	err := s.updateWant(context.Background(), &pb.MasterWant{Level: pb.MasterWant_UNKNOWN, Active: true})

	if err == nil {
		t.Errorf("Bad update: %v", err)
	}
}

func TestUpdateWantBasicUNKNOWN(t *testing.T) {
	s := InitTestServer()
	err := s.updateWant(context.Background(), &pb.MasterWant{Level: pb.MasterWant_UNKNOWN, Active: true})

	if err != nil {
		t.Errorf("bad update: %v", err)
	}
}

func TestUpdateWantStateFailUNKNOWN(t *testing.T) {
	s := InitTestServer()
	s.recordGetter = &testRecordGetter{fail: true}
	err := s.updateWant(context.Background(), &pb.MasterWant{Level: pb.MasterWant_STAGED_TO_BE_ADDED, Active: true})

	if err == nil {
		t.Errorf("Bad update: %v", err)
	}
}

func TestUpdateWantBasicLISTED(t *testing.T) {
	s := InitTestServer()
	err := s.updateWant(context.Background(), &pb.MasterWant{Level: pb.MasterWant_STAGED_TO_BE_ADDED, Active: true})

	if err != nil {
		t.Errorf("bad update: %v", err)
	}
}

func TestUpdateWantQuick(t *testing.T) {
	s := InitTestServer()
	err := s.updateWant(context.Background(), &pb.MasterWant{Level: pb.MasterWant_ALWAYS, Dirty: true})

	if err != nil {
		t.Errorf("bad update: %v", err)
	}
}

func TestUpdateFailRead(t *testing.T) {
	s := InitTestServer()
	s.GoServer.KSclient.Fail = true

	_, err := s.updateWantState(context.Background())

	if err == nil {
		t.Errorf("Bad load no fail")
	}
}

func TestUpdateWantStateANYTIMEUp(t *testing.T) {
	s := InitTestServer()
	s.config.Budget = 10
	err := s.updateWant(context.Background(), &pb.MasterWant{Level: pb.MasterWant_ANYTIME, Active: false})

	if err != nil {
		t.Errorf("Bad update: %v", err)
	}
}

func TestUpdateWantBasicANYTIMEDown(t *testing.T) {
	s := InitTestServer()
	s.config.Budget = -10
	err := s.updateWant(context.Background(), &pb.MasterWant{Level: pb.MasterWant_ANYTIME, Active: true})

	if err != nil {
		t.Errorf("bad update: %v", err)
	}
}

func TestUpdateWantStateANYTIMEUpFail(t *testing.T) {
	s := InitTestServer()
	s.recordGetter = &testRecordGetter{fail: true}
	s.config.Budget = 10
	err := s.updateWant(context.Background(), &pb.MasterWant{Level: pb.MasterWant_ANYTIME, Active: false})

	if err == nil {
		t.Errorf("Bad update: %v", err)
	}
}

func TestUpdateWantBasicANYTIMEDownFail(t *testing.T) {
	s := InitTestServer()
	s.recordGetter = &testRecordGetter{fail: true}
	s.config.Budget = -10
	err := s.updateWant(context.Background(), &pb.MasterWant{Level: pb.MasterWant_ANYTIME, Active: true})

	if err == nil {
		t.Errorf("bad update: %v", err)
	}
}
