package grpc2rest

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/grpc"
)

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
	client := NewGRPC2RestPetStoreServiceClient(conn)

	switch *option {
	case "pet":
		id, err := strconv.Atoi(name)
		if err != nil {
			return nil, err
		}
		return petById(client, id)
	case "user":
		return userByName(client, name)
	case "petput":
		return petPUT(client, name)
	}

	return nil, nil
}

//UserByName comments
func userByName(client GRPC2RestPetStoreServiceClient, name string) (interface{}, error) {
	res, err := client.UserByName(context.Background(), &UserByNameRequest{Username: name})
	if err != nil {
		return nil, err
	}
	fmt.Println("res :", res)
	return res, nil
}

//PetById comments
func petById(client GRPC2RestPetStoreServiceClient, id int) (interface{}, error) {
	res, err := client.PetById(context.Background(), &PetByIdRequest{Id: int32(id)})
	if err != nil {
		return nil, err
	}
	fmt.Println("res :", res)
	return res, nil
}

func petPUT(client GRPC2RestPetStoreServiceClient, payload string) (interface{}, error) {

	name := strings.Split(payload, ",")[1]
	id, _ := strconv.ParseInt(strings.Split(payload, ",")[0], 0, 64)
	if id == 0 {
		return nil, errors.New("please provide valid id")
	}
	petReq := &PetRequest{Id: int32(id), Name: name}
	fmt.Println("requesting pay load", petReq)
	res, err := client.PetPUT(context.Background(), petReq)
	if err != nil {
		return nil, err
	}
	fmt.Println("res :", res)
	return res, nil
}
