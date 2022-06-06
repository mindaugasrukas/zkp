FROM alpine:3.14

# install build dependencies
RUN apk add bash go make protoc protobuf-dev; \
    go get -u github.com/golang/protobuf/proto; \
    go get -u github.com/golang/protobuf/protoc-gen-go

ENV PATH=$PATH:/root/go/bin
ENV WORKSPACE=/source

ADD . /source
WORKDIR /source

# Build the application
RUN make client; \
    mkdir /app; \
    cp /source/build/client /app

WORKDIR /app
ENTRYPOINT ["/app/client"]
