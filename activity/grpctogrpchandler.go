package activity

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"google.golang.org/grpc"

	"github.com/project-flogo/core/support/log"
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

				if input.GRPCMthdParamtrs["contextdata"] != nil {

					inputs := make([]reflect.Value, 2)

					inputs[0] = reflect.ValueOf(input.GRPCMthdParamtrs["contextdata"])
					inputs[1] = reflect.ValueOf(input.GRPCMthdParamtrs["reqdata"])

					resultArr := reflect.ValueOf(clientInterfaceObj).MethodByName(input.GRPCMthdParamtrs["methodName"].(string)).Call(inputs)

					res := resultArr[0]
					grpcErr := resultArr[1]
					if !grpcErr.IsNil() {
						erroString := fmt.Sprintf("%v", grpcErr.Interface())
						logger.Error("Propagating error to calling function:", erroString)
						erroString = "{\"error\":\"true\",\"details\":{\"error\":\"" + erroString + "\"}}"
						err := json.Unmarshal([]byte(erroString), &output.Body)
						if err != nil {
							return err
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
