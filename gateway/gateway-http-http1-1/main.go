package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"

	"grpc-post-body-test/chat"
)

// handleSayHello handles the HTTP request and forwards it to the gRPC server
func handleSayHello(w http.ResponseWriter, r *http.Request) {
	// Read and parse the request body
	var reqBody chat.Message
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Set up a connection to the gRPC server
	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		http.Error(w, "Failed to connect to gRPC server", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := chat.NewChatServiceClient(conn)

	// Call the SayHello method
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	resp, err := client.SayHello(ctx, &reqBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error calling gRPC method: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert the gRPC response to JSON and write it to the HTTP response
	respJSON, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}

func main() {
	http.HandleFunc("/api/sayhello", handleSayHello)

	log.Println("Starting HTTP API Gateway on port 8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
