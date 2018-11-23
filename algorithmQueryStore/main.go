package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"strings"

	"google.golang.org/grpc"

	_ "github.com/lib/pq"
	"github.com/zachlefevre/project_hopper_backend/com"
)

const (
	port             = ":50051"
	connectionstring = "postgresql://root@cockroachdb-public:26257/?sslmode=disable"
)

type store struct {
}

func (s store) GetAlgorithm(ctx context.Context, algo *pb.Algorithm) (*pb.Algorithm, error) {
	log.Print("query store: query algorithm request")
	return &pb.Algorithm{
		Name:       algo.Name + " but better",
		Version:    algo.Version,
		Id:         algo.Id,
		Status:     "created",
		FileIDs:    nil,
		DatasetIDs: nil,
	}, nil
}
func (s store) CreateAlgorithm(ctx context.Context, algo *pb.Algorithm) (*pb.Algorithm, error) {
	log.Print("query store: create algorithm request")

	db, err := sql.Open("postgres", connectionstring)
	defer db.Close()
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}

	var fileIDs []string
	var datasetIDs []string
	for _, id := range algo.FileIDs {
		fileIDs = append(fileIDs, `'`+id+`'`)
	}
	for _, id := range algo.DatasetIDs {
		datasetIDs = append(datasetIDs, `'`+id+`'`)
	}
	algoString := fmt.Sprintf("'%v', '%v', '%v', 'Created', ARRAY[%v], ARRAY[%v]",
		algo.Id,
		algo.Name,
		algo.Version,
		strings.Join(fileIDs, ","),
		strings.Join(datasetIDs, ","))
	sql := "INSERT INTO algorithm.algos VALUES (" + algoString + ")"
	log.Println("executing: ", sql)

	if resp, err := db.Exec(
		sql); err != nil {
		log.Fatal("Failed to persist algo to db", err)
	} else {
		log.Println("Persisted algorithm to db: ", resp)
	}

	return algo, nil
}

func main() {
	initDB()
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	log.Println("Algorithm query store is running on:", port)
	pb.RegisterAlgorithmQueryStoreServer(s, &store{})
	s.Serve(lis)
}

func initDB() {
	db, err := sql.Open("postgres", connectionstring)
	defer db.Close()
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}

	if resp, err := db.Exec(
		"CREATE DATABASE IF NOT EXISTS algorithm"); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Created Database: ", resp)
	}

	if resp, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS algorithm.algos
		(id UUID PRIMARY KEY,
			name STRING,
			version STRING,
			status STRING,
			fileIDs STRING[],
			datasetIDs STRING[])`); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Created table: ", resp)
	}
}
