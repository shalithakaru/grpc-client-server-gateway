package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"grpc-post-body-test/chat"

	"google.golang.org/grpc"
)

type ChatServiceServer struct {
	chat.UnimplementedChatServiceServer
}

type CallServiceServer struct {
	chat.UnimplementedCallServiceServer
}

// Unary RPC implementation
func (s *ChatServiceServer) UnaryChat(ctx context.Context, message *chat.Message) (*chat.Message, error) {
	response := &chat.Message{
		Body: "Hello From ChatService!",
	}
	return response, nil
}

// Client streaming RPC implementation
func (s *ChatServiceServer) ClientStreamChat(stream chat.ChatService_ClientStreamChatServer) error {
	var messages []string
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			// End of stream
			response := &chat.Message{
				Body: "Received messages: " + joinMessages(messages),
			}
			return stream.SendAndClose(response)
		}
		if err != nil {
			return err
		}
		messages = append(messages, message.Body)
	}
}

// Server streaming RPC implementation
func (s *ChatServiceServer) ServerStreamChat(message *chat.Message, stream chat.ChatService_ServerStreamChatServer) error {
	for i := 0; i < 5; i++ {
		response := &chat.Message{
			Body: "Message " + message.Body + " number " + fmt.Sprint(i+1),
		}
		if err := stream.Send(response); err != nil {
			return err
		}
	}
	return nil
}

// Bidirectional streaming RPC implementation
func (s *ChatServiceServer) BidirectionalStreamChat(stream chat.ChatService_BidirectionalStreamChatServer) error {
	// Goroutines are lightweight and scalable. They require minimal memory and startup time, allowing the server to handle many client streams concurrently without significant overhead.
	var wg sync.WaitGroup

	messageCh := make(chan *chat.Message)

	// Goroutine to receive messages from the client
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			msg, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					close(messageCh)
					return
				}
				log.Printf("Error receiving message: %v", err)
				return
			}
			messageCh <- msg
		}
	}()

	// Goroutine to send messages to the client
	wg.Add(1)
	go func() {
		defer wg.Done()
		for msg := range messageCh {
			log.Printf("Received message from client: %s", msg.Body)
			response := &chat.Message{Body: "Echo: " + msg.Body}
			if err := stream.Send(response); err != nil {
				log.Printf("Error sending message: %v", err)
				return
			}
		}
	}()

	// Wait for both goroutines to complete
	wg.Wait()
	return nil

	// TODO: Refactor. This was the original implementation wihout goroutines
	// for {
	// 	message, err := stream.Recv()
	// 	if err == io.EOF {
	// 		return nil
	// 	}
	// 	if err != nil {
	// 		return err
	// 	}

	// 	response := &chat.Message{
	// 		Body: "Received: " + message.Body,
	// 	}
	// 	if err := stream.Send(response); err != nil {
	// 		return err
	// 	}
	// }
}

func joinMessages(messages []string) string {
	result := ""
	for _, message := range messages {
		result += message + " "
	}
	return result
}

func (s *CallServiceServer) SayHello(ctx context.Context, message *chat.Message) (*chat.Message, error) {
	response := &chat.Message{
		Body: "Hello From ChatServiceTwo!",
	}
	return response, nil
}

func main() {
	fmt.Println("Go gRPC Beginners Tutorial!")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(serverInterceptor),
	)
	chat.RegisterChatServiceServer(grpcServer, &ChatServiceServer{})
	chat.RegisterCallServiceServer(grpcServer, &CallServiceServer{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
}
