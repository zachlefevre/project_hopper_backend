package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/zachlefevre/project_hopper_backend/com"
	"google.golang.org/grpc"
)

const (
	created_event = "create-algorithm"
	aggregate     = "algorithm"
	grpcUri       = "eventstore:50051"
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
	router.HandleFunc("/api", version).Methods("GET")
	return router
}
func version(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<h1>V0</h1>"))
}
func createAlgorithm(w http.ResponseWriter, r *http.Request) {
	var createCmd pb.CreateAlgorithmCommand
	err := json.NewDecoder(r.Body).Decode(&createCmd)
	if err != nil {
		http.Error(w, "Invalid Algorithm", 500)
		return
	}
	cmdID, _ := uuid.NewV4()
	createCmd.Id = cmdID.String()
	algoID, _ := uuid.NewV4()
	createCmd.Algorithm.Id = algoID.String()
	createCmd.CreatedOn = time.Now().Unix()
	err = createAlgorithmRPC(createCmd)
	if err != nil {
		log.Print(err)
		http.Error(w, "Failed to create algorithm", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	j, _ := json.Marshal(createCmd)
	w.Write(j)
}
func createAlgorithmRPC(cmd pb.CreateAlgorithmCommand) error {
	conn, err := grpc.Dial(grpcUri, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to Connect: %v", err)
	}
	defer conn.Close()
	uid, err := uuid.NewV4()
	client := pb.NewEventStoreClient(conn)
	cmdJSON, _ := json.Marshal(cmd)
	event := &pb.Event{
		EventId:       uid.String(),
		EventType:     created_event,
		AggregateId:   cmd.Id,
		AggregateType: aggregate,
		EventData:     string(cmdJSON),
		Channel:       created_event,
	}
	resp, err := client.CreateEvent(context.Background(), event)
	if err != nil {
		return errors.Wrap(err, "Error from RPC server")
	}
	if resp.IsSuccessful {
		return nil
	} else {
		return errors.Wrap(err, "Error from RPC server")
	}
}
