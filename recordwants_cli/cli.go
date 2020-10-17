package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

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
		_, err := client.ClientUpdate(ctx, &rcpb.ClientUpdateRequest{})
		fmt.Printf("Ping: %v", err)
	case "spend":
		res, err := client.GetSpending(ctx, &pb.SpendingRequest{})
		if err != nil {
			log.Fatalf("Error on GET: %v", err)
		}
		total := int32(0)
		for _, v := range res.Spends {
			fmt.Printf("%v = %v\n", v.Month, v.Spend)
			total += v.Spend
		}
		fmt.Printf("TOTAL = %v\n", total)
	case "want":
		iv, _ := strconv.Atoi(os.Args[2])
		w, err := client.Update(ctx, &pb.UpdateRequest{Want: &pbgd.Release{Id: int32(iv)}, Level: pb.MasterWant_ANYTIME})
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
	case "unwant":
		iv, _ := strconv.Atoi(os.Args[2])
		w, err := client.Update(ctx, &pb.UpdateRequest{Want: &pbgd.Release{Id: int32(iv)}, Level: pb.MasterWant_NEVER})
		if err != nil {
			log.Fatalf("Error on GET: %v", err)
		}
		fmt.Printf("%v\n", w)
	}
}
