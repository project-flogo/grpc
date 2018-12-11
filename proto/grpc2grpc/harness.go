package grpc2grpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

// ServerStrct is a stub for your Trigger implementation
type ServerStrct struct {
	Server     *grpc.Server
	petMapArr  map[int32]Pet
	userMapArr map[string]User
	Done       chan bool
}

var recUserArrClient = make([]User, 0)
var recUserArrServer = make([]User, 0)

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

var sendUserDataClient = []User{
	{
		Id:       22,
		Username: "user22c",
		Email:    "email22c",
		Phone:    "phone22c",
	},
	{
		Id:       23,
		Username: "user23c",
		Email:    "email23c",
		Phone:    "phone23c",
	},
	{
		Id:       24,
		Username: "user24c",
		Email:    "email24c",
		Phone:    "phone24c",
	},
}

var sendUserDataServer = []User{
	{
		Id:       32,
		Username: "user32s",
		Email:    "email32s",
		Phone:    "phone32s",
	},
	{
		Id:       33,
		Username: "user33s",
		Email:    "email33s",
		Phone:    "phone33s",
	},
	{
		Id:       34,
		Username: "user34s",
		Email:    "email34s",
		Phone:    "phone34s",
	},
}

func CallClient(port *string, option *string, name string, data func(data interface{}) bool) (interface{}, error) {
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
	client := NewPetStoreServiceClient(conn)

	switch *option {
	case "pet":
		id, _ := strconv.Atoi(name)
		return petById(client, id)
	case "user":
		return userByName(client, name)
	case "listusers":
		err := listUsers(client, data)
		return nil, err
	case "storeusers":
		err := storeUsers(client, data)
		return nil, err
	case "bulkusers":
		err := bulkUsers(client, data)
		return nil, err
	}

	return nil, nil
}

//userByName comments
func userByName(client PetStoreServiceClient, name string) (interface{}, error) {
	res, err := client.UserByName(context.Background(), &UserByNameRequest{Username: name})
	if err != nil {
		return nil, err
	}
	fmt.Println("res :", res)

	return res, nil
}

//petById comments
func petById(client PetStoreServiceClient, id int) (interface{}, error) {
	res, err := client.PetById(context.Background(), &PetByIdRequest{Id: int32(id)})
	if err != nil {
		return nil, err
	}
	fmt.Println("res :", res)

	return res, nil
}

func listUsers(client PetStoreServiceClient, data func(data interface{}) bool) error {
	stream, err := client.ListUsers(context.Background(), &EmptyReq{Msg: "list of users"})
	if err != nil {
		return err
	}
	for {
		users, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if data != nil && data(users) {
			break
		}
		fmt.Println("RECEIVED: ", users.GetUsername())
	}

	return nil
}

func storeUsers(client PetStoreServiceClient, data func(data interface{}) bool) error {
	stream, err := client.StoreUsers(context.Background())
	if err != nil {
		return err
	}
	var i = 0
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	waitc := make(chan struct{})
	go func() {
		for {
			select {
			case <-c:
				close(waitc)
				i = 0
				return
			default:
				i++
				user := User{
					Username: "cuser" + strconv.Itoa(i),
					Id:       int32(i),
					Email:    "cemail" + strconv.Itoa(i),
				}
				fmt.Println("SEND: ", user.Username)
				if err := stream.Send(&user); err != nil {
					fmt.Println("error while sending user", user, err)
					close(waitc)
					return
				}
				if data != nil && data(&user) {
					close(waitc)
					return
				}
			}
		}
	}()
	<-waitc
	reply, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Println("erorr occured in StoreUsers response", err)
	} else {
		fmt.Println("response received from server:", reply)
	}

	return err
}

func bulkUsers(client PetStoreServiceClient, data func(data interface{}) bool) error {
	stream, err := client.BulkUsers(context.Background())
	if err != nil {
		return err
	}
	var rc = 0
	waitc := make(chan struct{})
	go func() {
		for {
			time.Sleep(time.Second)
			user, err := stream.Recv()
			if err == io.EOF {
				// read done.
				fmt.Println("received total: ", strconv.Itoa(rc))
				recUserArrClient = nil
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a user : %v", err)
				close(waitc)
				return
			}
			if user != nil {
				rc++
				fmt.Println("RECEIVED:  ", user.GetUsername())
				if data != nil {
					data(user)
				}
			}
		}
	}()
	for i := 0; i < 10; i++ {
		user := User{
			Username: "cuser" + strconv.Itoa(i),
			Id:       int32(i),
			Email:    "cemail" + strconv.Itoa(i),
		}
		fmt.Println("SEND: ", user.Username)
		if err := stream.Send(&user); err != nil {
			fmt.Println("error while sending user", user, err)
			return err
		}
		time.Sleep(time.Second)
	}
	stream.CloseSend()
	<-waitc

	return nil
}

func CallServer() (*ServerStrct, error) {
	s := grpc.NewServer()
	server := ServerStrct{
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

	RegisterPetStoreServiceServer(s, &server)

	return &server, nil
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

func (t *ServerStrct) ListUsers(req *EmptyReq, sReq PetStoreService_ListUsersServer) error {
	fmt.Println("server ListUsers method called")
	fmt.Println("request received from client ", req)
	var i = 0

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)

	for {
		select {
		case <-t.Done:
			return nil
		case <-exit:
			fmt.Println("total users sent:", strconv.Itoa(i))
			close(exit)
			return nil
		default:
			i++
			user := User{
				Username: "suser" + strconv.Itoa(i),
				Id:       int32(i),
				Email:    "semail" + strconv.Itoa(i),
			}
			fmt.Println("SEND: ", user.Username)
			err := sReq.Send(&user)
			if err != nil {
				fmt.Println("unable to send:", err)
				return err
			}
		}
	}
}
func (t *ServerStrct) StoreUsers(cReq PetStoreService_StoreUsersServer) error {
	fmt.Println("server StoreUsers method called")
	var i = 0
	for {
		userData, err := cReq.Recv()
		if userData != nil {
			i++
			fmt.Println("RECEIVED: ", userData.GetUsername())
			if t.Done != nil {
				t.Done <- true
			}
			fmt.Println("over there")
		}
		if err == io.EOF {
			return cReq.SendAndClose(&EmptyRes{Msg: "received total users:" + strconv.Itoa(i)})
		}
		if err != nil {
			return err
		}
	}
}
func (t *ServerStrct) BulkUsers(bReq PetStoreService_BulkUsersServer) error {
	fmt.Println("server BulkUsers method called")
	var rc = 0
	waits := make(chan struct{})
	go func() {
		for {
			time.Sleep(time.Second)
			userData, err := bReq.Recv()
			if userData != nil {
				rc++
				fmt.Println("RECEIVED: ", userData.GetUsername())
				if t.Done != nil {
					t.Done <- true
				}
			}
			if err == io.EOF {
				fmt.Println("total records received: ", strconv.Itoa(rc))
				close(waits)
				return
			}
			if err != nil {
				fmt.Println("error occured while receiving", err)
				close(waits)
				return
			}
		}
	}()
	for i := 0; i < 10; i++ {
		user := User{
			Username: "suser" + strconv.Itoa(i),
			Id:       int32(i),
			Email:    "semail" + strconv.Itoa(i),
		}
		fmt.Println("SEND: ", user.Username)
		if err := bReq.Send(&user); err != nil {
			fmt.Println("error while sending user", user, err)
			return err
		}
		time.Sleep(time.Second)
	}
	<-waits
	if t.Done != nil {
		close(t.Done)
	}
	return nil
}
