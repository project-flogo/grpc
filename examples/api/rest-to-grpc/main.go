package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	trigger "github.com/project-flogo/contrib/trigger/rest"
	"github.com/project-flogo/core/api"
	"github.com/project-flogo/core/engine"
	"github.com/project-flogo/grpc/activity"
	"github.com/project-flogo/microgateway"
	microapi "github.com/project-flogo/microgateway/api"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var addr string
var clientAddr string

// ServerStrct is a stub for your Trigger implementation
type ServerStrct struct {
}

var petMapArr = make(map[int32]Pet)
var userMapArr = make(map[string]User)

var petArr = []Pet{
	{
		Id:   2,
		Name: "cat2",
	},
	{
		Id:   3,
		Name: "cat3",
	},
	{
		Id:   4,
		Name: "cat4",
	},
}
var userArr = []User{
	{
		Id:       2,
		Username: "user2",
		Email:    "email2",
		Phone:    "phone2",
	},
	{
		Id:       3,
		Username: "user3",
		Email:    "email3",
		Phone:    "phone3",
	},
	{
		Id:       4,
		Username: "user4",
		Email:    "email4",
		Phone:    "phone4",
	},
}

func main() {
	for _, pet := range petArr {
		petMapArr[pet.GetId()] = pet
	}

	for _, user := range userArr {
		userMapArr[user.GetUsername()] = user
	}

	clientMode := flag.Bool("client", false, "command to run client")
	serverMode := flag.Bool("server", false, "command to run server")
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
	} else if *serverMode {
		callServer()
		return
	}

	app := api.NewApp()

	gateway := microapi.New("PetStoreGETDispatch")
	service := gateway.NewService("PetStore", &activity.Activity{})
	service.SetDescription("Make calls to PetStore grpc end point")
	service.AddSetting("operatingMode", "rest-to-grpc")
	service.AddSetting("hosturl", "localhost:9000")
	step := gateway.NewStep(service)
	step.AddInput("protoName", "petstore")
	step.AddInput("serviceName", "PetStoreService")
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
		panic(err)
	}

	gateway = microapi.New("PetStorePUTDispatch")
	service = gateway.NewService("PetStore", &activity.Activity{})
	service.SetDescription("Make calls to PetStore grpc end point")
	service.AddSetting("operatingMode", "rest-to-grpc")
	service.AddSetting("hosturl", "localhost:9000")
	step = gateway.NewStep(service)
	step.AddInput("protoName", "petstore")
	step.AddInput("serviceName", "PetStoreService")
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
		panic(err)
	}

	trg := app.NewTrigger(&trigger.Trigger{}, &trigger.Settings{Port: 9096})

	handler, err := trg.NewHandler(&trigger.HandlerSettings{
		Method: "GET",
		Path:   "/petstore/method/:grpcMethodName",
	})
	if err != nil {
		panic(err)
	}
	_, err = handler.NewAction(&microgateway.Action{}, PetStoreGETDispatch)
	if err != nil {
		panic(err)
	}

	handler, err = trg.NewHandler(&trigger.HandlerSettings{
		Method: "PUT",
		Path:   "/petstore/PetPUT",
	})
	if err != nil {
		panic(err)
	}
	_, err = handler.NewAction(&microgateway.Action{}, PetStorePUTDispatch)
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
	client := NewPetStoreServiceClient(conn)

	switch *option {
	case "pet":
		petById(client, id)
	case "user":
		userByName(client, name)
	}
}

//UserByName comments
func userByName(client PetStoreServiceClient, name string) {
	res, err := client.UserByName(context.Background(), &UserByNameRequest{Username: name})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("res :", res)
}

//PetById comments
func petById(client PetStoreServiceClient, id int) {
	res, err := client.PetById(context.Background(), &PetByIdRequest{Id: int32(id)})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("res :", res)
}

func callServer() {
	addr = ":9000"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()

	RegisterPetStoreServiceServer(s, new(ServerStrct))

	log.Println("Starting server on port: ", addr)

	s.Serve(lis)
}

func (t *ServerStrct) PetById(ctx context.Context, req *PetByIdRequest) (*PetResponse, error) {

	fmt.Println("server PetById method called")
	for _, pet := range petMapArr {
		if pet.Id == req.Id {
			return &PetResponse{Pet: &pet}, nil
		}
	}

	return nil, errors.New("Pet not found")

}

func (t *ServerStrct) UserByName(ctx context.Context, req *UserByNameRequest) (*UserResponse, error) {

	fmt.Println("server UserByName method called")

	for _, user := range userMapArr {
		if req.Username == user.Username {
			return &UserResponse{User: &user}, nil
		}
	}

	return nil, errors.New("User not found")

}

func (t *ServerStrct) PetPUT(ctx context.Context, req *PetRequest) (*PetResponse, error) {

	fmt.Println("server PetPUT method called")
	if req.Pet == nil {
		return nil, errors.New("Content not found. Invalid Request")
	}
	if req.Pet.Id == 0 {
		return nil, errors.New("Invalid id provided")
	} else {
		fmt.Println("Request recieved for id:", req.Pet.Id)
	}

	for id := range petMapArr {
		if id == req.Pet.GetId() {
			petMapArr[id] = *req.GetPet()
		} else {
			petMapArr[req.GetPet().GetId()] = *req.GetPet()
		}
	}

	var pet = petMapArr[req.GetPet().GetId()]
	return &PetResponse{Pet: &pet}, nil
}

func (t *ServerStrct) UserPUT(ctx context.Context, req *UserRequest) (*UserResponse, error) {

	fmt.Println("server UserPUT method called")
	if req.User == nil {
		return nil, errors.New("Content not found. Invalid Request")
	}

	if len(req.User.Username) == 0 {
		return nil, errors.New("Invalid Username provided")
	} else {
		fmt.Println("Request recieved for username:", req.User.Username)
	}

	for name := range userMapArr {
		if strings.Compare(name, req.User.GetUsername()) == 0 {
			userMapArr[name] = *req.GetUser()
		} else {
			userMapArr[req.GetUser().GetUsername()] = *req.GetUser()
		}
	}

	var user = userMapArr[req.GetUser().GetUsername()]
	return &UserResponse{User: &user}, nil
}
