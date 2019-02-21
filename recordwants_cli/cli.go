package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/brotherlogic/goserver/utils"
	"google.golang.org/grpc"

	pbgd "github.com/brotherlogic/godiscogs"
	pb "github.com/brotherlogic/recordwants/proto"

	//Needed to pull in gzip encoding init
	_ "google.golang.org/grpc/encoding/gzip"
)

func main() {
	host, port, err := utils.Resolve("recordwants")
	if err != nil {
		log.Fatalf("Unable to reach organiser: %v", err)
	}
	conn, err := grpc.Dial(host+":"+strconv.Itoa(int(port)), grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		log.Fatalf("Unable to dial: %v", err)
	}

	client := pb.NewWantServiceClient(conn)
	ctx, cancel := utils.BuildContext("recordwants-cli", "recordwants")
	defer cancel()

	switch os.Args[1] {
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
		_, err := client.Update(ctx, &pb.UpdateRequest{Want: &pbgd.Release{Id: int32(iv)}, KeepWant: true})
		if err != nil {
			log.Fatalf("Error on GET: %v", err)
		}
	case "super":
		iv, _ := strconv.Atoi(os.Args[2])
		_, err := client.Update(ctx, &pb.UpdateRequest{Want: &pbgd.Release{Id: int32(iv)}, Super: true})
		if err != nil {
			log.Fatalf("Error on GET: %v", err)
		}

	case "unwant":
		iv, _ := strconv.Atoi(os.Args[2])
		_, err := client.Update(ctx, &pb.UpdateRequest{Want: &pbgd.Release{Id: int32(iv)}})
		if err != nil {
			log.Fatalf("Error on GET: %v", err)
		}

	}
}
