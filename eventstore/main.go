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
	baseClientId = "eventstore"
	natsUrl      = "nats:4222"
)

type Server struct {
	*natsutil.StreamingComponent
}

func main() {
	cid, _ := uuid.NewV4()
	clientId := baseClientId + cid.String()
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	comp := natsutil.NewStreamingComponent(clientId)

	err = comp.ConnectToNatsStreamingService(clusterID, stan.NatsURL(natsUrl), stan.ConnectWait(100*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	log.Println("Event store listening on port:", port)
	pb.RegisterEventStoreServer(s, &Server{
		StreamingComponent: comp,
	})
	s.Serve(lis)
}

func (s *Server) CreateEvent(ctx context.Context, evnt *pb.Event) (*pb.Response, error) {
	cmd_store := store.NewEventStore()
	err := cmd_store.PersistEvent(evnt)
	if err != nil {
		return nil, err
	}
	go publishEvent(s.StreamingComponent, evnt)
	return &pb.Response{IsSuccessful: true}, nil
}

func (s *Server) GetEvents(ctx context.Context, filter *pb.EventFilter) (*pb.EventResponse, error) {
	cmd_store := store.NewEventStore()
	resp := cmd_store.GetEvents()
	return &pb.EventResponse{
		Events: resp,
	}, nil
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
