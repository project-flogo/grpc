# gRPC to REST
This recipe demonstrates receiving request from a gRPC client and routing to REST end point based on method names.

## Installation
* Install [Go](https://golang.org/)
* Install `protoc-gen-go` library
```bash
go get github.com/golang/protobuf/protoc-gen-go
```
* Download protoc for your respective OS from [here](https://github.com/google/protobuf/releases).<br>Extract protoc-$VERSION-$PLATFORM.zip file get the `protoc` binary from bin folder and configure it in PATH.

## Setup
```bash
git clone https://github.com/project-flogo/grpc
cd grpc
go install
cd examples/json/grpc-to-rest
```

Create the gateway:
```bash
flogo create -f flogo.json
cd MyProxy
flogo install github.com/project-flogo/proto/grpc2rest
flogo build
```

## Testing
Start proxy gateway.
```bash
bin/MyProxy
```

### #1 Testing PetById method
Run sample gRPC client.
```bash
go run main.go -client -port 9096 -method pet -param 2
```
Output can be seen as below.
```
res : pet:<id:2 name:"cat2" >
```
### #2 Testing UserByName method
Run sample gRPC client.
```bash
go run main.go -client -port 9096 -method user -param user1
```
Output can be seen as below.
```
res : user:<id:1 username:"user1" email:"email1@test.com" phone:"123-456-7890" >
```
### #3 Testing PetPUT method
Run sample gRPC client.
```bash
go run main.go -client -port 9096 -method petput -param 2,testpet
```
Output can be seen as below.
```
res : user:<id:1 username:"user1" email:"email1@test.com" phone:"123-456-7890" >
```
