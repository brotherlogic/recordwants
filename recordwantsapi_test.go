package main

import (
	"fmt"
	"testing"

	"golang.org/x/net/context"

	pbrc "github.com/brotherlogic/recordcollection/proto"
	pb "github.com/brotherlogic/recordwants/proto"
)

func InitTestServer() *Server {
	s := Init()
	s.recordGetter = &testRecordGetter{}
	s.SkipLog = true
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
	return []*pbrc.Want{&pbrc.Want{}}, nil
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
