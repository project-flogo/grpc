package examples

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/project-flogo/core/engine"
	"github.com/project-flogo/grpc/proto/grpc2grpc"
	"github.com/project-flogo/grpc/util"
	"github.com/project-flogo/microgateway/api"
)

func testGRPC2GRPC(t *testing.T, e engine.Engine) {
	defer api.ClearResources()
	util.Drain("9000")
	server, err := grpc2grpc.CallServer()
	assert.Nil(t, err)
	go func() {
		server.Start()
	}()
	util.Pour("9000")
	defer server.Server.Stop()

	server.Done = make(chan bool, 16)

	util.Drain("9096")
	err = e.Start()
	assert.Nil(t, err)
	defer func() {
		err := e.Stop()
		assert.Nil(t, err)
	}()
	util.Pour("9096")

	port, method := "9096", "pet"
	response, err := grpc2grpc.CallClient(&port, &method, "2", nil)
	assert.Nil(t, err)
	assert.Equal(t, "cat2", response.(*grpc2grpc.PetResponse).Pet.Name)

	method = "user"
	response, err = grpc2grpc.CallClient(&port, &method, "user2", nil)
	assert.Nil(t, err)
	assert.Equal(t, "user2", response.(*grpc2grpc.UserResponse).User.Username)

	method, count := "listusers", 0
	_, err = grpc2grpc.CallClient(&port, &method, "", func(data interface{}) bool {
		count++
		if count == 100 {
			server.Done <- true
		}
		return false
	})
	assert.Nil(t, err)
	assert.Condition(t, func() (success bool) { return count >= 100 }, "count should be >= 100")

	method, count = "storeusers", 0
	_, err = grpc2grpc.CallClient(&port, &method, "", func(data interface{}) bool {
		<-server.Done
		count++
		if count == 100 {
			return true
		}
		return false
	})
	assert.Nil(t, err)
	assert.Condition(t, func() (success bool) { return count >= 100 }, "count should be >= 100")

	method, count = "bulkusers", 0
	_, err = grpc2grpc.CallClient(&port, &method, "", func(data interface{}) bool {
		count++
		return false
	})
	for range server.Done {
		count++
	}
	server.Done = nil
	assert.Nil(t, err)
	assert.Equal(t, 20, count)
}

func TestGRPC2GRPCAPI(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping API integration test in short mode")
	}

	e, err := GRPC2GRPCExample()
	assert.Nil(t, err)
	testGRPC2GRPC(t, e)
}

func TestGRPC2GRPCJSON(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping JSON integration test in short mode")
	}

	data, err := ioutil.ReadFile(filepath.FromSlash("./json/grpc-to-grpc/flogo.json"))
	assert.Nil(t, err)
	cfg, err := engine.LoadAppConfig(string(data), false)
	assert.Nil(t, err)
	e, err := engine.New(cfg)
	assert.Nil(t, err)
	testGRPC2GRPC(t, e)
}
