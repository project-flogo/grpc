package rest2grpc

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// ServerStrct is a stub for your Trigger implementation
type ServerStrct struct {
	Server     *grpc.Server
	petMapArr  map[int32]Pet
	userMapArr map[string]User
}

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

func CallClient(port *string, option *string, name string) (interface{}, error) {
	clientAddr := *port
	if len(*port) == 0 {
		clientAddr = ":9000"
	}
	clientAddr = ":" + *port
	conn, err := grpc.Dial("localhost"+clientAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := NewRest2GRPCPetStoreServiceClient(conn)

	switch *option {
	case "pet":
		id, err := strconv.Atoi(name)
		if err != nil {
			return nil, err
		}
		return petById(client, id)
	case "user":
		return userByName(client, name)
	}

	return nil, nil
}

//UserByName comments
func userByName(client Rest2GRPCPetStoreServiceClient, name string) (interface{}, error) {
	res, err := client.UserByName(context.Background(), &UserByNameRequest{Username: name})
	if err != nil {
		return nil, err
	}
	fmt.Println("res :", res)
	return res, nil
}

//PetById comments
func petById(client Rest2GRPCPetStoreServiceClient, id int) (interface{}, error) {
	res, err := client.PetById(context.Background(), &PetByIdRequest{Id: int32(id)})
	if err != nil {
		return nil, err
	}
	fmt.Println("res :", res)
	return res, nil
}

func CallServer() (*ServerStrct, error) {
	s := grpc.NewServer()
	server := &ServerStrct{
		Server:     s,
		petMapArr:  make(map[int32]Pet),
		userMapArr: make(map[string]User),
	}
	for _, pet := range petArr {
		server.petMapArr[pet.GetId()] = pet
	}

	for _, user := range userArr {
		server.userMapArr[user.GetUsername()] = user
	}

	RegisterRest2GRPCPetStoreServiceServer(s, server)

	return server, nil
}

func (t *ServerStrct) Start() error {
	addr := ":9000"
	log.Println("Starting server on port: ", addr)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return t.Server.Serve(lis)
}

func (t *ServerStrct) PetById(ctx context.Context, req *PetByIdRequest) (*PetResponse, error) {

	fmt.Println("server PetById method called")
	for _, pet := range t.petMapArr {
		if pet.Id == req.Id {
			return &PetResponse{Pet: &pet}, nil
		}
	}

	return nil, errors.New("Pet not found")

}

func (t *ServerStrct) UserByName(ctx context.Context, req *UserByNameRequest) (*UserResponse, error) {

	fmt.Println("server UserByName method called")

	for _, user := range t.userMapArr {
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

	for id := range t.petMapArr {
		if id == req.Pet.GetId() {
			t.petMapArr[id] = *req.GetPet()
		} else {
			t.petMapArr[req.GetPet().GetId()] = *req.GetPet()
		}
	}

	var pet = t.petMapArr[req.GetPet().GetId()]
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

	for name := range t.userMapArr {
		if strings.Compare(name, req.User.GetUsername()) == 0 {
			t.userMapArr[name] = *req.GetUser()
		} else {
			t.userMapArr[req.GetUser().GetUsername()] = *req.GetUser()
		}
	}

	var user = t.userMapArr[req.GetUser().GetUsername()]
	return &UserResponse{User: &user}, nil
}
