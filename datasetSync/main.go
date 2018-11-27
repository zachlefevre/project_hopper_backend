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
	baseClientID         = "dataset-repository"
	clusterID            = "test-cluster"
	aggregate            = "Dataset"
	natsURL              = "nats:4222"
	eventstoreURI        = "eventstore:50051"
	durableID            = "dataset-repository-durable"
	queryStoreURI        = "dataset-querystore:50051"
	addedEvent           = "dataset-added-to-query-store"
	fileAssociatedEvent  = "dataset-associated-with-file"
	createChannel        = "create-dataset"
	fileAssociateChannel = "associate-file"
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
			log.Println("dataset sync heard ", msg.Data)
			createCmd := pb.CreatedatasetCommand{}
			err := json.Unmarshal(msg.Data, &createCmd)
			if err != nil {
				log.Print(err)
				return
			}

			if err = persistDatasetToQueryStore(&createCmd); err != nil {
				log.Println("failed to persist dataset to query store", err)
			}
			if err := createDatasetCreatedEvent(&createCmd); err != nil {
				log.Println("failed to create dataset created event", err)
			}
		}, stan.DurableName(durableID),
			stan.MaxInflight(25),
			stan.SetManualAckMode(),
			stan.AckWait(aw),
		)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		sc.Subscribe(fileAssociateChannel, func(msg *stan.Msg) {
			msg.Ack()
			log.Println("file association sync heard ", msg.Data)
			associationCmd := pb.AssociateFileCommand{}
			err := json.Unmarshal(msg.Data, &associationCmd)
			if err != nil {
				log.Print(err)
				return
			}

			if err = associateFileInQueryStore(&associationCmd); err != nil {
				log.Println("failed to associate file query store", err)
			}
			if err = persistFileToQueryStore(&associationCmd); err != nil {
				log.Println("failed to persist file to query store", err)
			}
			if err := createFileCreatedEvent(&associationCmd); err != nil {
				log.Println("failed to create file associated event", err)
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

func createDatasetCreatedEvent(createCmd *pb.CreateDatasetCommand) error {
	conn, err := grpc.Dial(eventstoreURI, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewEventStoreClient(conn)

	createdEvent := pb.DatasetCreatedEvent{
		DatasetID: createCmd.Dataset.Id,
	}
	log.Println("Creating created event", addedEvent)
	createdEventJSON, _ := json.Marshal(createdEvent)
	eid, _ := uuid.NewV4()
	event := &pb.Event{
		EventId:       eid.String(),
		EventType:     addedEvent,
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

func persistDatasetToQueryStore(cmd *pb.CreateDatasetCommand) error {
	conn, err := grpc.Dial(queryStoreURI, grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "Unable to connect")
	}
	queryStoreClient := pb.NewDatasetQueryStoreClient(conn)
	added, err := queryStoreClient.CreateDataset(context.Background(), cmd.Dataset)
	if err != nil {
		return err
	}
	log.Println("persisted to query store: " + added.Id)
	return nil
}

//TODO
func createFileCreatedEvent(createCmd *pb.AssociateFileCommand) error {
	conn, err := grpc.Dial(eventstoreURI, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewEventStoreClient(conn)

	associatedEvent := pb.FileAssociatedWithDatasetEvent{
		DatasetID: createCmd.Dataset.Id,
		FileID:    createCmd.DatasetFile.Id,
	}
	log.Println("Creating created event", associatedEvent)
	associatedEventJSON, _ := json.Marshal(associatedEvent)
	eid, _ := uuid.NewV4()
	event := &pb.Event{
		EventId:       eid.String(),
		EventType:     fileAssociatedEvent,
		AggregateType: aggregate,
		EventData:     string(associatedEventJSON),
		Channel:       fileAssociatedEvent,
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

func associateFileInQueryStore(cmd *pb.AssociateFileCommand) error {
	conn, err := grpc.Dial(queryStoreURI, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		return errors.Wrap(err, "Unable to connect")
	}

	queryStoreClient := pb.NewDatasetQueryStoreClient(conn)
	created, err := queryStoreClient.CreateFile(context.Background(), cmd.DatasetFile)
	if err != nil {
		return err
	}
	log.Println("persisted to query store: " + created.Id)
	return nil
}

func persistFileToQueryStore(cmd *pb.AssociateFileCommand) error {
	conn, err := grpc.Dial(queryStoreURI, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		return errors.Wrap(err, "Unable to connect")
	}
	fileAndDataset := &pb.DatasetAndFile{
		File:    cmd.DatasetFile,
		Dataset: cmd.Dataset,
	}
	queryStoreClient := pb.NewDatasetQueryStoreClient(conn)
	updated, err := queryStoreClient.AssociateFile(context.Background(), fileAndDataset)
	if err != nil {
		return err
	}
	log.Println("persisted to query store: " + updated.Id)
	return nil
}
