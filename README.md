# ZKP PoC

Good explanation of Chaum-Pedersen Protocol
https://crypto.stackexchange.com/questions/99262/chaum-pedersen-protocol

## Development

### Directory Layout

    client - sample client application 
        app - client application
        cmd - CLI commands
        docker - docker configuration
        
    server - sample server application
        app - server application
        docker - docker configuration

    store - pluggable sample server storage

    zkp - ZKP protocol
        algorithm - ZKP algorithms
        pedersen - Chaum-Pedersen Protocol
        proto - protobuf messages

### Generate dependencies

```shell
$ make proto
```

### Test

```shell
$ make test
```

Get the coverage
```shell
$ make coverage
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

Run client docker container and accessing server container:
```shell
$ docker run -it --rm "zkp-client:0.1" register -s host.docker.internal:8080 -u user-id -p 123
$ docker run -it --rm "zkp-client:0.1" login -s host.docker.internal:8080 -u user-id -p 123
```

Run server using docker-compose:
```shell
$ docker-compose -f server/docker/docker-compose.yml up
```
