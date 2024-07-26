# Inventive Weave

Backend of an application that allows products to be created and shared with the world!

## 1. Install tools

### Go toolchain

#### Mac

Check whether you have this tool installed using:
```bash
go version
$ go version go1.22.4 darwin/arm64
```

If you get an error, install it using:
```bash
brew install go
```

### Protocol buffer compiler

#### Mac

Check whether you have this tool installed using:
```bash
protoc --version
$ libprotoc 27.1
```

If you get an error, install it using:
```bash
brew install protobuf
```

### Go protocol buffers protoc plugin

#### Mac

Check whether you have this tool installed using:
```bash
protoc-gen-go --version
$ protoc-gen-go v1.34.2
```

If you get an error, install it using:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

### Go gRPC protoc plugin

#### Mac

Check whether you have this tool installed using:
```bash
protoc-gen-go-grpc --version
$ protoc-gen-go-grpc 1.4.0
```

If you get an error, install it using:
```bash
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## 2. Generate files

If you modify .proto files, the related .go files need to be regenerated. Do this with:
```bash
go generate ./...
```

## 3. Run tests

Run the test suite with:
```bash
go test ./...
```

## Dependencies

This project uses some dependencies. They have been carefully selected. Dependencies are only added when needed. The cost of adding a dependency is valued 
higher in Go than in many other languages. Only add dependencies when they save a large amount of implementation and 
when a significant portion of the library is relevant to this project.

* `google.golang.org/grpc`
Official gRPC library.
* `google.golang.org/protobuf`
Official protobuf library.
* `github.com/stretchr/testify`
Very widely used testing helper library. Provides assertions which are less verbose than if statements, especially for 
collections, and provides nice failure output for easier debugging.

// todo document architectural decisions (no db, etc)
// todo add diagram
// todo describe project structure
// todo docker compose