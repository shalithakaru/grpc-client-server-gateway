# gRPC Client | Server | Gateway

This project is to learn gRPC from scratch. 

## Prerequisite 

```bash
# Below will install 
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

# Proto dependencies
git clone https://github.com/googleapis/googleapis.git

# Generating the protobuf files 
protoc -I. --go_out=plugins=grpc:chat --go-grpc_out=chat --grpc-gateway_out=chat --proto_path=./ --proto_path=./googleapis chat.proto
# or
protoc -I. --go_out=chat --go-grpc_out=chat --grpc-gateway_out=chat --proto_path=./ --proto_path=./googleapis chat.proto
# or for errors 
protoc -I. --go_out=chat --go_opt=paths=source_relative --go-grpc_out=chat --go-grpc_opt=paths=source_relative --grpc-gateway_out=chat --grpc-gateway_opt=logtostderr=true,paths=source_relative --proto_path=./ --proto_path=./googleapis chat.proto
```

## Commands

### Server
```bash
GOWORK=off go run server.go interceptor.go
```

### Gayeway HTTP/1.1
```bash
GOWORK=off go run gateway.go
``` 

## Test

### Testing via Gateway
```bash
curl -X POST -k http://localhost:8080/v1/sayhello -d '{"body": "Hello From HTTP/1.1!"}'
```

or 
```bash
GOWORK=off go run client.go interceptor.go --http
```

### Testing via Client
```bash
GOWORK=off go run client.go interceptor.go 
```

## History

```bash
protoc --go_out=plugins=grpc:chat --go-grpc_out=. chat.proto
```
