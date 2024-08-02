package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"grpc-post-body-test/chat"

	"google.golang.org/grpc"
)

type ChatServiceServer struct {
	chat.UnimplementedChatServiceServer
}

type ChatServiceTwoServer struct {
	chat.UnimplementedChatServiceTwoServer
}

// Unary RPC implementation
func (s *ChatServiceServer) SayHello(ctx context.Context, message *chat.Message) (*chat.Message, error) {
	response := &chat.Message{
		Body: "Hello From ChatService!",
	}
	return response, nil
}

// Client streaming RPC implementation
func (s *ChatServiceServer) ClientStream(stream chat.ChatService_ClientStreamServer) error {
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
func (s *ChatServiceServer) ServerStream(message *chat.Message, stream chat.ChatService_ServerStreamServer) error {
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
func (s *ChatServiceServer) BidirectionalStream(stream chat.ChatService_BidirectionalStreamServer) error {
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		response := &chat.Message{
			Body: "Received: " + message.Body,
		}
		if err := stream.Send(response); err != nil {
			return err
		}
	}
}

func joinMessages(messages []string) string {
	result := ""
	for _, message := range messages {
		result += message + " "
	}
	return result
}

func (s *ChatServiceTwoServer) SayHello(ctx context.Context, message *chat.Message) (*chat.Message, error) {
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
	chat.RegisterChatServiceTwoServer(grpcServer, &ChatServiceTwoServer{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
}
