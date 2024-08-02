# gRPC Client | Server | Gateway

This project is to learn gRPC from scratch. 

## Prerequisite 

### Install gRPC gateway and OpenAPI v2 plugins
```bash
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
```

### Clone proto dependencies
```bash
git clone https://github.com/googleapis/googleapis.git
```

### Generate protobuf files (already generated)
```bash
protoc -I. --go_out=plugins=grpc:chat --go-grpc_out=chat --grpc-gateway_out=chat --proto_path=./ --proto_path=./googleapis chat.proto
# or 
protoc -I. --go_out=chat --go-grpc_out=chat --grpc-gateway_out=chat --proto_path=./ --proto_path=./googleapis chat.proto
# or with error handling
protoc -I. --go_out=chat --go_opt=paths=source_relative --go-grpc_out=chat --go-grpc_opt=paths=source_relative --grpc-gateway_out=chat --grpc-gateway_opt=logtostderr=true,paths=source_relative --proto_path=./ --proto_path=./googleapis chat.proto
```

## Commands

### Server
Run the gRPC server.
```bash
GOWORK=off go run server/*.go  
```

### gRPC Gateway HTTP/1.1
Run the gRPC gateway that supports HTTP/1.1.
```bash
GOWORK=off go run gateway/gateway-grpc-http1-1/*.go
``` 

### HTTP Gateway HTTP/1.1
Run the HTTP gateway that supports HTTP/1.1.
```bash
GOWORK=off go run gateway/gateway-http-http1-1/*.go
``` 

### gRPC Gateway HTTP/2
Run the gRPC gateway that supports HTTP/2.
```bash
GOWORK=off go run gateway/gateway-grpc-http2/*.go
``` 

### HTTP Gateway HTTP/2
Run the HTTP gateway that supports HTTP/2.
```bash
GOWORK=off go run gateway/gateway-http-http2/*.go
``` 

## Test

### Testing via Gateway
Test the service using the gateway with a curl command.
```bash
curl -X POST -k http://localhost:8080/v1/sayhello -d '{"body": "Hello From HTTP/1.1!"}'
```

Use the client to test via gateways that support different protocols.
```bash
GOWORK=off go run client/*.go --http-gateway # call via HTTP gateway that supports HTTP/1.1
GOWORK=off go run client/*.go --grpc-gateway # call via gRPC gateway that supports HTTP/1.1
GOWORK=off go run client/*.go --http-gateway-http2 # call via HTTP gateway that supports HTTP/2
GOWORK=off go run client/*.go --grpc-gateway-http2 # call via gRPC gateway that supports HTTP/2
GOWORK=off go run client/*.go # call server directly 
```

## Other
This command explains how to use current directory.
```bash
protoc --go_out=plugins=grpc:chat --go-grpc_out=. chat.proto
```
