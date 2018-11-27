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
func (s store) GetAlgorithms(ctx context.Context, algos *pb.MultipleAlgorithms) (*pb.MultipleAlgorithms, error) {
	log.Print("query store: query algorithm request")

	for i, element := range algos.Algorithms {
		&pb.MultipleAlgorithms[i] {
			Name:       element.Name + " but better",
			Version:    element.Version,
			Id:         element.Id,
			Status:     "created",
			FileIDs:    nil,
			DatasetIDs: nil,
		}
	}

	return &pb.MultipleAlgorithms, nil
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
func (s store) AssociateFile(ctx context.Context, pair *pb.AlgorithmAndFile) (*pb.Algorithm, error) {
	log.Print("query store: association requested")

	db, err := sql.Open("postgres", connectionstring)
	defer db.Close()
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}

	sql := `UPDATE algorithm.algos SET fileIDs = array_append(fileIDs,` + pair.File.Id + `)
	WHERE algorithm.algos.id = ` + pair.Algorithm.Id
	log.Println("executing: ", sql)

	if resp, err := db.Exec(sql); err != nil {
		log.Fatal("Failed to add file ID to algorithm", err)
	} else {
		log.Println("Failed to add file ID to algorithm", resp)
	}

	return pair.Algorithm, nil
}
func (s store) CreateFile(ctx context.Context, file *pb.AlgorithmFile) (*pb.AlgorithmFile, error) {
	log.Print("query store: create algorithmFile request ", file)

	db, err := sql.Open("postgres", connectionstring)
	defer db.Close()
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}

	fileString := fmt.Sprintf("'%v', '%v', '%v', '%v'",
		file.Id,
		file.Content,
		file.Name,
		file.Filetype)
	log.Println("filestring: ", fileString)
	sql := "INSERT INTO algorithm.files VALUES (" + fileString + ")"
	log.Println("executing: ", sql)

	if resp, err := db.Exec(sql); err != nil {
		log.Fatal("Failed to persist algo to db", err)
	} else {
		log.Println("Persisted algorithm to db: ", resp)
	}

	return file, nil
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
		log.Println("Created algorithm table: ", resp)
	}

	if resp, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS algorithm.files
		(id UUID PRIMARY KEY,
			content STRING,
			name STRING,
			type STRING)`); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Created file table: ", resp)
	}
}
