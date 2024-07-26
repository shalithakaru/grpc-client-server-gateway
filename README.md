Server
GOWORK=off go run server.go interceptor.go

Client
GOWORK=off go run client.go client_interceptor.go

Protoc
protoc --go_out==plugins=grpc:chat --go-grpc_out=. chat.proto