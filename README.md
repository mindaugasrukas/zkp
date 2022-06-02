# ZKP PoC

## Development

### Generate dependencies

```shell
$ protoc --go_out=./zkp ./zkp/*.proto
```

### Test

```shell
$ go test ./...
```

### Build

```shell
$ go build -o=./build/server ./server/...
$ go build -o=./build/client ./client/...
```
