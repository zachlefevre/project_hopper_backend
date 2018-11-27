package main

import (
	"bytes"
	"fmt"
	"strings"

	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	pb "github.com/zachlefevre/project_hopper_backend/com_ds"
)

const (
	aggregate = "dataset"
	grpcURI   = "dataset-aggregate:50051"
)

func main() {
	server := &http.Server{
		Addr:    ":80",
		Handler: initRoutes(),
	}
	log.Println("Http Server Listening...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
func initRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/datasets", createDataset).Methods("POST")
	router.HandleFunc("/api/datasets", getDataset).Methods("GET")
	router.HandleFunc("/api/datasets/file", addFile).Methods("POST")
	router.HandleFunc("/api", sig).Methods("GET")
	return router
}
func sig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<h1>Made by zachlefevre@gmail.com</h1>"))
}
func createDataset(w http.ResponseWriter, r *http.Request) {
	var data pb.Dataset
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid Dataset", 500)
		return
	}
	if data.Name == "" {
		http.Error(w, "Invalid Dataset. Provide name", 500)
		return
	}
	if data.Version == "" {
		http.Error(w, "Invalid Dataset. Provide version", 500)
		return
	}

	cmdID, _ := uuid.NewV4()
	createCmd := pb.CreateDatasetCommand{
		Dataset:   &data,
		CreatedOn: time.Now().Unix(),
		Id:        cmdID.String(),
	}

	resp, err := createDatasetRPC(&createCmd)
	if err != nil {
		log.Print(err)
		http.Error(w, "Failed to create dataset", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	j, _ := json.Marshal(resp)
	w.Write(j)
}
func createDatasetRPC(cmd *pb.CreateDatasetCommand) (*pb.Dataset, error) {
	conn, err := grpc.Dial(grpcURI, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewDatasetAggregateClient(conn)
	return client.CreateDataset(context.Background(), cmd)
}
func getDataset(w http.ResponseWriter, r *http.Request) {
	var data pb.Dataset
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid dataset", 500)
		return
	}

	queryID, _ := uuid.NewV4()
	getQuery := pb.GetDatasetQuery{
		Dataset:   &data,
		CreatedOn: time.Now().Unix(),
		Id:        queryID.String(),
	}
	resp, err := getDatasetRPC(&getQuery)
	if err != nil {
		log.Print(err)
		http.Error(w, "Failed to get dataset", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if resp == nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusFound)
	}
	j, _ := json.Marshal(resp)
	w.Write(j)
}
func getDatasetRPC(query *pb.GetDatasetQuery) (*pb.Dataset, error) {
	conn, err := grpc.Dial(grpcURI, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewDatasetAggregateClient(conn)
	return client.GetDataset(context.Background(), query)
}

func addFile(w http.ResponseWriter, r *http.Request) {
	dataName := r.URL.Query().Get("dataset")
	version := r.URL.Query().Get("version")
	if dataName == "" {
		w.WriteHeader(http.StatusNotFound)
	}
	if version == "" {
		w.WriteHeader(http.StatusNotFound)
	}

	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fileContent := make([]byte, 100)
	_, err = file.Read(fileContent)
	if err != nil {
		log.Print(err)
		http.Error(w, "Failed to read file", 500)
		return
	}
	fileContent = bytes.Trim(fileContent, "\x00")
	fileStr := string(fileContent)

	cd := strings.Split(handler.Header.Get("Content-Disposition"), "; ")
	nameSection := cd[2]
	log.Println(cd, nameSection)
	name := strings.Split(nameSection, "=")[1]
	name = name[1 : len(name)-1]
	log.Println(name)
	queryID, _ := uuid.NewV4()
	getQuery := pb.GetDatasetQuery{
		Dataset: &pb.Dataset{
			Name:    dataName,
			Version: version,
		},
		CreatedOn: time.Now().Unix(),
		Id:        queryID.String(),
	}
	data, err := getDatasetRPC(&getQuery)
	if err != nil {
		log.Print(err)
		http.Error(w, "Failed to get dataset", 500)
		return
	}
	log.Println(handler.Header)
	dataFile := pb.DatasetFile{
		Content:  fileStr,
		Name:     name,
		Filetype: handler.Header.Get("Content-Type"),
	}

	addFileCmd := pb.AssociateFileCommand{
		Dataset:     data,
		DatasetFile: &dataFile,
	}
	resp, err := addFileRPC(&addFileCmd)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	j, _ := json.Marshal(resp)
	w.Write(j)
}
func addFileRPC(cmd *pb.AssociateFileCommand) (*pb.Dataset, error) {
	conn, err := grpc.Dial(grpcURI, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewDatasetAggregateClient(conn)
	return client.AssociateFile(context.Background(), cmd)
}
