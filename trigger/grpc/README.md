# gRPC Trigger
gRPC trigger for mashling gateway supports method name based routing.

## Schema
settings, outputs and handler:

```json
  "settings": [
    {
      "name": "port",
      "type": "integer",
      "required": true
    },
    {
      "name": "protoName",
      "type": "string",
      "required": true
    },
    {
      "name": "protoFile",
      "type": "string"
    },
    {
      "name": "enableTLS",
      "type": "boolean"
    },
    {
      "name": "serverCert",
      "type": "string"
    },
    {
      "name": "serverKey",
      "type": "string"
    }
  ],
  "outputs": [
    {
      "name": "params",
      "type": "object"
    },
    {
      "name": "grpcData",
      "type": "object"
    },
    {
      "name": "content",
      "type": "any"
    }
  ],
  "handler": {
    "settings": [
      {
        "name": "serviceName",
        "type": "string"
      },
      {
        "name": "methodName",
        "type": "string"
      }
    ]
  }
```
### Settings
| Key    | Description   |
|:-----------|:--------------|
| port | The port to listen on |
| protoName | The name of the proto file|
| protoFile| The content of the proto file|
| enableTLS | true - To enable TLS (Transport Layer Security), false - No TLS security  |
| serverCert | Server certificate file in PEM format. Need to provide file name along with path. Path can be relative to gateway binary location. |
| serverKey | Server private key file in PEM format. Need to provide file name along with path. Path can be relative to gateway binary location. |

### Outputs
| Key    | Description   |
|:-----------|:--------------|
| params | Request params |
| content | HTTP request payload |
| grpcData | gRPC Method parameters |

### Handler settings
| Key    | Description   |
|:-----------|:--------------|
| serviceName | The name of the service mentioned in proto file|
| methodName | Name of the method |


### Sample Mashling Gateway Recipie

Following is the example mashling gateway descriptor uses a grpc trigger.

```json
{
  "name": "MyProxy",
  "type": "flogo:app",
  "version": "1.0.0",
  "description": "This is a simple proxy.",
  "properties": null,
  "channels": null,
  "triggers": [
    {
      "name": "flogo-grpc",
      "id": "MyProxy",
      "ref": "github.com/project-flogo/grpc/trigger/grpc",
      "settings": {
        "port": "9096",
        "protoName": "petstore"
      },
      "handlers": [
        {
          "settings": {
            "serviceName": "GRPC2RestPetStoreService",
            "methodName": "PetById"
          },
          "actions": [
            {
              "id": "microgateway:petByIdDispatch"
            }
          ]
        },
        {
          "settings": {
            "methodName": "PetPUT"
          },
          "actions": [
            {
              "id": "microgateway:petPutDispatch"
            }
          ]
        },
        {
          "settings": {
            "methodName": "UserByName"
          },
          "actions": [
            {
              "id": "microgateway:userByNameDispatch"
            }
          ]
        }
      ]
    }
  ],
  "resources": [
    {
      "id": "microgateway:petByIdDispatch",
      "compressed": false,
      "data": {
        "name": "Pets",
        "steps": [
          {
            "service": "PetStorePets",
            "input": {
              "pathParams.id": "=$.payload.params.id"
            }
          }
        ],
        "responses": [
          {
            "if": "$.PetStorePets.outputs.data.type == 'error'",
            "error": true,
            "output": {
              "code": 404,
              "data": {
                "error": "=$.PetStorePets.outputs.data.message"
              }
            }
          },
          {
            "if": "$.PetStorePets.outputs.data.type != 'error'",
            "error": false,
            "output": {
              "code": 200,
              "data": {
                "pet": "=$.PetStorePets.outputs.data"
              }
            }
          }
        ],
        "services": [
          {
            "name": "PetStorePets",
            "description": "Make calls to find pets",
            "ref": "github.com/project-flogo/contrib/activity/rest",
            "settings": {
              "uri": "http://petstore.swagger.io/v2/pet/:id",
              "method": "GET"
            }
          }
        ]
      }
    },
    {
      "id": "microgateway:petPutDispatch",
      "compressed": false,
      "data": {
        "name": "Pets",
        "steps": [
          {
            "service": "PetStorePetsPUT",
            "input": {
              "content": "=$.payload.content"
            }
          }
        ],
        "responses": [
          {
            "if": "$.PetStorePetsPUT.outputs.data.type == 'error'",
            "error": true,
            "output": {
              "code": 404,
              "data": {
                "error": "=$.PetStorePetsPUT.outputs.data.message"
              }
            }
          },
          {
            "if": "$.PetStorePetsPUT.outputs.data.type != 'error'",
            "error": false,
            "output": {
              "code": 200,
              "data": {
                "pet": "=$.PetStorePetsPUT.outputs.data"
              }
            }
          }
        ],
        "services": [
          {
            "name": "PetStorePetsPUT",
            "description": "Make calls to petstore",
            "ref": "github.com/project-flogo/contrib/activity/rest",
            "settings": {
              "uri": "http://petstore.swagger.io/v2/pet",
              "method": "PUT",
              "headers": {
                "Content-Type": "application/json"
              }
            }
          }
        ]
      }
    },
    {
      "id": "microgateway:userByNameDispatch",
      "compressed": false,
      "data": {
        "name": "Pets",
        "steps": [
          {
            "service": "PetStoreUsersByName",
            "input": {
              "pathParams.username": "=$.payload.params.username"
            }
          }
        ],
        "responses": [
          {
            "if": "$.PetStoreUsersByName.outputs.data.type == 'error'",
            "error": true,
            "output": {
              "code": 404,
              "data": {
                "error": "=$.PetStoreUsersByName.outputs.data.message"
              }
            }
          },
          {
            "if": "$.PetStoreUsersByName.outputs.data.type != 'error'",
            "error": false,
            "output": {
              "code": 200,
              "data": {
                "user": "=$.PetStoreUsersByName.outputs.data"
              }
            }
          }
        ],
        "services": [
          {
            "name": "PetStoreUsersByName",
            "description": "Make calls to find users",
            "ref": "github.com/project-flogo/contrib/activity/rest",
            "settings": {
              "uri": "http://petstore.swagger.io/v2/user/:username",
              "method": "GET"
            }
          }
        ]
      }
    }
  ],
  "actions": [
    {
      "ref": "github.com/project-flogo/microgateway",
      "settings": {
        "uri": "microgateway:petByIdDispatch"
      },
      "id": "microgateway:petByIdDispatch",
      "metadata": null
    },
    {
      "ref": "github.com/project-flogo/microgateway",
      "settings": {
        "uri": "microgateway:petPutDispatch"
      },
      "id": "microgateway:petPutDispatch",
      "metadata": null
    },
    {
      "ref": "github.com/project-flogo/microgateway",
      "settings": {
        "uri": "microgateway:userByNameDispatch"
      },
      "id": "microgateway:userByNameDispatch",
      "metadata": null
    }
  ]
}
```
### Sample Usage
This trigger depends on support files which can be generated with the grpc tool by passing proto file.

Sample demonstration of this trigger can be found in [examples](examples/).

#### Note
Currently This Trigger handles.<br>
1. Unary methods propagation.
2. REST path/query params can be mapped through params output key.
3. Routing can be done based on method names.
