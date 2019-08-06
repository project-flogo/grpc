package grpc

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/golang/protobuf/jsonpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/golang/protobuf/proto"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/trigger"
)

var (
	addr            string
	triggerMetadata = trigger.NewMetadata(&Settings{}, &HandlerSettings{}, &Output{}, &Reply{})
)

func init() {
	trigger.Register(&Trigger{}, &Factory{})
}

// Factory is a gRPC Trigger factory
type Factory struct {
	metadata *trigger.Metadata
}

// New creates a new trigger instance
func (*Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	s := &Settings{}
	err := metadata.MapToStruct(config.Settings, s, true)
	if err != nil {
		return nil, err
	}
	return &Trigger{settings: s, config: config}, nil
}

// Metadata get trigger metadata
func (f *Factory) Metadata() *trigger.Metadata {
	return triggerMetadata
}

// Handler struct
type Handler struct {
	handler  trigger.Handler
	settings *HandlerSettings
}

// Trigger is a stub for your gRPC Trigger implementation
type Trigger struct {
	config         *trigger.Config
	settings       *Settings
	handlers       map[string]*Handler
	defaultHandler *Handler
	server         *grpc.Server
	Logger         log.Logger
}

// Metadata implements trigger.Trigger.Metadata
func (t *Trigger) Metadata() *trigger.Metadata {
	return triggerMetadata
}

// Initialize implements trigger.Trigger.Initialize
func (t *Trigger) Initialize(ctx trigger.InitContext) error {
	logger := ctx.Logger()
	t.Logger = logger

	addr = ":" + strconv.Itoa(t.settings.Port)

	handlers := make(map[string]*Handler)
	for _, handler := range ctx.GetHandlers() {
		settings := &HandlerSettings{}
		err := metadata.MapToStruct(handler.Settings(), settings, true)
		if err != nil {
			return err
		}
		if settings.MethodName == "" && t.defaultHandler == nil {
			t.defaultHandler = &Handler{
				handler:  handler,
				settings: settings,
			}
		} else {
			handlers[settings.ServiceName+"_"+settings.MethodName] = &Handler{
				handler:  handler,
				settings: settings,
			}
		}
	}
	t.handlers = handlers
	t.Logger.Debugf("Enable TLS: %t", t.settings.EnableTLS)
	if t.settings.EnableTLS {
		// decode server cert and server key
		serverCert, err := t.decodeCertificate(t.settings.ServerCert)
		if err != nil {
			t.Logger.Errorf("Error decoding server certificate: %s", err.Error())
			return err
		}
		serverKey, err := t.decodeCertificate(t.settings.ServerKey)
		if err != nil {
			t.Logger.Errorf("Error decoding server key: %s", err.Error())
			return err
		}
		t.settings.ServerCert = string(serverCert)
		t.settings.ServerKey = string(serverKey)
	}
	return nil
}

// Stop implements trigger.Trigger.Start
func (t *Trigger) Stop() error {
	// stop the trigger
	t.server.GracefulStop()
	return nil
}

// Start implements trigger.Trigger.Start
func (t *Trigger) Start() error {
	// start the trigger
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		t.Logger.Error(err)
		return err
	}

	opts := []grpc.ServerOption{}

	if t.settings.EnableTLS {
		cert, err := tls.X509KeyPair([]byte(t.settings.ServerCert), []byte(t.settings.ServerKey))
		if err != nil {
			t.Logger.Error(err)
			return err
		}
		creds := credentials.NewServerTLSFromCert(&cert)
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	t.server = grpc.NewServer(opts...)

	protoName := t.settings.ProtoName
	protoName = strings.Split(protoName, ".")[0]

	// Register each serviceName + protoName
	if len(ServiceRegistery.ServerServices) != 0 {
		for k, service := range ServiceRegistery.ServerServices {
			servRegFlag := false
			if strings.Compare(k, protoName+service.ServiceInfo().ServiceName) == 0 {
				t.Logger.Infof("Registered Proto [%v] and Service [%v]", protoName, service.ServiceInfo().ServiceName)
				service.RunRegisterServerService(t.server, t)
				servRegFlag = true
			}
			if !servRegFlag {
				t.Logger.Errorf("Proto [%s] and Service [%s] not registered", protoName, service.ServiceInfo().ServiceName)
				return fmt.Errorf("Proto [%s] and Service [%s] not registered", protoName, service.ServiceInfo().ServiceName)
			}
		}

	} else {
		t.Logger.Error("gRPC server services not registered")
		return errors.New("gRPC server services not registered")
	}

	t.Logger.Debug("Starting server on port", addr)

	go func() {
		t.server.Serve(lis)
	}()

	t.Logger.Infof("Server started on port: [%d]", t.settings.Port)
	return nil
}

// CallHandler is to call a particular handler based on method name
func (t *Trigger) CallHandler(grpcData map[string]interface{}) (int, interface{}, error) {
	t.Logger.Debug("CallHandler method invoked")

	params := make(map[string]interface{})
	var content interface{}
	m := jsonpb.Marshaler{OrigName: true}
	// blocking the code for streaming requests
	if grpcData["contextdata"] != nil {
		// getting values from inputrequestdata and mapping it to params which can be used in different services like HTTP pathparams etc.
		s := reflect.ValueOf(grpcData["reqdata"]).Elem()
		typeOfS := s.Type()
		for i := 0; i < s.NumField(); i++ {
			f := s.Field(i)
			fieldName := proto.GetProperties(typeOfS).Prop[i].OrigName
			if !strings.HasPrefix(fieldName, "XXX_") {
				// XXX_ fields will not be mapped
				if _, ok := f.Interface().(proto.Message); ok {
					jsonString, err := m.MarshalToString(f.Interface().(proto.Message))
					if err != nil {
						t.Logger.Errorf("Marshal failed on field: %s with value: %v", fieldName, f.Interface())
					}
					t.Logger.Debugf("Marshaled FieldName: [%s] Value: [%s]", fieldName, jsonString)
					var paramValue map[string]interface{}
					json.Unmarshal([]byte(jsonString), &paramValue)
					params[fieldName] = paramValue
				} else {
					t.Logger.Debugf("Field name: [%s] Value: [%v]", fieldName, f.Interface())
					params[fieldName] = f.Interface()
				}
			}
		}

		// assign req data content to trigger content
		dataBytes, err := json.Marshal(grpcData["reqdata"])
		if err != nil {
			t.Logger.Error("Marshal failed on grpc request data")
			return 0, nil, err
		}

		err = json.Unmarshal(dataBytes, &content)
		if err != nil {
			t.Logger.Error("Unmarshal failed on grpc request data")
			return 0, nil, err
		}
	}

	handler, ok := t.handlers[grpcData["serviceName"].(string)+"_"+grpcData["methodName"].(string)]
	if !ok {
		handler = t.defaultHandler
	}

	if handler != nil {
		grpcData["protoName"] = t.settings.ProtoName

		out := &Output{
			Params:   params,
			GrpcData: grpcData,
			Content:  content,
		}

		t.Logger.Debug("Dispatch Found for ", handler.settings.ServiceName+"_"+handler.settings.MethodName)
		results, err := handler.handler.Handle(context.Background(), out)
		if err != nil {
			return 0, nil, err
		}
		reply := &Reply{}
		err = reply.FromMap(results)
		t.Logger.Debugf("Result from handler: %v", reply.Data)
		return reply.Code, reply.Data, err
	}

	t.Logger.Error("Dispatch not found")
	return 0, nil, errors.New("Dispatch not found")
}

func (t *Trigger) decodeCertificate(cert string) ([]byte, error) {
	if cert == "" {
		return nil, fmt.Errorf("Certificate is Empty")
	}

	// case 1: if certificate comes from fileselctor it will be base64 encoded
	if strings.HasPrefix(cert, "{") {
		t.Logger.Debug("Certificate received from file selector")
		certObj, err := coerce.ToObject(cert)
		if err == nil {
			certValue, ok := certObj["content"].(string)
			if !ok || certValue == "" {
				return nil, fmt.Errorf("No content found for certificate")
			}
			return base64.StdEncoding.DecodeString(strings.Split(certValue, ",")[1])
		}
		return nil, err
	}

	// case 2: if the certificate is defined as application property in the format "<encoding>,<encodedCertificateValue>"
	index := strings.IndexAny(cert, ",")
	if index > -1 {
		//some encoding is there
		t.Logger.Debug("Certificate received from application property with encoding")
		encoding := cert[:index]
		certValue := cert[index+1:]

		if strings.EqualFold(encoding, "base64") {
			return base64.StdEncoding.DecodeString(certValue)
		}
		return nil, fmt.Errorf("Error parsing the certificate or given encoding may not be supported")
	}

	// case 3: if the certificate is defined as application property that points to a file
	if strings.HasPrefix(cert, "file://") {
		// app property pointing to a file
		t.Logger.Debug("Certificate received from application property pointing to a file")
		fileName := t.settings.ServerCert[7:]
		return ioutil.ReadFile(fileName)
	}

	// case 4: if certificate is defined as path to a file (in oss)
	if strings.Contains(cert, "/") || strings.Contains(cert, "\\") {
		t.Logger.Debug("Certificate received from settings as file path")
		_, err := os.Stat(cert)
		if err != nil {
			t.Logger.Errorf("Cannot find certificate file: %s", err.Error())
		}
		return ioutil.ReadFile(cert)
	}

	t.Logger.Debug("Certificate received from application property without encoding")
	return []byte(cert), nil
}
