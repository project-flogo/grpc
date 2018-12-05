package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/project-flogo/contrib/activity/rest"
	"github.com/project-flogo/core/api"
	"github.com/project-flogo/core/engine"
	"github.com/project-flogo/grpc/proto/grpc2rest"
	trigger "github.com/project-flogo/grpc/trigger"
	"github.com/project-flogo/microgateway"
	microapi "github.com/project-flogo/microgateway/api"
	"google.golang.org/grpc"
)

var addr string
var clientAddr string

func main() {
	clientMode := flag.Bool("client", false, "command to run client")
	methodName := flag.String("method", "pet", "method name")
	paramVal := flag.String("param", "user2", "method param")
	port := flag.String("port", "", "port value")

	flag.Parse()
	var id int

	if *clientMode {
		if strings.Compare(*methodName, "pet") == 0 {
			id, _ = strconv.Atoi(*paramVal)
		}
		callClient(port, methodName, id, *paramVal)

		return
	}

	app := api.NewApp()

	gateway := microapi.New("petByIdDispatch")
	service := gateway.NewService("PetStorePets", &rest.Activity{})
	service.SetDescription("Make calls to find pets")
	service.AddSetting("uri", "http://petstore.swagger.io/v2/pet/:id")
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
		panic(err)
	}

	gateway = microapi.New("petPutDispatch")
	service = gateway.NewService("PetStorePetsPUT", &rest.Activity{})
	service.SetDescription("Make calls to petstore")
	service.AddSetting("uri", "http://petstore.swagger.io/v2/pet")
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
		panic(err)
	}

	gateway = microapi.New("userByNameDispatch")
	service = gateway.NewService("PetStoreUsersByName", &rest.Activity{})
	service.SetDescription("Make calls to find users")
	service.AddSetting("uri", "http://petstore.swagger.io/v2/user/:username")
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
		panic(err)
	}

	triggerSettings := &trigger.Settings{
		Port:        9096,
		ProtoName:   "petstore",
		ServiceName: "PetStoreService",
	}
	trg := app.NewTrigger(&trigger.Trigger{}, triggerSettings)

	action := &microgateway.Action{}

	handler, err := trg.NewHandler(&trigger.HandlerSettings{
		AutoIdReply:     false,
		UseReplyHandler: false,
		MethodName:      "PetById",
	})
	if err != nil {
		panic(err)
	}
	_, err = handler.NewAction(action, petByIdDispatch)
	if err != nil {
		panic(err)
	}

	handler, err = trg.NewHandler(&trigger.HandlerSettings{
		AutoIdReply:     false,
		UseReplyHandler: false,
		MethodName:      "PetPUT",
	})
	if err != nil {
		panic(err)
	}
	_, err = handler.NewAction(action, petPutDispatch)
	if err != nil {
		panic(err)
	}

	handler, err = trg.NewHandler(&trigger.HandlerSettings{
		AutoIdReply:     false,
		UseReplyHandler: false,
		MethodName:      "UserByName",
	})
	if err != nil {
		panic(err)
	}
	_, err = handler.NewAction(action, userByNameDispatch)
	if err != nil {
		panic(err)
	}

	e, err := api.NewEngine(app)
	if err != nil {
		panic(err)
	}
	engine.RunEngine(e)
}

func callClient(port *string, option *string, id int, name string) {
	clientAddr = *port
	if len(*port) == 0 {
		clientAddr = ":9000"
	}
	clientAddr = ":" + *port
	conn, err := grpc.Dial("localhost"+clientAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := grpc2rest.NewPetStoreServiceClient(conn)

	switch *option {
	case "pet":
		petById(client, id)
	case "user":
		userByName(client, name)
	case "petput":
		petPUT(client, name)

	}
}

//UserByName comments
func userByName(client grpc2rest.PetStoreServiceClient, name string) {
	res, err := client.UserByName(context.Background(), &grpc2rest.UserByNameRequest{Username: name})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("res :", res)
}

//PetById comments
func petById(client grpc2rest.PetStoreServiceClient, id int) {
	res, err := client.PetById(context.Background(), &grpc2rest.PetByIdRequest{Id: int32(id)})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("res :", res)
}

func petPUT(client grpc2rest.PetStoreServiceClient, payload string) {

	name := strings.Split(payload, ",")[1]
	id, _ := strconv.ParseInt(strings.Split(payload, ",")[0], 0, 64)
	if id == 0 {
		log.Fatal("please provide valid id")
	}
	petReq := &grpc2rest.PetRequest{Id: int32(id), Name: name}
	fmt.Println("requesting pay load", petReq)
	res, err := client.PetPUT(context.Background(), petReq)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("res :", res)

}
