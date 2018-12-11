package examples

import (
	"github.com/project-flogo/core/api"
	"github.com/project-flogo/core/engine"
	"github.com/project-flogo/grpc/activity"
	"github.com/project-flogo/grpc/trigger"
	"github.com/project-flogo/microgateway"
	microapi "github.com/project-flogo/microgateway/api"
)

func GRPC2GRPCExample() (engine.Engine, error) {
	app := api.NewApp()

	gateway := microapi.New("PetStoreDispatch")
	service := gateway.NewService("PetStoregRPCServer", &activity.Activity{})
	service.SetDescription("Make calls to sample PetStore gRPC server")
	service.AddSetting("operatingMode", "grpc-to-grpc")
	service.AddSetting("hosturl", "localhost:9000")
	step := gateway.NewStep(service)
	step.AddInput("grpcMthdParamtrs", "=$.payload.grpcData")
	response := gateway.NewResponse(true)
	response.SetIf("$.PetStoregRPCServer.outputs.body.error == 'true'")
	response.SetCode(404)
	response.SetData("=$.PetStoregRPCServer.outputs.body.details")
	response = gateway.NewResponse(false)
	response.SetCode(200)
	response.SetData("=$.PetStoregRPCServer.outputs.body")
	PetStoreDispatch, err := gateway.AddResource(app)
	if err != nil {
		return nil, err
	}

	triggerSettings := &trigger.Settings{
		Port:        9096,
		ProtoName:   "petstore",
		ServiceName: "PetStoreService",
	}
	trg := app.NewTrigger(&trigger.Trigger{}, triggerSettings)

	handler, err := trg.NewHandler(&trigger.HandlerSettings{
		AutoIdReply:     false,
		UseReplyHandler: false,
	})
	if err != nil {
		return nil, err
	}
	_, err = handler.NewAction(&microgateway.Action{}, PetStoreDispatch)
	if err != nil {
		return nil, err
	}

	e, err := api.NewEngine(app)
	if err != nil {
		return nil, err
	}

	return e, nil
}
