package main

import (
	"context"
	"crypto/tls"
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

	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		http.Error(w, "Failed to connect to gRPC server", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := chat.NewCallServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	resp, err := client.SayHello(ctx, &reqBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error calling gRPC method: %v", err), http.StatusInternalServerError)
		return
	}

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

	srv := &http.Server{
		Addr:    ":8444",
		Handler: http.DefaultServeMux,
		TLSConfig: &tls.Config{
			NextProtos: []string{"h2", "http/1.1"},
		},
	}

	log.Println("Starting HTTP API Gateway on https://localhost:8444 with HTTP/2")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
