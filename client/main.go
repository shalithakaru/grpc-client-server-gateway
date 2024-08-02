package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"grpc-post-body-test/chat"
)

// These will call the server directly
func callUnary(client chat.ChatServiceClient) {
	message := &chat.Message{Body: "Hello From Client!"}
	response, err := client.UnaryChat(context.Background(), message)
	if err != nil {
		log.Fatalf("Error calling SayHello: %v", err)
	}
	log.Printf("Response from server: %s", response.Body)
}

func callClientStream(client chat.ChatServiceClient) {
	stream, err := client.ClientStreamChat(context.Background())
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}

	messages := []string{"First message", "Second message", "Third message"}

	for _, msg := range messages {
		err := stream.Send(&chat.Message{Body: msg})
		if err != nil {
			log.Fatalf("Error sending message: %v", err)
		}
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	log.Printf("Response from server: %s", response.Body)
}

func callServerStream(client chat.ChatServiceClient) {
	message := &chat.Message{Body: "Hello From Client!"}
	stream, err := client.ServerStreamChat(context.Background(), message)
	if err != nil {
		log.Fatalf("Error calling ServerStream: %v", err)
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error receiving response: %v", err)
		}
		log.Printf("Response from server: %s", response.Body)
	}
}

func callBidirectionalStream(client chat.ChatServiceClient) {
	stream, err := client.BidirectionalStreamChat(context.Background())
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}

	messages := []string{"First message", "Second message", "Third message"}

	waitc := make(chan struct{})

	// Send messages
	go func() {
		for _, msg := range messages {
			if err := stream.Send(&chat.Message{Body: msg}); err != nil {
				log.Fatalf("Error sending message: %v", err)
			}
			time.Sleep(1 * time.Second)
		}
		stream.CloseSend()
	}()

	// Receive messages
	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving response: %v", err)
			}
			log.Printf("Response from server: %s", response.Body)
		}
		close(waitc)
	}()

	<-waitc
}

func callServerDirectlyAllServices() {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure(), grpc.WithUnaryInterceptor(clientInterceptor))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	client := chat.NewChatServiceClient(conn)

	log.Println("Calling unary SayHello...")
	callUnary(client)

	log.Println("Calling client streaming ClientStream...")
	callClientStream(client)

	log.Println("Calling server streaming ServerStream...")
	callServerStream(client)

	log.Println("Calling bidirectional streaming BidirectionalStream...")
	callBidirectionalStream(client)
}

func callServerDirectly() {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure(), grpc.WithUnaryInterceptor(clientInterceptor))
	if err != nil {
		log.Fatalf("Did not connect: %s", err)
	}
	defer conn.Close()

	client := chat.NewChatServiceClient(conn)
	response, err := client.UnaryChat(context.Background(), &chat.Message{Body: "Hello From Client!"})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Printf("Response from server: %s", response.Body)
}

func callGRPCGatewayHTTP2() {
	message := chat.Message{
		Body: "Hello From Client via gRPC-Gateway HTTP/2!",
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Failed to marshal message to JSON: %v", err)
	}

	req, err := http.NewRequestWithContext(context.Background(), "POST", "http://localhost:8443/v1/sayhello", bytes.NewBuffer(messageJSON))
	if err != nil {
		log.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Received non-OK HTTP status: %s", resp.Status)
	}

	var responseMessage chat.Message
	if err := json.NewDecoder(resp.Body).Decode(&responseMessage); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}

	fmt.Printf("Response from server: %s\n", responseMessage.Body)
}

func callGRPCGateway() {
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

// HTTP Client implementation for new API Gateway
func callHTTPGateway() {
	// Define the message to send
	message := chat.Message{
		Body: "Hello From Client via New API Gateway!",
	}

	// Convert the message to JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Failed to marshal message to JSON: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequestWithContext(context.Background(), "POST", "http://localhost:8081/api/sayhello", bytes.NewBuffer(messageJSON))
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

func callHTTPGatewayHTTP2() {
	message := chat.Message{
		Body: "Hello From Client via HTTP API Gateway HTTP/2!",
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Failed to marshal message to JSON: %v", err)
	}

	req, err := http.NewRequestWithContext(context.Background(), "POST", "http://localhost:8444/api/sayhello", bytes.NewBuffer(messageJSON))
	if err != nil {
		log.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Received non-OK HTTP status: %s", resp.Status)
	}

	var responseMessage chat.Message
	if err := json.NewDecoder(resp.Body).Decode(&responseMessage); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}

	fmt.Printf("Response from server: %s\n", responseMessage.Body)
}

func main() {
	// Use a flag to determine which client to use
	useHTTPGateway := flag.Bool("http-gateway", false, "Use HTTP client")
	useGRPCGatway := flag.Bool("grpc-gateway", false, "Use HTTP client for new API Gateway")
	useHTTPGatewayHTTP2 := flag.Bool("http-gateway-http2", false, "Use HTTP client for HTTP API Gateway HTTP/2")
	useGRPCGatewayHTTP2 := flag.Bool("grpc-gateway-http2", false, "Use HTTP client for gRPC-Gateway HTTP/2")

	flag.Parse()

	if *useHTTPGateway {
		callHTTPGateway()
	} else if *useGRPCGatway {
		callGRPCGateway()
	} else if *useHTTPGatewayHTTP2 {
		callHTTPGatewayHTTP2()
	} else if *useGRPCGatewayHTTP2 {
		callGRPCGatewayHTTP2()
	} else {
		// callServerDirectly()
		callServerDirectlyAllServices()
	}
}
