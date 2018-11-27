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
	port                 = ":50051"
	clusterID            = "test-cluster"
	eventStoreURI        = "eventstore:50051"
	queryStoreURI        = "dataset-querystore:50051"
	createDatasetChannel = "create-dataset"
	aggregate            = "Dataset"
	associateFileChannel = "associate-file"
)

type DataServer struct {
}

func (a *DataServer) CreateDataset(ctx context.Context, cmd *pb.CreateDatasetCommand) (*pb.Dataset, error) {
	log.Printf("Dataset Aggregate Received: ", cmd)
	var conn *grpc.ClientConn
	var err error
	for conn, err = grpc.Dial(eventStoreURI, grpc.WithInsecure()); err != nil; time.Sleep(time.Second * 5) {
		log.Printf(eventStoreURI + " is not available. Trying again")
		conn, err = grpc.Dial(eventStoreURI, grpc.WithInsecure())
	}

	dataID, _ := uuid.NewV4()
	cmd.Dataset.Id = dataID.String()
	eventID, _ := uuid.NewV4()
	cmdJSON, _ := json.Marshal(cmd)
	event := &pb.Event{
		EventId:       eventID.String(),
		EventType:     createDatasetChannel,
		AggregateType: aggregate,
		EventData:     string(cmdJSON),
		Channel:       createDatasetChannel,
	}
	client := pb.NewEventStoreClient(conn)
	log.Println("Dataset Aggregate: Sending Event")
	response, err := client.CreateEvent(ctx, event)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to add to event store")
	}
	if !response.IsSuccessful {
		return nil, errors.Wrap(err, "Failed to add to event store")
	}
	return cmd.Dataset, nil
}
func (a *DataServer) GetDataset(ctx context.Context, qry *pb.GetDatasetQuery) (*pb.Dataset, error) {
	log.Printf("Dataset Query Received: ", qry)
	var conn *grpc.ClientConn
	var err error
	for conn, err = grpc.Dial(queryStoreURI, grpc.WithInsecure()); err != nil; time.Sleep(time.Second * 5) {
		log.Printf(queryStoreURI + " is not available. Trying again")
		conn, err = grpc.Dial(queryStoreURI, grpc.WithInsecure())
	}
	queryStore := pb.NewDatasetQueryStoreClient(conn)
	response, err := queryStore.GetDataset(ctx, qry.Dataset)
	log.Print("received dataset: ", response)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get from dataset query store")
	}
	return response, nil
}
func (a *DataServer) AssociateFile(ctx context.Context, cmd *pb.AssociateFileCommand) (*pb.Dataset, error) {
	log.Printf("Dataset File Association Request Recieved: ", cmd)
	var conn *grpc.ClientConn
	var err error
	for conn, err = grpc.Dial(eventStoreURI, grpc.WithInsecure()); err != nil; time.Sleep(time.Second * 5) {
		log.Printf(eventStoreURI + " is not available. Trying again")
		conn, err = grpc.Dial(eventStoreURI, grpc.WithInsecure())
	}

	fileID, _ := uuid.NewV4()
	cmd.DatasetFile.Id = fileID.String()
	eventID, _ := uuid.NewV4()
	cmdJSON, _ := json.Marshal(cmd)
	event := &pb.Event{
		EventId:       eventID.String(),
		EventType:     associateFileChannel,
		AggregateType: aggregate,
		EventData:     string(cmdJSON),
		Channel:       associateFileChannel,
	}
	client := pb.NewEventStoreClient(conn)
	log.Println("Dataset Aggregate: Sending Event")
	response, err := client.CreateEvent(ctx, event)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to add to event store")
	}
	if !response.IsSuccessful {
		return nil, errors.Wrap(err, "Failed to add to event store")
	}
	return cmd.Dataset, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	log.Println("Dataset aggregate is running on:", port)
	pb.RegisterDatasetAggregateServer(s, &DataServer{})
	s.Serve(lis)
}
