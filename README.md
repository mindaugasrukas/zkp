# ZKP PoC

Good explanation of Chaum-Pedersen Protocol
https://crypto.stackexchange.com/questions/99262/chaum-pedersen-protocol

## Development

### Generate dependencies

```shell
$ make proto
```

### Test

```shell
$ make test
```

### Build OS native application

```shell
$ make server
$ make client
```

### Build docker images

```shell
$ make server-image
$ make client-image
```

### Run docker images

Run server with default setting
```shell
$ make server-run
```

Run client docker container and accessing server container
```shell
$ docker run -it --rm "zkp-client:0.1" register -s host.docker.internal:8080 -u user-id -p 123
$ docker run -it --rm "zkp-client:0.1" login -s host.docker.internal:8080 -u user-id -p 123
```
