package main

import (
	"context"
	"encoding/json"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

// clientLoggingInterceptor is a gRPC client interceptor for logging requests and responses
func clientLoggingInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	// Log the method and the request
	log.Printf("Method: %s", method)

	// Serialize the request to log raw bytes
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

	// Convert the request to JSON for readability
	reqJSON, err := json.Marshal(req)
	if err == nil {
		log.Printf("Request JSON: %s", reqJSON)
	} else {
		log.Printf("Failed to marshal request to JSON: %v", err)
	}

	// Invoke the RPC call
	err = invoker(ctx, method, req, reply, cc, opts...)

	// Log the response
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

	// Convert the response to JSON for readability
	respJSON, err := json.Marshal(reply)
	if err == nil {
		log.Printf("Response JSON: %s", respJSON)
	} else {
		log.Printf("Failed to marshal response to JSON: %v", err)
	}

	return err
}