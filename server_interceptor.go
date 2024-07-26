package main

import (
	"context"
	"encoding/json"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/proto"
)

func loggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}

	p, _ := peer.FromContext(ctx)

	log.Printf("Metadata: %v", md)
	log.Printf("Peer: %v", p)

	log.Printf("Request - Method:%s Peer:%v Metadata:%v", info.FullMethod, p, md)

	if msg, ok := req.(proto.Message); ok {
		rawRequest, err := proto.Marshal(msg)
		if err == nil {
			log.Printf("Serialized Request: %x", rawRequest)
		} else {
			log.Printf("Failed to serialize request: %v", err)
		}
	} else {
		log.Printf("Request is not a proto.Message: %T", req)
	}

	reqJSON, err := json.Marshal(req)
	if err == nil {
		log.Printf("Request JSON: %s", reqJSON)
	} else {
		log.Printf("Failed to marshal request to JSON: %v", err)
	}

	log.Printf("Raw Request: %v", req)

	h, err := handler(ctx, req)

	return h, err
}
