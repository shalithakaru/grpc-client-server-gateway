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

### Testing via gateways or server directly
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
### Testing via frontend application
Please note gRPC-web currently supports 2 RPC modes 
- Unary RPCs (example)
- Server-side Streaming RPCs (example) (NOTE: Only when grpcwebtext mode is used.)

Client-side and Bi-directional streaming is not currently supported and you can see chat_pb.js and chat_grpc_web_pb.js doesn't have the impelemntation for it even we try to generate it.

#### Prerequisite 
1. Node.js
2. Revery Proxy
This is a reverse proxy that can front existing gRPC servers and expose their functionality using gRPC-Web protocol, allowing for the gRPC services to be consumed from browsers.
More information here https://github.com/improbable-eng/grpc-web/blob/master/go/grpcwebproxy/README.md 
```bash
git clone https://github.com/improbable-eng/grpc-web.git
go install ./go/grpcwebproxy
# Run the proxy. Please note we are running the gRPC server on port 9000
grpcwebproxy --backend_addr=localhost:9000 --run_tls_server=false --allow_all_origins
```
3. Tools to generate protobuf (This will be installed globally)
```
npm install -g grpc-tools
npm install -g protoc-gen-grpc-web
```
4. Generate protobuf for frontend
```bash
protoc -I=. chat.proto \
  --js_out=import_style=commonjs,binary:./client-web/src/grpc \
  --grpc-web_out=import_style=commonjs,mode=grpcwebtext:./client-web/src/grpc --proto_path=./googleapis

# please note there's an small issue with the above genearete as it refers 
# var google_api_annotations_pb in both files and we need to comment it for now.
# need to find how to properly generate ./google/api/annotations_pb.js later
```

#### Run client-web application
```bash
cd client-web
npm start
```

## Other
This command explains how to use current directory.
```bash
protoc --go_out=plugins=grpc:chat --go-grpc_out=. chat.proto
```
