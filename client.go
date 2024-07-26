package main

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"grpc-post-body-test/chat"
)

func main() {

	var conn *grpc.ClientConn
	conn, err := grpc.NewClient(":9000", grpc.WithInsecure(), grpc.WithUnaryInterceptor(clientLoggingInterceptor))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := chat.NewChatServiceClient(conn)

	response, err := c.SayHello(context.Background(), &chat.Message{Body: "Hello From Client!"})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Printf("Response from server: %s", response.Body)

}
