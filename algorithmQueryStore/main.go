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

	db, err := sql.Open("postgres", connectionstring)
	defer db.Close()
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	var sqlString string
	if algo.Id != "" {
		sqlString = "SELECT * FROM algorithm.algos WHERE ID = " + "'" + algo.Id + "'"
	} else {
		sqlString = "SELECT * FROM algorithm.algos WHERE NAME = " + "'" + algo.Name + "'" + "AND Version = " + "'" + algo.Version + "'"
	}

	log.Println("executing: ", sqlString)

	resp, err := db.Query(sqlString)
	if err != nil {
		log.Fatal("Failed to get algorithm to db", err)
	} else {
		log.Println("Queried algorithm from db: ", resp)
	}

	defer resp.Close()
	ds := pb.Algorithm{}
	fileIDsB := []uint8{}
	datasetIDsB := []uint8{}
	for resp.Next() {
		err := resp.Scan(&ds.Id, &ds.Name, &ds.Version, &ds.Status, &fileIDsB, &datasetIDsB)
		for _, fid := range fileIDsB {
			ds.FileIDs = append(ds.FileIDs, (string(fid)))
		}

		log.Println("fileIDs", ds.FileIDs)
		if err != nil {
			log.Fatal(err)
		}

		for _, dsid := range datasetIDsB {
			ds.DatasetIDs = append(ds.DatasetIDs, (string(dsid)))
		}

		log.Println("datasetIDs", ds.DatasetIDs)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("returning " + ds.Name)
		return &ds, nil
	}
	err = resp.Err()
	if err != nil {
		log.Fatal(err)
	}
	return &pb.Algorithm{}, nil
}
func (s store) GetAlgorithms(ctx context.Context, algos *pb.Algorithm) (*pb.MultipleAlgorithms, error) {
	log.Print("query store: query algorithm request")

	return nil, nil
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
