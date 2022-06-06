WORKSPACE = $(shell git rev-parse --show-toplevel)
SERVER-VERSION = 0.1
CLIENT-VERSION = 0.1

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
	go build -o=./build/server ./server/main.go

.PHONY: clean
clean:
	rm -rf ./build
	rm -rf ./zkp/gen

.PHONY: test
test: all
	go test ./...

.PHONY: server-image
server-image: server
	docker build -t "zkp-server:$(SERVER-VERSION)" -f "server/docker/Dockerfile" $(WORKSPACE)

.PHONY: client-image
client-image: client
	docker build -t "zkp-client:$(CLIENT-VERSION)" -f "client/docker/Dockerfile" $(WORKSPACE)

.PHONY: run-server
run-server: server-image
	docker run -it --rm -p 8080:8080 "zkp-server:$(SERVER-VERSION)"
