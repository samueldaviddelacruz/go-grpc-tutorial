gen:
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative proto/*.proto --go-grpc_out=pb --go-grpc_opt=paths=source_relative proto/*.proto


clean:
	rm pb/*.go 

server:
	go run cmd/server/main.go -port 7777

client:
	go run cmd/client/main.go -address 0.0.0.0:7777

test:
	go test -cover -race ./...

.PHONY: gen clean server client test