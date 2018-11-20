package main

import (
	"context"
	"log"

	"github.com/pkg/errors"
	"github.com/zachlefevre/project_hopper_backend/com"
	"google.golang.org/grpc"
)

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
