.PHONY: all install_deps proto server gateway client

all: install_deps proto server gateway client client-web

# Install necessary tools and dependencies
install_deps:
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	git clone https://github.com/googleapis/googleapis.git
	npm install -g grpc-tools
	npm install -g protoc-gen-grpc-web

# Generate protobuf files
proto:
	protoc -I. --go_out=chat --go_opt=paths=source_relative \
		--go-grpc_out=chat --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=chat --grpc-gateway_opt=logtostderr=true,paths=source_relative \
		--proto_path=./ --proto_path=./googleapis chat.proto
	protoc -I=. chat.proto \
		--js_out=import_style=commonjs,binary:./client-web/src/grpc \
		--grpc-web_out=import_style=commonjs,mode=grpcwebtext:./client-web/src/grpc --proto_path=./googleapis

# Run the gRPC server
server:
	GOWORK=off go run server/*.go

# Run the gRPC gateways
gateway: gateway-grpc-http1-1 gateway-http-http1-1 gateway-grpc-http2 gateway-http-http2

gateway-grpc-http1-1:
	GOWORK=off go run gateway/gateway-grpc-http1-1/*.go

gateway-http-http1-1:
	GOWORK=off go run gateway/gateway-http-http1-1/*.go

gateway-grpc-http2:
	GOWORK=off go run gateway/gateway-grpc-http2/*.go

gateway-http-http2:
	GOWORK=off go run gateway/gateway-http-http2/*.go

# Run the gRPC client
client:
	GOWORK=off go run client/*.go --http-gateway
	GOWORK=off go run client/*.go --grpc-gateway
	GOWORK=off go run client/*.go --http-gateway-http2
	GOWORK=off go run client/*.go --grpc-gateway-http2
	GOWORK=off go run client/*.go

# Run the client-web application
client-web: client-web-setup client-web-proxy client-web-run

client-web-setup:
	cd client-web && npm install

client-web-proxy:
	git clone https://github.com/improbable-eng/grpc-web.git
	go install ./go/grpcwebproxy
	grpcwebproxy --backend_addr=localhost:9000 --run_tls_server=false --allow_all_origins &

client-web-run:
	cd client-web && npm start
