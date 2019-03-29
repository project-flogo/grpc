package examples

import (
	"github.com/project-flogo/contrib/activity/rest"
	resttrigger "github.com/project-flogo/contrib/trigger/rest"
	"github.com/project-flogo/core/api"
	"github.com/project-flogo/core/engine"
	"github.com/project-flogo/grpc/activity"
	"github.com/project-flogo/grpc/trigger"
	"github.com/project-flogo/microgateway"
	microapi "github.com/project-flogo/microgateway/api"
)

var petStoreURL = "http://petstore.swagger.io/v2"

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

	handler, err := trg.NewHandler(&trigger.HandlerSettings{})
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

func GRPC2RestExample() (engine.Engine, error) {
	app := api.NewApp()

	gateway := microapi.New("petByIdDispatch")
	service := gateway.NewService("PetStorePets", &rest.Activity{})
	service.SetDescription("Make calls to find pets")
	service.AddSetting("uri", petStoreURL+"/pet/:id")
	service.AddSetting("method", "GET")
	step := gateway.NewStep(service)
	step.AddInput("pathParams.id", "=$.payload.params.Id")
	response := gateway.NewResponse(true)
	response.SetIf("$.PetStorePets.outputs.data.type == 'error'")
	response.SetCode(404)
	response.SetData(map[string]interface{}{
		"error": "=$.PetStorePets.outputs.data.message",
	})
	response = gateway.NewResponse(false)
	response.SetIf("$.PetStorePets.outputs.data.type != 'error'")
	response.SetCode(200)
	response.SetData(map[string]interface{}{
		"pet": "=$.PetStorePets.outputs.data",
	})
	petByIdDispatch, err := gateway.AddResource(app)
	if err != nil {
		return nil, err
	}

	gateway = microapi.New("petPutDispatch")
	service = gateway.NewService("PetStorePetsPUT", &rest.Activity{})
	service.SetDescription("Make calls to petstore")
	service.AddSetting("uri", petStoreURL+"/pet")
	service.AddSetting("method", "PUT")
	service.AddSetting("headers", map[string]string{"Content-Type": "application/json"})
	step = gateway.NewStep(service)
	step.AddInput("content", "=$.payload.content")
	response = gateway.NewResponse(true)
	response.SetIf("$.PetStorePetsPUT.outputs.data.type == 'error'")
	response.SetCode(404)
	response.SetData(map[string]interface{}{
		"error": "=$.PetStorePetsPUT.outputs.data.message",
	})
	response = gateway.NewResponse(false)
	response.SetIf("$.PetStorePetsPUT.outputs.data.type != 'error'")
	response.SetCode(200)
	response.SetData(map[string]interface{}{
		"pet": "=$.PetStorePetsPUT.outputs.data",
	})
	petPutDispatch, err := gateway.AddResource(app)
	if err != nil {
		return nil, err
	}

	gateway = microapi.New("userByNameDispatch")
	service = gateway.NewService("PetStoreUsersByName", &rest.Activity{})
	service.SetDescription("Make calls to find users")
	service.AddSetting("uri", petStoreURL+"/user/:username")
	service.AddSetting("method", "GET")
	step = gateway.NewStep(service)
	step.AddInput("pathParams.username", "=$.payload.params.Username")
	response = gateway.NewResponse(true)
	response.SetIf("$.PetStoreUsersByName.outputs.data.type == 'error'")
	response.SetCode(404)
	response.SetData(map[string]interface{}{
		"error": "=$.PetStoreUsersByName.outputs.data.message",
	})
	response = gateway.NewResponse(false)
	response.SetIf("$.PetStoreUsersByName.outputs.data.type != 'error'")
	response.SetCode(200)
	response.SetData(map[string]interface{}{
		"user": "=$.PetStoreUsersByName.outputs.data",
	})
	userByNameDispatch, err := gateway.AddResource(app)
	if err != nil {
		return nil, err
	}

	triggerSettings := &trigger.Settings{
		Port:        9096,
		ProtoName:   "petstore",
		ServiceName: "GRPC2RestPetStoreService",
	}
	trg := app.NewTrigger(&trigger.Trigger{}, triggerSettings)

	action := &microgateway.Action{}

	handler, err := trg.NewHandler(&trigger.HandlerSettings{
		MethodName: "PetById",
	})
	if err != nil {
		return nil, err
	}
	_, err = handler.NewAction(action, petByIdDispatch)
	if err != nil {
		return nil, err
	}

	handler, err = trg.NewHandler(&trigger.HandlerSettings{
		MethodName: "PetPUT",
	})
	if err != nil {
		return nil, err
	}
	_, err = handler.NewAction(action, petPutDispatch)
	if err != nil {
		return nil, err
	}

	handler, err = trg.NewHandler(&trigger.HandlerSettings{
		MethodName: "UserByName",
	})
	if err != nil {
		return nil, err
	}
	_, err = handler.NewAction(action, userByNameDispatch)
	if err != nil {
		return nil, err
	}

	return api.NewEngine(app)
}

func Rest2GRPCExample() (engine.Engine, error) {
	app := api.NewApp()

	gateway := microapi.New("PetStoreGETDispatch")
	service := gateway.NewService("PetStore", &activity.Activity{})
	service.SetDescription("Make calls to PetStore grpc end point")
	service.AddSetting("operatingMode", "rest-to-grpc")
	service.AddSetting("hosturl", "localhost:9000")
	step := gateway.NewStep(service)
	step.AddInput("protoName", "petstore")
	step.AddInput("serviceName", "Rest2GRPCPetStoreService")
	step.AddInput("methodName", "=$.payload.pathParams.grpcMethodName")
	step.AddInput("queryParams", "=$.payload.queryParams")
	response := gateway.NewResponse(true)
	response.SetIf("$.PetStore.outputs.body.error == 'true'")
	response.SetCode(404)
	response.SetData("=$.PetStore.outputs.body.details")
	response = gateway.NewResponse(false)
	response.SetCode(200)
	response.SetData("=$.PetStore.outputs.body")
	PetStoreGETDispatch, err := gateway.AddResource(app)
	if err != nil {
		return nil, err
	}

	gateway = microapi.New("PetStorePUTDispatch")
	service = gateway.NewService("PetStore", &activity.Activity{})
	service.SetDescription("Make calls to PetStore grpc end point")
	service.AddSetting("operatingMode", "rest-to-grpc")
	service.AddSetting("hosturl", "localhost:9000")
	step = gateway.NewStep(service)
	step.AddInput("protoName", "petstore")
	step.AddInput("serviceName", "Rest2GRPCPetStoreService")
	step.AddInput("methodName", "PetPUT")
	step.AddInput("content", "=$.payload.content")
	response = gateway.NewResponse(true)
	response.SetIf("$.PetStore.outputs.body.error == 'true'")
	response.SetCode(404)
	response.SetData("=$.PetStore.outputs.body.details")
	response = gateway.NewResponse(false)
	response.SetCode(200)
	response.SetData("=$.PetStore.outputs.body")
	PetStorePUTDispatch, err := gateway.AddResource(app)
	if err != nil {
		return nil, err
	}

	trg := app.NewTrigger(&resttrigger.Trigger{}, &resttrigger.Settings{Port: 9096})

	handler, err := trg.NewHandler(&resttrigger.HandlerSettings{
		Method: "GET",
		Path:   "/petstore/method/:grpcMethodName",
	})
	if err != nil {
		return nil, err
	}
	_, err = handler.NewAction(&microgateway.Action{}, PetStoreGETDispatch)
	if err != nil {
		return nil, err
	}

	handler, err = trg.NewHandler(&resttrigger.HandlerSettings{
		Method: "PUT",
		Path:   "/petstore/PetPUT",
	})
	if err != nil {
		return nil, err
	}
	_, err = handler.NewAction(&microgateway.Action{}, PetStorePUTDispatch)
	if err != nil {
		return nil, err
	}

	return api.NewEngine(app)
}
