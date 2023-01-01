package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/brotherlogic/goserver/utils"

	pbgd "github.com/brotherlogic/godiscogs"
	rcpb "github.com/brotherlogic/recordcollection/proto"
	pb "github.com/brotherlogic/recordwants/proto"

	//Needed to pull in gzip encoding init
	_ "google.golang.org/grpc/encoding/gzip"
)

func main() {
	ctx, cancel := utils.BuildContext("recordwants-cli", "recordwants")
	defer cancel()

	conn, err := utils.LFDialServer(ctx, "recordwants")

	if err != nil {
		log.Fatalf("Unable to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewWantServiceClient(conn)

	switch os.Args[1] {
	case "ping":
		client := rcpb.NewClientUpdateServiceClient(conn)
		id, _ := strconv.Atoi(os.Args[2])
		_, err := client.ClientUpdate(ctx, &rcpb.ClientUpdateRequest{InstanceId: int32(id)})
		fmt.Printf("Ping: %v", err)
	case "want":
		iv, _ := strconv.Atoi(os.Args[2])
		w, err := client.Update(ctx, &pb.UpdateRequest{Want: &pbgd.Release{Id: int32(iv)}, Level: pb.MasterWant_WANT_DIGITAL})
		if err != nil {
			log.Fatalf("Error on GET: %v", err)
		}
		fmt.Printf("%v\n", w)
	case "list":
		iv, _ := strconv.Atoi(os.Args[2])
		_, err := client.Update(ctx, &pb.UpdateRequest{Want: &pbgd.Release{Id: int32(iv)}, Level: pb.MasterWant_LIST})
		if err != nil {
			log.Fatalf("Error on GET: %v", err)
		}
	case "all":
		wa, err := client.GetWants(ctx, &pb.GetWantsRequest{})
		if err != nil {
			log.Fatalf("Error on GET: %v", err)
		}
		for i, w := range wa.GetWant() {
			fmt.Printf("%v. %v\n", i, w)
		}
	case "pingall":
		wa, err := client.GetWants(ctx, &pb.GetWantsRequest{})
		if err != nil {
			log.Fatalf("Error on GET: %v", err)
		}

		conn2, err2 := utils.LFDialServer(ctx, "recordcollection")
		if err2 != nil {
			log.Fatalf("Pah: %v", err2)
		}
		defer conn2.Close()
		rcclient := rcpb.NewRecordCollectionServiceClient(conn2)
		rclient := rcpb.NewClientUpdateServiceClient(conn)

		for i, w := range wa.GetWant() {
			if w.GetLevel() == pb.MasterWant_ANYTIME_LIST || w.GetLevel() == pb.MasterWant_LIST {
				fmt.Printf("%v. %v\n", i, w)
				recs, err := rcclient.QueryRecords(ctx, &rcpb.QueryRecordsRequest{Query: &rcpb.QueryRecordsRequest_ReleaseId{int32(w.Release.Id)}})
				if err != nil {
					log.Fatalf("Bad get: %v", err)
				}
				for _, r := range recs.GetInstanceIds() {
					_, err := rclient.ClientUpdate(ctx, &rcpb.ClientUpdateRequest{InstanceId: r})
					fmt.Printf("%v\n", err)
				}
			}
		}
	case "random":
		wa, err := client.GetWants(ctx, &pb.GetWantsRequest{})
		if err != nil {
			log.Fatalf("Error on GET: %v", err)
		}

		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(wa.GetWant()), func(i, j int) { wa.Want[i], wa.Want[j] = wa.Want[j], wa.Want[i] })

		count := 0
		for i, want := range wa.GetWant() {
			if want.GetLevel() == pb.MasterWant_WANT_OG || want.GetLevel() == pb.MasterWant_WANT_DIGITAL {
				fmt.Printf("%v. %v\n", i, want)
				count++
			}

			if count > 5 {
				return
			}
		}
	case "sync":
		wa, err := client.Sync(ctx, &pb.SyncRequest{})
		fmt.Printf("%v and %v\n", wa, err)
	case "next":
		count := 0

		if len(os.Args) > 2 {
			v, err := strconv.Atoi(os.Args[2])
			if err != nil {
				log.Fatalf("Conversion: %v", err)
			}
			count = v
		}

		wa, err := client.GetWants(ctx, &pb.GetWantsRequest{})
		if err != nil {
			log.Fatalf("Error on GET: %v", err)
		}
		for i, w := range wa.GetWant() {
			if w.GetLevel() == pb.MasterWant_UNKNOWN {
				fmt.Printf("%v. %v = %v\n", i, w.Release.GetId(), w)
				if count == 0 {
					return
				}
				count--
			}
		}
	case "clearall":
		wa, err := client.GetWants(ctx, &pb.GetWantsRequest{})
		if err != nil {
			log.Fatalf("Error on GET: %v", err)
		}
		for i, w := range wa.GetWant() {
			if w.GetLevel() == pb.MasterWant_ANYTIME {
				_, err := client.Update(ctx, &pb.UpdateRequest{Want: w.GetRelease(), Level: pb.MasterWant_NEVER})
				if err != nil {
					log.Fatalf("ERROR ON UPDATE: %v", err)
				}
				fmt.Printf("%v. %v\n", i, w)
			}
		}
	case "get":
		iv, _ := strconv.Atoi(os.Args[2])
		wa, err := client.GetWants(ctx, &pb.GetWantsRequest{ReleaseId: []int32{int32(iv)}})
		if err != nil {
			log.Fatalf("Error on GET: %v", err)
		}
		fmt.Printf("GOT: %v\n", wa.GetWant())
	case "super":
		iv, _ := strconv.Atoi(os.Args[2])
		_, err := client.Update(ctx, &pb.UpdateRequest{Want: &pbgd.Release{Id: int32(iv)}, Super: true})
		if err != nil {
			log.Fatalf("Error on GET: %v", err)
		}
	case "always":
		iv, _ := strconv.Atoi(os.Args[2])
		_, err := client.Update(ctx, &pb.UpdateRequest{Want: &pbgd.Release{Id: int32(iv)}, Level: pb.MasterWant_ALWAYS})
		if err != nil {
			log.Fatalf("Error on GET: %v", err)
		}
	case "never":
		iv, _ := strconv.Atoi(os.Args[2])
		w, err := client.Update(ctx, &pb.UpdateRequest{Want: &pbgd.Release{Id: int32(iv)}, Level: pb.MasterWant_NEVER})
		if err != nil {
			log.Fatalf("Error on GET: %v", err)
		}
		fmt.Printf("%v\n", w)
	}
}
