package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"grpc-post-body-test/chat"
)

func callServerDirectly() {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure(), grpc.WithUnaryInterceptor(clientInterceptor))
	if err != nil {
		log.Fatalf("Did not connect: %s", err)
	}
	defer conn.Close()

	client := chat.NewChatServiceClient(conn)
	response, err := client.SayHello(context.Background(), &chat.Message{Body: "Hello From Client!"})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Printf("Response from server: %s", response.Body)
}

func callGateway() {
	// HTTP Client implementation
	// Define the message to send
	message := chat.Message{
		Body: "Hello From Client via HTTP!",
	}

	// Convert the message to JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Failed to marshal message to JSON: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequestWithContext(context.Background(), "POST", "http://localhost:8080/v1/sayhello", bytes.NewBuffer(messageJSON))
	if err != nil {
		log.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Set the appropriate headers
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client with a timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Received non-OK HTTP status: %s", resp.Status)
	}

	// Decode the response
	var responseMessage chat.Message
	if err := json.NewDecoder(resp.Body).Decode(&responseMessage); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}

	// Print the response
	fmt.Printf("Response from server: %s\n", responseMessage.Body)
}

func main() {
	// Use a flag to determine which client to use
	useHTTP := flag.Bool("http", false, "Use HTTP client")
	flag.Parse()

	if *useHTTP {
		callGateway()
	} else {
		callServerDirectly()
	}
}
