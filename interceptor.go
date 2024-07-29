package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/proto"
)

// loggingInterceptor is a gRPC server interceptor for logging requests and responses.
func serverInterceptor(
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

// clientLoggingInterceptor is a gRPC client interceptor for logging requests and responses.
func clientInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	start := time.Now()

	log.Printf("Method: %s", method)

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

	err = invoker(ctx, method, req, reply, cc, opts...)
	log.Printf("Method: %s Duration: %s Error: %v", method, time.Since(start), err)

	if msg, ok := reply.(proto.Message); ok {
		rawResponse, err := proto.Marshal(msg)
		if err == nil {
			log.Printf("Serialized Response: %x", rawResponse)
		} else {
			log.Printf("Failed to serialize response: %v", err)
		}
	} else {
		log.Printf("Response is not a proto.Message: %T", reply)
	}

	respJSON, err := json.Marshal(reply)
	if err == nil {
		log.Printf("Response JSON: %s", respJSON)
	} else {
		log.Printf("Failed to marshal response to JSON: %v", err)
	}

	return err
}
