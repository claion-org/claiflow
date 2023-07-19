# README gRPC

## Install protoc

### apt

```sh
apt install -y protobuf-compiler
protoc --version
```

### brew

```sh
brew install protobuf
protoc --version
```

### manual

```sh
https://github.com/protocolbuffers/protobuf/releases
```

## Install grpc

```sh
go get -u google.golang.org/grpc
```

## Install protoc-gen-go

```sh
go get -u google.golang.org/protobuf
go install google.golang.org/protobuf/cmd/protoc-gen-go
protoc-gen-go --version
```

## Install protoc-gen-go-grpc

```sh
go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
protoc-gen-go-grpc --version
```

## Build gRPC

```sh
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative *.proto
```
