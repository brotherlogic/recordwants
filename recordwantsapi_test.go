package main

import (
	"fmt"
	"testing"

	"github.com/brotherlogic/keystore/client"
	"golang.org/x/net/context"

	pbgd "github.com/brotherlogic/godiscogs"
	pbrc "github.com/brotherlogic/recordcollection/proto"
	pb "github.com/brotherlogic/recordwants/proto"
)

func InitTestServer() *Server {
	s := Init()
	s.recordGetter = &testRecordGetter{}
	s.alerter = &testAlerter{}
	s.SkipLog = true
	s.GoServer.KSclient = *keystoreclient.GetTestClient(".test")
	return s
}

type testRecordGetter struct {
	fail bool
}

func (t *testRecordGetter) getRecords(ctx context.Context) ([]*pbrc.Record, error) {
	if t.fail {
		return make([]*pbrc.Record, 0), fmt.Errorf("Built to fail")
	}
	return []*pbrc.Record{&pbrc.Record{Metadata: &pbrc.ReleaseMetadata{DateAdded: 1515888000, Cost: 1234}}}, nil
}

func (t *testRecordGetter) getWants(ctx context.Context) ([]*pbrc.Want, error) {
	if t.fail {
		return make([]*pbrc.Want, 0), fmt.Errorf("Built to fail")
	}
	return []*pbrc.Want{
		&pbrc.Want{Release: &pbgd.Release{Id: 123}, Metadata: &pbrc.WantMetadata{Active: true}},
		&pbrc.Want{Release: &pbgd.Release{Id: 124}, Metadata: &pbrc.WantMetadata{Active: true}},
	}, nil
}

func (t *testRecordGetter) unwant(ctx context.Context, want *pb.MasterWant) error {
	if t.fail {
		return fmt.Errorf("Built to fail")
	}
	return nil
}

func (t *testRecordGetter) want(ctx context.Context, want *pb.MasterWant) error {
	if t.fail {
		return fmt.Errorf("Built to fail")
	}
	return nil
}

func TestGetSpending(t *testing.T) {
	s := InitTestServer()
	spends, err := s.GetSpending(context.Background(), &pb.SpendingRequest{})

	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if len(spends.Spends) != 12 || spends.Spends[0].Spend != 1234 {
		t.Errorf("Error in spend amount! %v", spends.Spends)
	}
}

func TestGetSpendingFail(t *testing.T) {
	s := InitTestServer()
	s.recordGetter = &testRecordGetter{fail: true}
	_, err := s.GetSpending(context.Background(), &pb.SpendingRequest{})

	if err == nil {
		t.Errorf("Did not fail")
	}
}

func TestSimpleUpdate(t *testing.T) {
	s := InitTestServer()
	s.config.Wants = append(s.config.Wants, &pb.MasterWant{Release: &pbgd.Release{Id: 123}})
	s.Update(context.Background(), &pb.UpdateRequest{Want: &pbgd.Release{Id: 123}, KeepWant: true})
}

func TestSimpleUpdateFail(t *testing.T) {
	s := InitTestServer()
	s.recordGetter = &testRecordGetter{fail: true}
	s.config.Wants = append(s.config.Wants, &pb.MasterWant{Release: &pbgd.Release{Id: 123}})
	s.Update(context.Background(), &pb.UpdateRequest{Want: &pbgd.Release{Id: 123}, KeepWant: true})
}
