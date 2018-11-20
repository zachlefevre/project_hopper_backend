package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/zachlefevre/project_hopper_backend/com"
)

const (
	port = ":50051"
)

type store struct {
}

func (s store) GetAlgorithm(ctx context.Context, algo *pb.Algorithm) (*pb.Algorithm, error) {
	log.Print("query store: query algorithm request")
	return &pb.Algorithm{
		Name:    algo.Name + "but better",
		Version: algo.Version,
		Status:  "created",
		Files:   nil,
	}, nil
}
func (s store) CreateAlgorithm(ctx context.Context, algo *pb.Algorithm) (*pb.Algorithm, error) {
	log.Print("query store: create algorithm request")
	return &pb.Algorithm{
		Name:    "tst",
		Version: "v0",
		Status:  "",
		Files:   nil,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	log.Println("Algorithm query store is running on:", port)
	pb.RegisterAlgorithmQueryStoreServer(s, &store{})
	s.Serve(lis)
}
