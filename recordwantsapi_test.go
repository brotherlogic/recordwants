package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/brotherlogic/keystore/client"
	"golang.org/x/net/context"

	pbgd "github.com/brotherlogic/godiscogs"
	pbrc "github.com/brotherlogic/recordcollection/proto"
	pb "github.com/brotherlogic/recordwants/proto"
)

func InitTestServer() *Server {
	s := Init()
	s.recordGetter = &testRecordGetter{}
	s.recordAdder = &testRecordAdder{}
	s.alerter = &testAlerter{}
	s.SkipLog = true
	s.GoServer.KSclient = *keystoreclient.GetTestClient(".test")
	s.GoServer.KSclient.Save(context.Background(), KEY, &pb.Config{})
	return s
}

type testRecordAdder struct {
	fail bool
}

func (t *testRecordAdder) getAdds(ctx context.Context) ([]int32, error) {
	if t.fail {
		return []int32{}, fmt.Errorf("Built to fail")
	}

	return []int32{123}, nil
}

type testRecordGetter struct {
	fail       bool
	failGet    bool
	lastUnwant int32
	lastWant   int32
}

func (t *testRecordGetter) getRecordsSince(ctx context.Context, since int64) ([]int32, error) {
	if t.fail {
		return []int32{}, fmt.Errorf("Built to fail")
	}
	return []int32{12}, nil
}

func (t *testRecordGetter) getRecords(ctx context.Context, id int32) ([]int32, error) {
	if t.fail {
		return []int32{}, fmt.Errorf("Built to fail")
	}
	return []int32{12}, nil
}

func (t *testRecordGetter) getRecord(ctx context.Context, id int32) (*pbrc.Record, error) {
	if t.failGet {
		return nil, fmt.Errorf("Built to fail")
	}
	return &pbrc.Record{Metadata: &pbrc.ReleaseMetadata{Cost: 100, DateAdded: time.Now().Unix()}}, nil
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
	t.lastUnwant = want.GetRelease().GetId()
	return nil
}

func (t *testRecordGetter) want(ctx context.Context, want *pb.MasterWant) error {
	if t.fail {
		return fmt.Errorf("Built to fail")
	}
	t.lastWant = want.GetRelease().GetId()
	return nil
}

func TestGetSpending(t *testing.T) {
	s := InitTestServer()
	spends, err := s.GetSpending(context.Background(), &pb.SpendingRequest{})

	if err == nil {
		t.Fatalf("Error: %v", spends)
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
	s.Update(context.Background(), &pb.UpdateRequest{Want: &pbgd.Release{Id: 123}, KeepWant: true, Level: pb.MasterWant_LIST})
}

func TestSimpleUpdateFail(t *testing.T) {
	s := InitTestServer()
	s.recordGetter = &testRecordGetter{fail: true}
	s.Update(context.Background(), &pb.UpdateRequest{Want: &pbgd.Release{Id: 123}, KeepWant: true})
}

func TestUpdateEmpty(t *testing.T) {
	s := InitTestServer()
	res, err := s.Update(context.Background(), &pb.UpdateRequest{Want: &pbgd.Release{Id: 1233455}})
	if err == nil {
		t.Errorf("Bad return: %v", res)
	}
}

func TestAddWant(t *testing.T) {
	s := InitTestServer()
	_, err := s.AddWant(context.Background(), &pb.AddWantRequest{ReleaseId: 123, Superwant: true})
	if err != nil {
		t.Errorf("Error adding want: %v", err)
	}
	_, err = s.AddWant(context.Background(), &pb.AddWantRequest{ReleaseId: 123, Superwant: true})
	if err != nil {
		t.Errorf("Error adding want: %v", err)
	}

	_, err = s.GetWant(context.Background(), &pb.GetWantRequest{ReleaseId: 123})
	if err != nil {
		t.Errorf("Fail: %v", err)
	}
}

func TestEmptyGet(t *testing.T) {
	s := InitTestServer()
	want, err := s.GetWant(context.Background(), &pb.GetWantRequest{ReleaseId: 123})
	if err == nil {
		t.Errorf("Fail: %v", want)
	}

}
