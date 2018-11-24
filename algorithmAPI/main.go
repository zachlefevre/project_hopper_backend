package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/zachlefevre/project_hopper_backend/com"
)

const (
	aggregate = "algorithm"
	grpcURI   = "algorithm-aggregate:50051"
)

func main() {
	server := &http.Server{
		Addr:    ":3000",
		Handler: initRoutes(),
	}
	log.Println("Http Server Listening...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
func initRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/algorithms", createAlgorithm).Methods("POST")
	router.HandleFunc("/api/algorithms", getAlgorithm).Methods("GET")
	router.HandleFunc("/api/algorithms/file", addFile).Methods("POST")
	router.HandleFunc("/api", sig).Methods("GET")
	return router
}
func sig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<h1>Made by zachlefevre@gmail.com</h1>"))
}
func createAlgorithm(w http.ResponseWriter, r *http.Request) {
	var algo pb.Algorithm
	err := json.NewDecoder(r.Body).Decode(&algo)
	if err != nil {
		http.Error(w, "Invalid Algorithm", 500)
		return
	}

	cmdID, _ := uuid.NewV4()
	createCmd := pb.CreateAlgorithmCommand{
		Algorithm: &algo,
		CreatedOn: time.Now().Unix(),
		Id:        cmdID.String(),
	}

	resp, err := createAlgorithmRPC(&createCmd)
	if err != nil {
		log.Print(err)
		http.Error(w, "Failed to create algorithm", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	j, _ := json.Marshal(resp)
	w.Write(j)
}
func createAlgorithmRPC(cmd *pb.CreateAlgorithmCommand) (*pb.Algorithm, error) {
	conn, err := grpc.Dial(grpcURI, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewAlgorithmAggregateClient(conn)
	return client.CreateAlgorithm(context.Background(), cmd)
}
func getAlgorithm(w http.ResponseWriter, r *http.Request) {
	var algo pb.Algorithm
	err := json.NewDecoder(r.Body).Decode(&algo)
	if err != nil {
		http.Error(w, "Invalid algoritm", 500)
		return
	}

	queryID, _ := uuid.NewV4()
	getQuery := pb.GetAlgorithmQuery{
		Algorithm: &algo,
		CreatedOn: time.Now().Unix(),
		Id:        queryID.String(),
	}
	resp, err := getAlgorithmRPC(&getQuery)
	if err != nil {
		log.Print(err)
		http.Error(w, "Failed to get algorithm", 500)
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
func getAlgorithmRPC(query *pb.GetAlgorithmQuery) (*pb.Algorithm, error) {
	conn, err := grpc.Dial(grpcURI, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewAlgorithmAggregateClient(conn)
	return client.GetAlgorithm(context.Background(), query)
}
func createFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	algo := r.URL.Query().Get("version")
	if algo == "" {
		w.WriteHeader(http.StatusNotFound)
	}
	log.Printf(algo)
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Print(err)
		http.Error(w, "File upload failed", 500)
		return
	}
	defer file.Close()
	mimeType := header.Header.Get("Content-Type")
	switch mimeType {
	case "image/jpeg":
		log.Println(file)
	case "image/png":
		log.Println(file)
	default:
		log.Print(err)
		http.Error(w, "File type unsupported", 500)
	}
}

func addFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	algoName := r.URL.Query().Get("algorithm")
	version := r.URL.Query().Get("version")
	if algoName == "" {
		w.WriteHeader(http.StatusNotFound)
	}
	if version == "" {
		w.WriteHeader(http.StatusNotFound)
	}

	queryID, _ := uuid.NewV4()
	getQuery := pb.GetAlgorithmQuery{
		Algorithm: &pb.Algorithm{
			Name:    algoName,
			Version: version,
		},
		CreatedOn: time.Now().Unix(),
		Id:        queryID.String(),
	}
	algo, err := getAlgorithmRPC(&getQuery)
	if err != nil {
		log.Print(err)
		http.Error(w, "Failed to get algorithm", 500)
		return
	}

	var file pb.AlgorithmFile
	err = json.NewDecoder(r.Body).Decode(&file)
	if err != nil {
		http.Error(w, "Invalid file", 500)
		return
	}

	addFileCmd := pb.AssociateFileCommand{
		Algorithm:     algo,
		AlgorithmFile: &file,
	}
	resp, err := addFileRPC(&addFileCmd)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	j, _ := json.Marshal(resp)
	w.Write(j)
}
func addFileRPC(cmd *pb.AssociateFileCommand) (*pb.Algorithm, error) {
	conn, err := grpc.Dial(grpcURI, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewAlgorithmAggregateClient(conn)
	return client.AssociateFile(context.Background(), cmd)
}
