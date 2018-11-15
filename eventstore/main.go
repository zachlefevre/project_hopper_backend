package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/gofrs/uuid"

	"github.com/nats-io/go-nats-streaming"
	"google.golang.org/grpc"

	"github.com/zachlefevre/project_hopper_backend/natsutil"

	"github.com/zachlefevre/project_hopper_backend/com"
	"github.com/zachlefevre/project_hopper_backend/store"
)

const (
	port         = ":50051"
	clusterID    = "test-cluster"
	baseClientID = "eventstore"
	natsUrl      = "nats:4222"
)

type Nats struct {
	comp *natsutil.StreamingComponent
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	cid, _ := uuid.NewV4()
	clientId := baseClientID + cid.String()
	comp := natsutil.NewStreamingComponent(clientId)

	err = comp.ConnectToNatsStreamingService(clusterID, stan.NatsURL(natsUrl), stan.ConnectWait(100*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	log.Println("Event store listening on port:", port)
	pb.RegisterEventStoreServer(s, &Nats{
		comp: comp,
	})
	s.Serve(lis)
}

func publishEvent(comp *natsutil.StreamingComponent, event *pb.Event) {
	sc := comp.NATS()
	channel := event.Channel
	eventMsg := []byte(event.EventData)
	err := sc.Publish(channel, eventMsg)
	if err != nil {
		log.Fatal("failed to publish", err)
	}
	log.Println("Published message on channel: " + channel)
}

func (s *Nats) CreateEvent(ctx context.Context, evnt *pb.Event) (*pb.Response, error) {
	cmd_store := store.NewEventStore()
	err := cmd_store.PersistEvent(evnt)
	if err != nil {
		return nil, err
	}
	go publishEvent(s.comp, evnt)
	return &pb.Response{IsSuccessful: true}, nil
}

func (s *Nats) GetEvents(ctx context.Context, filter *pb.EventFilter) (*pb.EventResponse, error) {
	cmd_store := store.NewEventStore()
	resp := cmd_store.GetEvents()
	return &pb.EventResponse{
		Events: resp,
	}, nil
}
