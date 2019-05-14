# gRPC to gRPC
This recipe is a proxy gateway for gRPC end points.

## Installation
* Install [Go](https://golang.org/)

## Setup
```bash
git clone https://github.com/project-flogo/grpc
cd grpc/examples/api/grpc-to-grpc
go build
```

## Testing
Start proxy gateway:
```bash
FLOGO_RUNNER_TYPE=DIRECT ./grpc-to-grpc
```

Start sample gRPC server.
```bash
./grpc-to-grpc -server
```

### #1 Testing PetById method
Run sample gRPC client.
```bash
./grpc-to-grpc -client -port 9096 -method pet -param 2
```
Now you should see logs in proxy gateway terminal and sample gRPC server terminal. Sample output in client terminal can be seen as below.
```
res : pet:<id:2 name:"cat2" >
```
### #2 Testing UserByName method
Run sample gRPC client.
```bash
./grpc-to-grpc -client -port 9096 -method user -param user2
```
Output can be seen as below.
```
res : user:<id:2 username:"user2" email:"email2" phone:"phone2" >
```
### #3 Testing ListUsers(Server Streaming) method
Run sample gRPC client.
```bash
./grpc-to-grpc -client -port 9096 -method listusers
```
Streaming data can be observed at client side. Interrupt the sample server to stop streaming.

### #4 Testing StoreUsers(Client Streaming) method
Run sample gRPC client.
```bash
./grpc-to-grpc -client -port 9096 -method storeusers
```
Streaming data can be observed at server side. Interrupt the client to stop streaming.

### #5 Testing BulkUsers(Bidirectional Streaming) method
Run sample gRPC client.
```bash
./grpc-to-grpc -client -port 9096 -method bulkusers
```
Streaming data can be observed at both client and server side. One second delay is kept for veiwing purpose.
