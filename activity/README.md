#### gRPC

The `grpc` service type works as gRPC client and will communicate to the given address with given method parameters.

The service `settings` for the request are as follows:

| Name   |  Type   | Description   |
|:-----------|:--------|:--------------|
| operatingMode | string | Either 'grpc-to-grpc' or 'rest-to-grpc' |
| hosturl | string | A gRPC end point url with port |
| enableTLS | bool | true - To enable TLS (Transport Layer Security), false - No TLS security  |
| clientCert | string | Server certificate file in PEM format. Need to provide file name along with path. Path can be relative to gateway binary location. |

The available `input` for the request are as follows:

| Name   |  Type   | Description   |
|:-----------|:--------|:--------------|
| grpcMthdParamtrs | JSON object | A grpcMthdParamtrs payload which holds full information like method parameters, service name, proto name, method name etc.|
| header | JSON object | HTTP request header params|
| serviceName | string | Name of the service present in proto |
| protoName | string | Name of the proto file used |
| methodName | string | rpc method name present inside service |
| params | JSON object | HTTP request params |
| queryParams | JSON object | HTTP request query params |
| content | JSON object | HTTP request paylod |
| pathParams | JSON object | HTTP request path params |

The available response `outputs` are as follows:

| Name   |  Type   | Description   |
|:-----------|:--------|:--------------|
|body | JSON object | The response object from gRPC end server |

A sample `service` definition is:

```json
{
    "name": "PetStoreUsers",
    "description": "Make calls to grpc end point",
    "ref": "github.com/project-flogo/grpc/activity",
    "settings": {
        "hosturl": "localhost:9000",
        "enableTLS": "true"
    }
}
```

An example `step` that invokes the above `PetStoreUsers` service using `grpcMthdParamtrs` is:

```json
{
 "service": "PetStoreUsers",
 "input": {
    "grpcMthdParamtrs": "=$.payload.grpcData",
    "clientCert": "=$.env.CLIENT_CERT"
 }
}
```

Response handler:

```json
{
  "error": false,
  "output": {
      "code": 200,
      "data": "=$.PetStoreUsers.response.body"
  }
}
```

#### Note
Support files for this service are generated using proto file with grpc command. Unary methods are allowed in all grpc gateway recipes. Streaming methods are allowed only in case of grpc-to-grpc gateway.
