# ZKP PoC

Good explanation of Chaum-Pedersen Protocol
https://crypto.stackexchange.com/questions/99262/chaum-pedersen-protocol

## Development

### Generate dependencies

```shell
$ protoc --go_out=./zkp ./zkp/proto/*.proto
```

### Test

```shell
$ go test ./...
```

### Build

```shell
$ go build -o=./build/server ./server/...
$ go build -o=./build/client ./client/main.go
```
