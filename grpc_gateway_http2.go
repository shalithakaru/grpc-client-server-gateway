package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"grpc-post-body-test/chat"
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := chat.RegisterChatServiceTwoHandlerFromEndpoint(ctx, mux, "localhost:9000", opts)
	if err != nil {
		return err
	}

	srv := &http.Server{
		Addr:    ":8443",
		Handler: mux,
		TLSConfig: &tls.Config{
			NextProtos: []string{"h2", "http/1.1"},
		},
	}

	log.Println("Serving gRPC-Gateway on https://localhost:8443 with HTTP/2")
	return srv.ListenAndServe()
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
