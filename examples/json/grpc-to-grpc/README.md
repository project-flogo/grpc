# gRPC to gRPC
This recipe is a proxy gateway for gRPC end points.

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
cd examples/json/grpc-to-grpc
```

Create the gateway:
```bash
flogo create -f flogo.json
cd MyProxy
flogo install github.com/project-flogo/grpc/proto/grpc2grpc
flogo build
```

## Testing
Start proxy gateway.
```bash
bin/MyProxy
```

Start sample gRPC server.
```bash
go run main.go -server
```

### #1 Testing PetById method
Run sample gRPC client.
```bash
go run main.go -client -port 9096 -method pet -param 2
```
Now you should see logs in proxy gateway terminal and sample gRPC server terminal. Sample output in client terminal can be seen as below.
```
res : pet:<id:2 name:"cat2" >
```
### #2 Testing UserByName method
Run sample gRPC client.
```bash
go run main.go -client -port 9096 -method user -param user2
```
Output can be seen as below.
```
res : user:<id:2 username:"user2" email:"email2" phone:"phone2" >
```
### #3 Testing ListUsers(Server Streaming) method
Run sample gRPC client.
```bash
go run main.go -client -port 9096 -method listusers
```
Streaming data can be observed at client side. Interrupt the sample server to stop streaming.

### #4 Testing StoreUsers(Client Streaming) method
Run sample gRPC client.
```bash
go run main.go -client -port 9096 -method storeusers
```
Streaming data can be observed at server side. Interrupt the client to stop streaming.

### #5 Testing BulkUsers(Bidirectional Streaming) method
Run sample gRPC client.
```bash
go run main.go -client -port 9096 -method bulkusers
```
Streaming data can be observed at both client and server side. One second delay is kept for veiwing purpose.
