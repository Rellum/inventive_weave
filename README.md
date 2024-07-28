# Inventive Weave

Demo of an application that allows products to be created and shared with the world!

### Architecture overview

1. Written in Golang (fast compiling, type safety, small deployables, performant concurrency)
1. Microservice design (independently deploy and scale components)
1. gRPC connections between services (api-first, type safety, high performance, code generation)
1. Monorepo (single version stream, easier code sharing, smaller dev and tooling footprint)
1. No data persistence yet (defer data storage design until required)
1. Prometheus for metrics (industry standard, possibility of self-hosting, safer pull strategy)
1. Structured logging (easier indexing, aggregation, machine readable)
1. Multi-stage container build (smaller images, reduces storage requirements on nodes)
1. Static files embedded into Go binaries (single file deployables)
1. Additional architecture decisions are not made at this stage:
   1. Production platform (possibilities include Kubernetes or a managed solution)
   1. Internal networking security (server certs, possibly client certs, possibly service mesh network)
   1. Frontend (website, web apis, public apis, mobile apis if needed)
   1. Frontend architecture (servers, load balancing, auth, security, etc)
   1. CI/CD (recommend continuous deployment. Needs some process and ownership, and a rollback plan)

### Local development

#### Install tools

Instructions for Mac. Linux or Windows commands will differ - use Google and your favourite package manager.

*Go toolchain*

Check whether you have this tool installed using:
```bash
go version
$ go version go1.22.4 darwin/arm64
```

If you get an error, install it using:
```bash
brew install go
```

*Protocol buffer compiler*

Check whether you have this tool installed using:
```bash
protoc --version
$ libprotoc 27.1
```

If you get an error, install it using:
```bash
brew install protobuf
```

*Go protocol buffers protoc plugin*

Check whether you have this tool installed using:
```bash
protoc-gen-go --version
$ protoc-gen-go v1.34.2
```

If you get an error, install it using:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

*Go gRPC protoc plugin*

Check whether you have this tool installed using:
```bash
protoc-gen-go-grpc --version
$ protoc-gen-go-grpc 1.4.0
```

If you get an error, install it using:
```bash
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

*Docker*

Different options are possible. Look into `podman`, `Docker desktop`, `rancher`
and `colima`

Make sure that both `docker` and `docker compose` commands are supported. You may need to install command line clients 
separately.

#### Generate files

If you modify .proto files, the related .go files need to be regenerated. Do this with:
```bash
go generate ./...
```

#### Run tests

Run the whole test suite with:
```bash
go test ./...
```

#### Run backend with docker compose:

You can build and run the whole backend with:
```bash
docker-compose up --build
```

Then you will be able to reach the following endpoints:
1. http://localhost:9071/metrics (the prometheus metrics of the creators service)
1. http://localhost:9073/metrics (the prometheus metrics of the demo fe service)
1. http://localhost:9071 (a _very_ basic demo form to upload creator data)

With the backend running, you can also use the `creatorclient` cli tool to test the creator service. For example:
```bash
cat ./data/example1.json | go run ./tools/creatorclient --creators_svc localhost:9070
```

Shut down the backend with `^c` in the terminal running docker compose, or with:
```bash
docker-compose down
```

#### Dependencies

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

#### Project structure

```
 |-tools                    // cli applications 
 |-svc                      // services
 | |-servicename            // root has main package of server
 | | |-doc.go               // service-level documentation goes here
 | | |-static               // static files and assets
 | | |-templates            // go templates
 | | |-types                // types related to this service. Can be imported from outside of the service.
 | | |-server               // gRPC server implementation
 | | |-servicenamepb        // gRPC and protobuf definitions
 | | | |-generate.go        // implements go generate command
 | | | |-servicename.proto  // definitions
 | | | |-servicename.pb.go  // generated files are checked in. Generate toolchain not needed for general Go development.
 | | | |-convert.go         // Code to convert between the proto and normal Go types
 |-data                     // Data files used in multiple tests and examples
 |-pkg                      // Go library code. Organised in packages. No "util" or mixed bag packages.

```

// todo add diagram