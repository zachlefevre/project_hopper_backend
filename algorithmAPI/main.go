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
	createdEvent = "create-algorithm"
	aggregate    = "algorithm"
	grpcURI      = "algorithm-aggregate:50051"
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
	var algoIDBytes []byte
	_, err := r.Body.Read(algoIDBytes)
	if err != nil {
		log.Print(err)
		http.Error(w, "Failed to get algorithm", 500)
		return
	}
	algoID := string(algoIDBytes)

}
