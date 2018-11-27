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

func (s store) GetDataset(ctx context.Context, data *pb.Dataset) (*pb.Dataset, error) {
	log.Print("query store: query dataset request")
	return &pb.Dataset{
		Name:       data.Name + " but better",
		Version:    data.Version,
		Id:         data.Id,
		Status:     "created",
		FileIDs:    nil,
		DatasetIDs: nil,
	}, nil
}
func (s store) GetDatasets(ctx context.Context, data *pb.MultipleDatasets) (*pb.MultipleDatasets, error) {
	log.Print("query store: query algorithm request")

	for i, element := range data.Datasets {
		&pb.MultipleDatasets[i] {
			Name:       element.Name + " but better",
			Version:    element.Version,
			Id:         element.Id,
			Status:     "created",
			FileIDs:    nil,
			DatasetIDs: nil,
		}
	}

	return &pb.MultipleDatasets, nil
}
func (s store) CreateDataset(ctx context.Context, data *pb.Dataset) (*pb.Dataset, error) {
	log.Print("query store: create Dataset request")

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
	dataString := fmt.Sprintf("'%v', '%v', '%v', 'Created', ARRAY[%v], ARRAY[%v]",
		data.Id,
		data.Name,
		data.Version,
		strings.Join(fileIDs, ","),
		strings.Join(datasetIDs, ","))
	sql := "INSERT INTO dataset.datas VALUES (" + dataString + ")"
	log.Println("executing: ", sql)

	if resp, err := db.Exec(
		sql); err != nil {
		log.Fatal("Failed to persist data to db", err)
	} else {
		log.Println("Persisted dataset to db: ", resp)
	}

	return data, nil
}
func (s store) AssociateFile(ctx context.Context, pair *pb.DatasetAndFile) (*pb.Dataset, error) {
	log.Print("query store: association requested")

	db, err := sql.Open("postgres", connectionstring)
	defer db.Close()
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}

	sql := `UPDATE dataset.datas SET fileIDs = array_append(fileIDs,` + pair.File.Id + `)
	WHERE dataset.datas.id = ` + pair.Dataset.Id
	log.Println("executing: ", sql)

	if resp, err := db.Exec(sql); err != nil {
		log.Fatal("Failed to add file ID to dataset", err)
	} else {
		log.Println("Failed to add file ID to dataset", resp)
	}

	return pair.Dataset, nil
}
func (s store) CreateFile(ctx context.Context, file *pb.DatasetFile) (*pb.DatasetFile, error) {
	log.Print("query store: create datsetFile request ", file)

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
	sql := "INSERT INTO dataset.files VALUES (" + fileString + ")"
	log.Println("executing: ", sql)

	if resp, err := db.Exec(sql); err != nil {
		log.Fatal("Failed to persist data to db", err)
	} else {
		log.Println("Persisted dataset to db: ", resp)
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
	log.Println("Dataset query store is running on:", port)
	pb.RegisterDatasetQueryStoreServer(s, &store{})
	s.Serve(lis)
}

func initDB() {
	db, err := sql.Open("postgres", connectionstring)
	defer db.Close()
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}

	if resp, err := db.Exec(
		"CREATE DATABASE IF NOT EXISTS dataset"); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Created Database: ", resp)
	}

	if resp, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS dataset.datas
		(id UUID PRIMARY KEY,
			name STRING,
			version STRING,
			status STRING,
			fileIDs STRING[],
			datasetIDs STRING[])`); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Created dataset table: ", resp)
	}

	if resp, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS dataset.files
		(id UUID PRIMARY KEY,
			content STRING,
			name STRING,
			type STRING)`); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Created file table: ", resp)
	}
}
