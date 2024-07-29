Server
GOWORK=off go run server.go server_interceptor.go

Client
GOWORK=off go run client.go client_interceptor.go

Protoc
protoc --go_out==plugins=grpc:chat --go-grpc_out=. chat.proto

------
Gateway HTTP/1.1

git clone https://github.com/googleapis/googleapis.git
protoc -I. --go_out==plugins=grpc:chat --go-grpc_out=chat --grpc-gateway_out=chat --proto_path=./ --proto_path=./googleapis chat.proto

GOWORK=off go run gateway.go
GOWORK=off go run server.go server_interceptor.go

curl -X POST -k http://localhost:8080/v1/sayhello -d '{"body": "Hello From HTTP/1.1!"}'
