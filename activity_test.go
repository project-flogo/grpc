package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/mapper"
	"github.com/project-flogo/core/data/metadata"
	logger "github.com/project-flogo/core/support/log"
	grpcactivity "github.com/project-flogo/grpc/activity/grpc"
	"github.com/project-flogo/grpc/proto/rest2grpc"
)

// ServerStrct is a stub for your Trigger implementation
type ServerStrct struct {
}

var (
	petMapArr  = make(map[int32]rest2grpc.Pet)
	userMapArr = make(map[string]rest2grpc.User)
)

// PetById gets a pet by id
func (t *ServerStrct) PetById(ctx context.Context, req *rest2grpc.PetByIdRequest) (*rest2grpc.PetResponse, error) {
	for _, pet := range petMapArr {
		if pet.Id == req.Id {
			return &rest2grpc.PetResponse{Pet: &pet}, nil
		}
	}
	return nil, errors.New("Pet not found")
}

// UserByName gets a user by name
func (t *ServerStrct) UserByName(ctx context.Context, req *rest2grpc.UserByNameRequest) (*rest2grpc.UserResponse, error) {
	for _, user := range userMapArr {
		if req.Username == user.Username {
			return &rest2grpc.UserResponse{User: &user}, nil
		}
	}
	return nil, errors.New("User not found")
}

func (t *ServerStrct) PetPUT(ctx context.Context, req *rest2grpc.PetRequest) (*rest2grpc.PetResponse, error) {

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
	return &rest2grpc.PetResponse{Pet: &pet}, nil
}

func (t *ServerStrct) UserPUT(ctx context.Context, req *rest2grpc.UserRequest) (*rest2grpc.UserResponse, error) {

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
	return &rest2grpc.UserResponse{User: &user}, nil
}

type initContext struct {
	settings map[string]interface{}
}

func newInitContext(values map[string]interface{}) *initContext {
	if values == nil {
		values = make(map[string]interface{})
	}
	return &initContext{
		settings: values,
	}
}

func (i *initContext) Settings() map[string]interface{} {
	return i.settings
}

func (i *initContext) MapperFactory() mapper.Factory {
	return nil
}

func (i *initContext) Logger() logger.Logger {
	return logger.RootLogger()
}

type activityContext struct {
	input  map[string]interface{}
	output map[string]interface{}
}

func newActivityContext(values map[string]interface{}) *activityContext {
	if values == nil {
		values = make(map[string]interface{})
	}
	return &activityContext{
		input:  values,
		output: make(map[string]interface{}),
	}
}

func (a *activityContext) ActivityHost() activity.Host {
	return a
}

func (a *activityContext) Name() string {
	return "test"
}

func (a *activityContext) GetInput(name string) interface{} {
	return a.input[name]
}

func (a *activityContext) SetOutput(name string, value interface{}) error {
	a.output[name] = value
	return nil
}

func (a *activityContext) GetInputObject(input data.StructValue) error {
	return input.FromMap(a.input)
}

func (a *activityContext) SetOutputObject(output data.StructValue) error {
	a.output = output.ToMap()
	return nil
}

func (a *activityContext) GetSharedTempData() map[string]interface{} {
	return nil
}

func (a *activityContext) ID() string {
	return "test"
}

func (a *activityContext) IOMetadata() *metadata.IOMetadata {
	return nil
}

func (a *activityContext) Reply(replyData map[string]interface{}, err error) {

}

func (a *activityContext) Return(returnData map[string]interface{}, err error) {

}

func (a *activityContext) Scope() data.Scope {
	return nil
}

func (a *activityContext) Logger() logger.Logger {
	return logger.RootLogger()
}

func TestGRPC(t *testing.T) {
	var petArr = []rest2grpc.Pet{
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
	var userArr = []rest2grpc.User{
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
	for _, pet := range petArr {
		petMapArr[pet.GetId()] = pet
	}
	for _, user := range userArr {
		userMapArr[user.GetUsername()] = user
	}

	addr := ":9000"
	socket, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	server := grpc.NewServer()
	rest2grpc.RegisterRest2GRPCPetStoreServiceServer(server, &ServerStrct{})

	done := make(chan bool, 1)
	go func() {
		server.Serve(socket)
		done <- true
	}()
	defer func() {
		server.GracefulStop()
		<-done
	}()

	activity, err := grpcactivity.New(newInitContext(map[string]interface{}{
		"operatingMode": "grpc-to-grpc",
		"hosturl":       "localhost:9000",
	}))
	assert.Nil(t, err)

	grpcData := map[string]interface{}{
		"methodName":  "PetById",
		"contextdata": context.Background(),
		"reqdata":     &rest2grpc.PetByIdRequest{Id: 2},
		"serviceName": "Rest2GRPCPetStoreService",
		"protoName":   "petstore",
	}
	ctx := newActivityContext(map[string]interface{}{
		"grpcMthdParamtrs": grpcData,
	})
	_, err = activity.Eval(ctx)
	assert.Nil(t, err)

	body := ctx.output["body"]
	if pet, ok := body.(*rest2grpc.PetResponse); !ok {
		t.Fatal("should be pet response")
	} else if pet.Pet.Name != "cat2" {
		t.Fatal("didn't get correct pet")
	}

	activity, err = grpcactivity.New(newInitContext(map[string]interface{}{
		"operatingMode": "rest-to-grpc",
		"hosturl":       "localhost:9000",
	}))
	assert.Nil(t, err)

	ctx = newActivityContext(map[string]interface{}{
		"protoName":   "petstore",
		"serviceName": "Rest2GRPCPetStoreService",
		"methodName":  "PetById",
		"queryParams": map[string]string{
			"id": "2",
		},
	})
	_, err = activity.Eval(ctx)
	assert.Nil(t, err)

	if a, ok := ctx.output["body"].(map[string]interface{}); !ok {
		t.Fatal("should be a map")
	} else if b := a["pet"]; b == nil {
		t.Fatal("pet should not be nil")
	} else if c, ok := b.(map[string]interface{}); !ok {
		t.Fatal("pet should be a map")
	} else if d := c["name"]; d == nil {
		t.Fatal("name should not be nil")
	} else if e, ok := d.(string); !ok {
		t.Fatal("name should be a string")
	} else if e != "cat2" {
		t.Fatal("name should be equal to cat2")
	}
}
