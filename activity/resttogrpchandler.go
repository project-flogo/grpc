package activity

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/grpc"

	"github.com/project-flogo/core/support/log"
)

func (a *Activity) restTogRPCHandler(input *Input, output *Output, logger log.Logger, conn *grpc.ClientConn) error {
	// check for method name
	if len(input.MethodName) == 0 {
		if len(input.PathParams["grpcMethodName"]) == 0 {
			logger.Error("Method name not provided in json/pathParams")
			return errors.New("Method name not provided")
		}
		input.MethodName = input.PathParams["grpcMethodName"]
		logger.Debug("Method name: ", input.MethodName)
	}

	clServFlag := false
	if len(ClientServiceRegistery.ClientServices) != 0 {
		for k, service := range ClientServiceRegistery.ClientServices {
			if strings.Compare(k, input.ProtoName+input.ServiceName) == 0 {
				logger.Debugf("client service object found for proto [%v] and service [%v]", input.ProtoName, input.ServiceName)

				InvokeMethodData := make(map[string]interface{})
				InvokeMethodData["ClientObject"] = service.GetRegisteredClientService(conn)
				InvokeMethodData["MethodName"] = input.MethodName
				if len(input.PathParams) != 0 {
					InvokeMethodData["PathParams"] = input.PathParams
				}
				if len(input.Params) != 0 {
					InvokeMethodData["Params"] = input.Params
				}
				if len(input.QueryParams) != 0 {
					InvokeMethodData["QueryParams"] = input.QueryParams
				}
				if input.Content != nil {
					InvokeMethodData["Content"] = input.Content
				}
				InvokeMethodData["Mode"] = "rest-to-grpc"
				resMap := service.InvokeMethod(InvokeMethodData)
				if resMap["Response"] != nil && strings.Compare(string(resMap["Response"].([]byte)), "null") != 0 {
					err := json.Unmarshal(resMap["Response"].([]byte), &output.Body)
					if err != nil {
						return err
					}
				} else {
					erroString := fmt.Sprintf("%v", resMap["Error"])
					erroString = "{\"error\":\"true\",\"details\":{\"error\":\"" + erroString + "\"}}"
					err := json.Unmarshal([]byte(erroString), &output.Body)
					if err != nil {
						return err
					}
				}
				clServFlag = true
			}
		}
		if !clServFlag {
			logger.Errorf("client service object not found for proto [%v] and service [%v]", input.ProtoName, input.ServiceName)
		}
	} else {
		logger.Errorf("gRPC Client services not registered")
	}

	return nil
}
