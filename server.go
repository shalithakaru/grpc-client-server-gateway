package main

import (
	"context"
	"fmt"
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

func (s *ChatServiceServer) SayHello(ctx context.Context, message *chat.Message) (*chat.Message, error) {
	response := &chat.Message{
		Body: "Hello From ChatService!",
	}
	return response, nil
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
