package main

import (
	"context"
	"encoding/json"
	"log"
	"runtime"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/zachlefevre/project_hopper_backend/com"

	"github.com/nats-io/go-nats-streaming"
	"github.com/zachlefevre/project_hopper_backend/natsutil"
	"google.golang.org/grpc"
)

const (
	client_id         = "algorithm-repository"
	cluster_id        = "test-cluster"
	subscribe_channel = "create-algorithm"
	grpcUri           = "localhost:50051"
	event             = "algorithm-added-to-registry"
	aggregate         = "algorithms"
	durableID         = "algorithm-repository-durable"
)

func main() {
	comp := natsutil.NewStreamingComponent(client_id)

	err := comp.ConnectToNatsStreamingService(cluster_id, stan.NatsURL(stan.DefaultNatsURL))
	if err != nil {
		log.Fatal(err)
	}

	sc := comp.NATS()

	aw, _ := time.ParseDuration("60s")
	sc.Subscribe(subscribe_channel, func(msg *stan.Msg) {
		msg.Ack()
		createCmd := pb.CreateAlgorithmCommand{}
		err := json.Unmarshal(msg.Data, &createCmd)
		if err != nil {
			log.Print(err)
			return
		}
		createdCmd := pb.AlgorithmCreatedCommand{
			Algorithm: createCmd.Algorithm,
		}
		log.Println("Algorithm has been added to the algorithm repository", createCmd.Algorithm.Name)
		if err := createAlgorithmCreatedEvent(createdCmd); err != nil {

		}
	}, stan.DurableName(durableID),
		stan.MaxInflight(25),
		stan.SetManualAckMode(),
		stan.AckWait(aw),
	)
	runtime.Goexit()
}

func createAlgorithmCreatedEvent(cmd pb.AlgorithmCreatedCommand) error {
	conn, err := grpc.Dial(grpcUri, grpc.WithInsecure())
	if err != nil {
		log.Fatal("Unable to connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewEventStoreClient(conn)
	cmdJson, _ := json.Marshal(cmd)
	uuid, _ := uuid.NewV4()
	event := &pb.Event{
		EventId:       uuid.String(),
		EventType:     event,
		AggregateId:   cmd.Id,
		AggregateType: aggregate,
		EventData:     string(cmdJson),
		Channel:       event,
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
