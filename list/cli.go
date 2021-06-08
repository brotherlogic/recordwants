package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/brotherlogic/goserver/utils"

	pbgd "github.com/brotherlogic/godiscogs"
	rcpb "github.com/brotherlogic/recordcollection/proto"
	pb "github.com/brotherlogic/recordwants/proto"

	//Needed to pull in gzip encoding init
	"google.golang.org/grpc/codes"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/status"
)

func main() {
	ctx, cancel := utils.ManualContext("recordwants-cli", time.Minute*10)
	defer cancel()

	conn, err := utils.LFDialServer(ctx, "recordwants")
	if err != nil {
		log.Fatalf("Unable to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewWantServiceClient(conn)

	conn2, err := utils.LFDialServer(ctx, "recordcollection")
	if err != nil {
		log.Fatalf("Unable to dial: %v", err)
	}
	defer conn2.Close()
	rcclient := rcpb.NewRecordCollectionServiceClient(conn2)

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	dateStr := scanner.Text()
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Fatalf("Error parsing date in %v -> %v", os.Args[1], err)
	}

	for scanner.Scan() {
		id, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatalf("Error parsing text: %v -> %v", err, scanner.Text())
		}

		// Check we doon't already own this
		ids, err := rcclient.QueryRecords(ctx, &rcpb.QueryRecordsRequest{Query: &rcpb.QueryRecordsRequest_ReleaseId{int32(id)}})
		if err != nil {
			log.Fatalf("Error getting record: %v", err)
		}

		if len(ids.GetInstanceIds()) == 0 {
			_, err = client.AddWant(ctx, &pb.AddWantRequest{
				ReleaseId:   int32(id),
				Level:       pb.MasterWant_LIST,
				RetireTime:  date.Unix(),
				RetireLevel: pb.MasterWant_NEVER,
			})

			if status.Convert(err).Code() == codes.FailedPrecondition {
				_, err = client.Update(ctx, &pb.UpdateRequest{
					Want:        &pbgd.Release{Id: int32(id)},
					Level:       pb.MasterWant_LIST,
					RetireTime:  date.Unix(),
					RetireLevel: pb.MasterWant_NEVER,
				})
			}

			if err != nil {
				log.Fatalf("Error adding want: %v", err)
			}
		}
	}
}
