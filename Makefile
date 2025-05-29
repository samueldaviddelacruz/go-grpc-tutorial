gen: 
	protoc \
	--proto_path=proto \
	--go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=openapiv2 \
	proto/auth_service.proto proto/laptop_service.proto


clean:
	rm pb/*.go 

server1:
	go run cmd/server/main.go -port 7778
server2:
	go run cmd/server/main.go -port 7779

server1-tls:
	go run cmd/server/main.go -port 7778 -tls
server2-tls:
	go run cmd/server/main.go -port 7779 -tls

server:
	go run cmd/server/main.go -port 7777

rest:
	go run cmd/server/main.go -port 7776 -type rest

client:
	go run cmd/client/main.go -address 0.0.0.0:8080

client-tls:
	go run cmd/client/main.go -address 0.0.0.0:8080 -tls

test:
	go test -cover -race ./...

cert:
	cd cert; ./gen.sh; cd ..

.PHONY: gen clean server client test cert