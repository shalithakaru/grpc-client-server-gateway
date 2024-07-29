package main

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"grpc-post-body-test/chat"
)

func main() {
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
