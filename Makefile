gen:
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative proto/*.proto --go-grpc_out=pb --go-grpc_opt=paths=source_relative proto/*.proto


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

client:
	go run cmd/client/main.go -address 0.0.0.0:8080

client-tls:
	go run cmd/client/main.go -address 0.0.0.0:8080 -tls

test:
	go test -cover -race ./...

cert:
	cd cert; ./gen.sh; cd ..

.PHONY: gen clean server client test cert