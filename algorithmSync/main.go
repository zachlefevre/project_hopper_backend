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
	baseClientID  = "algorithm-repository"
	clusterID     = "test-cluster"
	createChannel = "create-algorithm"
	grpcURI       = "store:50051"
	addedEvent    = "algorithm-added-to-query-store"
	aggregate     = "Alorithm"
	durableID     = "algorithm-repository-durable"
	natsURL       = "nats:4222"
	queryStoreURI = "querystore:50051"
)

func main() {
	cid, _ := uuid.NewV4()
	clientID := baseClientID + cid.String()
	comp := natsutil.NewStreamingComponent(clientID)

	err := comp.ConnectToNatsStreamingService(clusterID, stan.NatsURL(natsURL), stan.ConnectWait(100*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	sc := comp.NATS()

	var wg sync.WaitGroup

	aw, _ := time.ParseDuration("60s")
	wg.Add(1)
	go func() {
		sc.Subscribe(createChannel, func(msg *stan.Msg) {
			msg.Ack()
			log.Println("algorithm sync heard ", msg.Data)
			createCmd := pb.CreateAlgorithmCommand{}
			err := json.Unmarshal(msg.Data, &createCmd)
			if err != nil {
				log.Print(err)
				return
			}

			if err = persistToQueryStore(&createCmd); err != nil {
				log.Println("failed to persist to query store", err)
			}
			if err := createAlgorithmCreatedEvent(&createCmd); err != nil {
				log.Println("failed to create algorithm created event", err)
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

func createAlgorithmCreatedEvent(createCmd *pb.CreateAlgorithmCommand) error {
	conn, err := grpc.Dial(grpcURI, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewEventStoreClient(conn)
	eventID, _ := uuid.NewV4()
	createdEvent := pb.AlgorithmCreatedEvent{
		Algorithm: createCmd.Algorithm,
		Id:        eventID.String(),
	}
	createdEventJSON, _ := json.Marshal(createdEvent)
	eid, _ := uuid.NewV4()
	event := &pb.Event{
		EventId:       eid.String(),
		EventType:     addedEvent,
		AggregateId:   createdEvent.Id,
		AggregateType: aggregate,
		EventData:     string(createdEventJSON),
		Channel:       addedEvent,
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

func persistToQueryStore(cmd *pb.CreateAlgorithmCommand) error {
	conn, err := grpc.Dial(queryStoreURI, grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "Unable to connect")
	}
	queryStoreClient := pb.NewAlgorithmQueryStoreClient(conn)
	added, err := queryStoreClient.CreateAlgorithm(context.Background(), cmd.Algorithm)
	if err != nil {
		return err
	}
	log.Println("persisted to query store: " + added.Id)
	return nil
}
