# REST to gRPC
This recipe demonstrates on receiving request from a Rest client and routing to gRPC server.

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
cd examples/json/rest-to-grpc
```

Create the gateway:
```bash
flogo create -f flogo.json
cd MyProxy
flogo install github.com/project-flogo/proto/rest2grpc
flogo build
```

## Testing
Start gateway.
```bash
bin/MyProxy
```

Start sample gRPC server.
```bash
go run main.go -server
```

### #1 Testing PetById method with GET request
Sample GET request.
```curl
curl --request GET http://localhost:9096/petstore/method/PetById?id=2
```
Now you should see logs in gateway terminal and sample gRPC server terminal. Output in curl request terminal can be seen as below.
```json
{
 "pet": {
  "id": 2,
  "name": "cat2"
 }
}
```
### #2 Testing UserByName method with GET request
Sample GET request.
```curl
curl --request GET http://localhost:9096/petstore/method/UserByName?username=user2
```
Output can be seen as below.
```json
{
 "user": {
  "email": "email2",
  "id": 2,
  "phone": "phone2",
  "username": "user2"
 }
}
```
### #3 Testing PetPUT method with PUT request
Payload
```json
{
 "pet": {
  "id": 12,
  "name": "mycat12"
 }
}
```
Curl command
```curl
curl -X PUT "http://localhost:9096/petstore/PetPUT" -H "accept: application/xml" -H "Content-Type: application/json" -d '{"pet": {"id": 12,"name": "mycat12"}}'
```
Output can be seen as below.
```json
{
 "pet": {
  "id": 12,
  "name": "mycat12"
 }
}
```
