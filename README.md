# gRPC Client | Server | Gateway

This project aims to provide hands-on experience with basic gRPC, Golang concurrency, and telemetry. The primary focus is on understanding the fundamentals of these technologies, their integration, and their practical applications in modern software development.

- [gRPC Client | Server | Gateway](#grpc-client--server--gateway)
  - [Roadmap](#roadmap)
  - [Features](#features)
    - [gRPC](#grpc)
    - [Golang Concurrency](#golang-concurrency)
    - [Frontend application in `React.js`](#frontend-application-in-reactjs)
    - [OpenTelemetry with Traces, Logs, Metrics](#opentelemetry-with-traces-logs-metrics)
    - [Kubernetes](#kubernetes)
  - [gRPC Server \& Gateway](#grpc-server--gateway)
    - [Prerequisite](#prerequisite)
      - [Install gRPC gateway and related components](#install-grpc-gateway-and-related-components)
      - [Clone proto dependencies](#clone-proto-dependencies)
      - [Generate protobuf files (already generated)](#generate-protobuf-files-already-generated)
    - [Server](#server)
      - [Run Server](#run-server)
      - [Run Gateways for Server](#run-gateways-for-server)
  - [gRPC Client](#grpc-client)
    - [Testing via gateways or calling the server directly](#testing-via-gateways-or-calling-the-server-directly)
    - [Testing via frontend application](#testing-via-frontend-application)
      - [Prerequisite](#prerequisite-1)
      - [Run client-web application](#run-client-web-application)
  - [Telemetry](#telemetry)
  - [References](#references)

## Roadmap
A list of tasks or features that need to be completed, serving as a checklist or roadmap for the project.

-  Golang Concurrency
-  OpenTelemetry Client SDKs and Collector
-  Kubernetes Cluster
-  Kubernetes Operator

## Features 
### gRPC
Explains the implementation of both unary and streaming based [RPCs](https://book.systemsapproach.org/e2e/rpc.html).
  - `Unary RPC` - The client sends a single request to the server and gets a single response back,  similar to a traditional function call.
  - `Streaming RPC`: Allows for more complex interactions. There are three types:
    - `Server Streaming RPC`: The client sends a single request and receives a stream of responses.
    - `Client Streaming RPC`: The client sends a stream of requests and receives a single response.
    - `Bidirectional Streaming RPC`: Both client and server send a stream of messages to each other.

### Golang Concurrency
Highlights the use of goroutines, channels, and the select statement for concurrent programming.
  - `Goroutines`: Lightweight threads managed by the Go runtime, allowing for efficient concurrency.
  - `Channels`: Used for communication between goroutines, facilitating safe data exchange.
  - `Select Statement`: Enables waiting on multiple channel operations, helping in building concurrent and responsive applications.

### Frontend application in `React.js`
This frontend application was implemented using React.js to demonstrate how frontend applications can use gPRC based services.

  Please note gRPC-web currently supports 2 RPC modes.
  - Unary RPCs
  - Server-side Streaming RPCs (NOTE: Only when grpcwebtext mode is used.)

  `Client-side` and `Bi-directional` streaming is not currently supported and you can see in the generated `chat_pb.js` and `chat_grpc_web_pb.j` there's no impelemntation for it even we try to generate it using `protoc`.

### OpenTelemetry with Traces, Logs, Metrics
  OpenTelemetry is an observability framework, providing instrumentation to collect telemetry data (traces, logs, metrics).
  - `Traces`: Provide insights into the request paths and performance of the application by tracing the flow through different services.
  - `Logs`: Capture application events and errors for debugging and monitoring.
  - `Metrics`: Quantitative data about the system's performance and health (e.g., request count, latency).

### Kubernetes
Discusses the deployment and management of the application in a Kubernetes environment.
  - `Kubernetes`: An open-source container orchestration platform that automates the deployment, scaling, and management of containerized applications.
  - `Monitoring and Logging`: Integrating Kubernetes with Prometheus and OpenTelemetry to monitor and log the application's performance and health.
  - `Operator`: Details the implementation of a Kubernetes Operator to manage custom resources.


## gRPC Server & Gateway

### Prerequisite 

#### Install gRPC gateway and related components
```bash
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
```

#### Clone proto dependencies
```bash
git clone https://github.com/googleapis/googleapis.git
```

#### Generate protobuf files (already generated)
```bash
protoc -I. --go_out=plugins=grpc:chat --go-grpc_out=chat --grpc-gateway_out=chat --proto_path=./ --proto_path=./googleapis chat.proto
# or 
protoc -I. --go_out=chat --go-grpc_out=chat --grpc-gateway_out=chat --proto_path=./ --proto_path=./googleapis chat.proto
# or with error handling
protoc -I. --go_out=chat --go_opt=paths=source_relative --go-grpc_out=chat --go-grpc_opt=paths=source_relative --grpc-gateway_out=chat --grpc-gateway_opt=logtostderr=true,paths=source_relative --proto_path=./ --proto_path=./googleapis chat.proto
```

### Server

Please note if you are not using Golang workspaces you don't need `GOWORK=off`

#### Run Server
Run the gRPC server.
```bash
GOWORK=off go run server/*.go  
```

#### Run Gateways for Server
Below gateways were implemented to undestand and explain how gRPC acts under different protocols.
TODO: Need to add interceptors to explain how HTTP gateways decode serialised requests and responses.

```bash
GOWORK=off go run gateway/gateway-grpc-http1-1/*.go # Run the gRPC gateway that supports HTTP/1.1.
GOWORK=off go run gateway/gateway-http-http1-1/*.go # Run the HTTP gateway that supports HTTP/1.1.
GOWORK=off go run gateway/gateway-grpc-http2/*.go   # Run the gRPC gateway that supports HTTP/2.
GOWORK=off go run gateway/gateway-http-http2/*.go   # Run the HTTP gateway that supports HTTP/2.
``` 

## gRPC Client

### Testing via gateways or calling the server directly
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

#### Prerequisite 
1. Node.js (Tested in Node.js `v18.17.1`)
2. Revery Proxy: This is a reverse proxy that can front existing gRPC servers and expose their functionality using gRPC-Web protocol, allowing for the gRPC services to be consumed from browsers.
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

## Telemetry
Run `docker compose --env-file .env -f docker-compose.yaml up`

## References
1. gRPC Example https://tutorialedge.net/golang/go-grpc-beginners-tutorial/
2. OpenTelemetry Collectors https://github.com/open-telemetry/opentelemetry-demo/blob/main/docker-compose.yml


