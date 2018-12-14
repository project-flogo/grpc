# REST to gRPC
This recipe demonstrates on receiving request from a Rest client and routing to gRPC server.

## Installation
* Install [Go](https://golang.org/)

## Setup
```bash
git clone https://github.com/project-flogo/grpc
cd grpc/examples/api/rest-to-grpc
go build
```

## Testing
Start gateway.
```bash
./rest-to-grpc
```

Start sample gRPC server.
```bash
./rest-to-grpc -server
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
