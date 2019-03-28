package examples

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/project-flogo/core/engine"
	"github.com/project-flogo/grpc/proto/grpc2grpc"
	"github.com/project-flogo/grpc/proto/grpc2rest"
	"github.com/project-flogo/grpc/proto/rest2grpc"
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

func testGRPC2Rest(t *testing.T, e engine.Engine) {
	defer api.ClearResources()

	util.Drain("9096")
	err := e.Start()
	assert.Nil(t, err)
	defer func() {
		err := e.Stop()
		assert.Nil(t, err)
	}()
	util.Pour("9096")

	server := &http.Server{Addr: ":8181"}
	r := mux.NewRouter()
	r.HandleFunc("/pet/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		petJSONPayload := "{\n \"Name\": \"testpet\" , \"Id\": 2 \n}"
		io.WriteString(w, petJSONPayload)
	})
	r.HandleFunc("/pet", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		petJSONPayload := "{\n \"name\": \"testpet\" , \"id\": 3 \n}"
		io.WriteString(w, petJSONPayload)
	})
	r.HandleFunc("/user/{username}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		userJSONPayload := "{\n \"username\": \"testuser\", \"id\": 2  \n}"
		io.WriteString(w, userJSONPayload)
	})
	done := make(chan bool, 1)
	go func() {
		server.Handler = r
		server.ListenAndServe()
		done <- true
	}()
	_, err = http.Get("http://localhost:8181/pet/1")
	for err != nil {
		_, err = http.Get("http://localhost:8181/pet/1")
	}
	defer func() {
		err := server.Shutdown(nil)
		if err != nil {
			t.Fatal(err)
		}
		<-done
	}()

	port, method := "9096", "petput"
	response, err := grpc2rest.CallClient(&port, &method, "3,testpet")
	assert.Nil(t, err)
	assert.Equal(t, int32(3), response.(*grpc2rest.PetResponse).Pet.Id)
	assert.Equal(t, "testpet", response.(*grpc2rest.PetResponse).Pet.Name)

	method = "pet"
	response, err = grpc2rest.CallClient(&port, &method, "2")
	assert.Nil(t, err)
	assert.Equal(t, int32(2), response.(*grpc2rest.PetResponse).Pet.Id)

	method = "user"
	response, err = grpc2rest.CallClient(&port, &method, "testuser")
	assert.Nil(t, err)
	assert.Equal(t, "testuser", response.(*grpc2rest.UserResponse).User.Username)
}

func TestGRPC2RestAPI(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping API integration test in short mode")
	}

	e, err := GRPC2RestExample()
	assert.Nil(t, err)
	testGRPC2Rest(t, e)
}

func TestGRPC2RestJSON(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping JSON integration test in short mode")
	}

	data, err := ioutil.ReadFile(filepath.FromSlash("./json/grpc-to-rest/flogo.json"))
	assert.Nil(t, err)
	cfg, err := engine.LoadAppConfig(string(data), false)
	assert.Nil(t, err)
	e, err := engine.New(cfg)
	assert.Nil(t, err)
	testGRPC2Rest(t, e)
}

func testRest2GRPC(t *testing.T, e engine.Engine) {
	defer api.ClearResources()
	util.Drain("9000")
	server, err := rest2grpc.CallServer()
	assert.Nil(t, err)
	go func() {
		server.Start()
	}()
	util.Pour("9000")
	defer server.Server.Stop()

	util.Drain("9096")
	err = e.Start()
	assert.Nil(t, err)
	defer func() {
		err := e.Stop()
		assert.Nil(t, err)
	}()
	util.Pour("9096")

	transport := &http.Transport{
		MaxIdleConns: 1,
	}
	defer transport.CloseIdleConnections()
	client := &http.Client{
		Transport: transport,
	}

	req, err := http.NewRequest(http.MethodGet, "http://localhost:9096/petstore/method/PetById?id=2", nil)
	assert.Nil(t, err)
	response, err := client.Do(req)
	assert.Nil(t, err)
	body, err := ioutil.ReadAll(response.Body)
	assert.Nil(t, err)
	response.Body.Close()
	var pet rest2grpc.PetResponse
	err = json.Unmarshal(body, &pet)
	assert.Nil(t, err)
	assert.Equal(t, int32(2), pet.Pet.Id)
	assert.Equal(t, "cat2", pet.Pet.Name)

	req, err = http.NewRequest(http.MethodGet, "http://localhost:9096/petstore/method/UserByName?username=user2", nil)
	assert.Nil(t, err)
	response, err = client.Do(req)
	assert.Nil(t, err)
	body, err = ioutil.ReadAll(response.Body)
	assert.Nil(t, err)
	response.Body.Close()
	var user rest2grpc.UserResponse
	err = json.Unmarshal(body, &user)
	assert.Nil(t, err)
	assert.Equal(t, int32(2), user.User.Id)
	assert.Equal(t, "user2", user.User.Username)

	payload := `{"pet": {"id": 12,"name": "mycat12"}}`
	req, err = http.NewRequest(http.MethodPut, "http://localhost:9096/petstore/PetPUT", bytes.NewReader([]byte(payload)))
	assert.Nil(t, err)
	req.Header.Add("Content-Type", "application/json")
	response, err = client.Do(req)
	assert.Nil(t, err)
	body, err = ioutil.ReadAll(response.Body)
	assert.Nil(t, err)
	response.Body.Close()
	err = json.Unmarshal(body, &pet)
	assert.Nil(t, err)
	assert.Equal(t, int32(12), pet.Pet.Id)
	assert.Equal(t, "mycat12", pet.Pet.Name)
}

func TestRest2GRPCAPI(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping API integration test in short mode")
	}

	e, err := Rest2GRPCExample()
	assert.Nil(t, err)
	testRest2GRPC(t, e)
}

func TestRest2GRPCJSON(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping JSON integration test in short mode")
	}

	data, err := ioutil.ReadFile(filepath.FromSlash("./json/rest-to-grpc/flogo.json"))
	assert.Nil(t, err)
	cfg, err := engine.LoadAppConfig(string(data), false)
	assert.Nil(t, err)
	e, err := engine.New(cfg)
	assert.Nil(t, err)
	testRest2GRPC(t, e)
}
