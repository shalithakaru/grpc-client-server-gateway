Server
GOWORK=off go run server.go server_interceptor.go

Client
GOWORK=off go run client.go client_interceptor.go

Protoc
protoc --go_out==plugins=grpc:chat --go-grpc_out=. chat.proto

------
after the gateway

git clone https://github.com/googleapis/googleapis.git
protoc -I. --go_out==plugins=grpc:chat --go-grpc_out=chat --grpc-gateway_out=chat --proto_path=./ --proto_path=./googleapis chat.proto
