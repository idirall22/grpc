gen:
	protoc -I proto proto/*.proto --go_out=plugins=grpc:.

clean:
	sudo rm -rf pb/*

test:
	go test -cover -race ./...

server:
	go run cmd/server/main.go -port 8080

client:
	go run cmd/client/main.go -address 0.0.0.0:8080 