package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"time"

	"github.com/pkg/errors"

	"github.com/gofrs/uuid"

	"google.golang.org/grpc"

	"github.com/zachlefevre/project_hopper_backend/com"
)

const (
	port                   = ":50051"
	clusterID              = "test-cluster"
	eventStoreURI          = "eventstore:50051"
	queryStoreURI          = "querystore:50051"
	createAlgorithmChannel = "create-algorithm"
	aggregate              = "Algorithm"
)

type AlgoServer struct {
}

func (a *AlgoServer) CreateAlgorithm(ctx context.Context, cmd *pb.CreateAlgorithmCommand) (*pb.Algorithm, error) {
	log.Printf("Algorithm Aggregate Received: ", cmd)
	var conn *grpc.ClientConn
	var err error
	for conn, err = grpc.Dial(eventStoreURI, grpc.WithInsecure()); err != nil; time.Sleep(time.Second * 5) {
		log.Printf(eventStoreURI + " is not available. Trying again")
		conn, err = grpc.Dial(eventStoreURI, grpc.WithInsecure())
	}

	for _, f := range cmd.Algorithm.Files {
		fileID, _ := uuid.NewV4()
		f.Id = fileID.String()
	}

	algoID, _ := uuid.NewV4()
	cmd.Algorithm.Id = algoID.String()
	eventID, _ := uuid.NewV4()
	cmdJSON, _ := json.Marshal(cmd)
	event := &pb.Event{
		EventId:       eventID.String(),
		EventType:     createAlgorithmChannel,
		AggregateType: aggregate,
		EventData:     string(cmdJSON),
		Channel:       createAlgorithmChannel,
	}
	log.Println("generating new client")
	client := pb.NewEventStoreClient(conn)
	log.Println("sending event")
	response, err := client.CreateEvent(ctx, event)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to add to event store")
	}
	if !response.IsSuccessful {
		return nil, errors.Wrap(err, "Failed to add to event store")
	}
	return cmd.Algorithm, nil
}
func (a *AlgoServer) GetAlgorithm(ctx context.Context, qry *pb.GetAlgorithmQuery) (*pb.Algorithm, error) {
	log.Printf("Algorithm Query Received: ", qry)
	var conn *grpc.ClientConn
	var err error
	for conn, err = grpc.Dial(queryStoreURI, grpc.WithInsecure()); err != nil; time.Sleep(time.Second * 5) {
		log.Printf(queryStoreURI + " is not available. Trying again")
		conn, err = grpc.Dial(queryStoreURI, grpc.WithInsecure())
	}
	queryStore := pb.NewAlgorithmQueryStoreClient(conn)
	response, err := queryStore.GetAlgorithm(ctx, qry.Algorithm)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get from algorithm query store")
	}
	return response, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	log.Println("Algorithm aggregate is running on:", port)
	pb.RegisterAlgorithmAggregateServer(s, &AlgoServer{})
	s.Serve(lis)
}
