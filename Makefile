.PHONY: all
all: server client

.PHONY: proto
proto:
	protoc --go_out=./zkp ./zkp/proto/*.proto

.PHONY: client
client: proto
	go build -o=./build/client ./client/main.go

.PHONY: server
server: proto
	go build -o=./build/server ./server/...

.PHONY: clean
clean:
	rm -rf ./build
	rm -rf ./zkp/gen

.PHONY: test
test: all
	go test ./...
