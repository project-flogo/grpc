package grpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/grpc/metadata"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/project-flogo/core/support/log"
	"google.golang.org/grpc"
)

var clientInterfaceObj interface{}

func (a *Activity) gRPCTogRPCHandler(input *Input, output *Output, logger log.Logger, conn *grpc.ClientConn) error {

	serviceName := input.GRPCMthdParamtrs["serviceName"].(string)
	protoName := input.GRPCMthdParamtrs["protoName"].(string)
	protoName = strings.Split(protoName, ".")[0]

	if len(serviceName) == 0 && len(protoName) == 0 {
		return errors.New("Service name and Proto name required")
	}

	clServFlag := false
	if len(ClientServiceRegistery.ClientServices) != 0 {
		for k, service := range ClientServiceRegistery.ClientServices {
			if strings.Compare(k, protoName+serviceName) == 0 {
				logger.Debugf("client service object found for proto [%v] and service [%v]", protoName, serviceName)
				clientInterfaceObj = service.GetRegisteredClientService(conn)
				clServFlag = true

				methodName := input.GRPCMthdParamtrs["methodName"].(string)
				var throwError bool
				if len(requests) > 0 {
					//For flogo case where all data comes from flogo.json
					md := metadata.New(input.Header)
					input.GRPCMthdParamtrs["contextdata"] = metadata.NewOutgoingContext(context.Background(), md)
					throwError = true
				}
				if input.GRPCMthdParamtrs["contextdata"] != nil {

					inputs := make([]reflect.Value, 2)

					inputs[0] = reflect.ValueOf(input.GRPCMthdParamtrs["contextdata"])
					if reqData, ok := input.GRPCMthdParamtrs["reqdata"]; ok {
						inputs[1] = reflect.ValueOf(reqData)
					} else {
						//No way to handle grpc to grpc
						//Remove all field except data
						inputData := make(map[string]interface{})
						for k, v := range input.GRPCMthdParamtrs {
							if k == "serviceName" || k == "protoName" || k == "contextdata" || k == "reqdata" || k == "methodName" {
								continue
							}
							inputData[k] = v
						}

						request := GetRequest(protoName + "-" + serviceName + "-" + methodName)
						if request != nil {
							err := mapstructure.Decode(inputData, request)
							if err != nil {
								return err
							}
						}
						inputs[1] = reflect.ValueOf(request)
					}

					if len(input.Header) > 0 {
						md := metadata.New(input.Header)
						inputs = append(inputs, reflect.ValueOf(grpc.Header(&md)))
					}
					resultArr := reflect.ValueOf(clientInterfaceObj).MethodByName(methodName).Call(inputs)

					res := resultArr[0]
					grpcErr := resultArr[1]
					if !grpcErr.IsNil() {
						if throwError {
							return fmt.Errorf("%v", grpcErr.Interface())
						} else {
							erroString := fmt.Sprintf("%v", grpcErr.Interface())
							logger.Error("Propagating error to calling function:", erroString)
							erroString = "{\"error\":\"true\",\"details\":{\"error\":\"" + erroString + "\"}}"
							err := json.Unmarshal([]byte(erroString), &output.Body)
							if err != nil {
								return err
							}
						}
					} else {
						output.Body = res.Interface()
					}
				} else {
					InvokeMethodData := make(map[string]interface{})
					InvokeMethodData["ClientObject"] = clientInterfaceObj
					InvokeMethodData["MethodName"] = input.GRPCMthdParamtrs["methodName"]
					InvokeMethodData["reqdata"] = input.GRPCMthdParamtrs["reqdata"]
					InvokeMethodData["strmReq"] = input.GRPCMthdParamtrs["strmReq"]

					resMap := service.InvokeMethod(InvokeMethodData)

					if resMap["Error"] != nil {
						logger.Errorf("Error occured:%v", resMap["Error"])
						erroString := fmt.Sprintf("%v", resMap["Error"])
						erroString = "{\"error\":\"true\",\"details\":{\"error\":\"" + erroString + "\"}}"
						err := json.Unmarshal([]byte(erroString), &output.Body)
						if err != nil {
							return err
						}
					}

				}

			}
		}
		if !clServFlag {
			logger.Errorf("client service object not found for proto [%v] and service [%v]", protoName, serviceName)
		}
	} else {
		logger.Errorf("gRPC Client services not registered")
	}
	return nil
}
