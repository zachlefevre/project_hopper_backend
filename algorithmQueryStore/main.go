package main

import (
	"context"
	"encoding/json"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/zachlefevre/project_hopper_backend/com"

	"github.com/nats-io/go-nats-streaming"
	"github.com/zachlefevre/project_hopper_backend/natsutil"
	"google.golang.org/grpc"
)

const (
	baseClient_id  = "algorithm-repository"
	cluster_id     = "test-cluster"
	create_channel = "create-algorithm"
	grpcUri        = "store:50051"
	added_event    = "algorithm-added-to-query-store"
	aggregate      = "Alorithm"
	durableID      = "algorithm-repository-durable"
	natsURL        = "nats:4222"
	queryStoreURI  = "queryStore:50051"
)

func main() {
	cid, _ := uuid.NewV4()
	clientID := baseClient_id + cid.String()
	comp := natsutil.NewStreamingComponent(clientID)

	err := comp.ConnectToNatsStreamingService(cluster_id, stan.NatsURL(natsURL), stan.ConnectWait(100*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	sc := comp.NATS()

	var wg sync.WaitGroup

	aw, _ := time.ParseDuration("60s")
	wg.Add(1)
	go func() {
		sc.Subscribe(create_channel, func(msg *stan.Msg) {
			msg.Ack()
			createCmd := pb.CreateAlgorithmCommand{}
			err := json.Unmarshal(msg.Data, &createCmd)
			if err != nil {
				log.Print(err)
				return
			}

			persistToQueryStore(&createCmd)

			eventID, _ := uuid.NewV4()
			createdEvent := pb.AlgorithmCreatedEvent{
				Algorithm: createCmd.Algorithm,
				Id:        eventID.String(),
			}
			if err := createAlgorithmCreatedEvent(createdEvent); err != nil {

			}
		}, stan.DurableName(durableID),
			stan.MaxInflight(25),
			stan.SetManualAckMode(),
			stan.AckWait(aw),
		)
		wg.Done()
	}()

	wg.Wait()
	runtime.Goexit()
}

func createAlgorithmCreatedEvent(createdEvent pb.AlgorithmCreatedEvent) error {
	conn, err := grpc.Dial(grpcUri, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewEventStoreClient(conn)
	createdEventJSON, _ := json.Marshal(createdEvent)
	eid, _ := uuid.NewV4()
	event := &pb.Event{
		EventId:       eid.String(),
		EventType:     added_event,
		AggregateId:   createdEvent.Id,
		AggregateType: aggregate,
		EventData:     string(createdEventJSON),
		Channel:       added_event,
	}

	resp, err := client.CreateEvent(context.Background(), event)
	if err != nil {
		return errors.Wrap(err, "errors from RPC server")
	}
	if resp.IsSuccessful {
		return nil
	} else {
		return errors.Wrap(err, "errors from RPC server")
	}
}
