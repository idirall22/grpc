gen:
	protoc -I proto proto/*.proto --go_out=plugins=grpc:.

clean:
	sudo rm -rf pb/*

test:
	go test -cover -race ./...