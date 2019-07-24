package trigger

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"

	"github.com/golang/protobuf/jsonpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/golang/protobuf/proto"
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

// GRPCTriggerFactory is a gRPC Trigger factory
type Factory struct {
	metadata *trigger.Metadata
}

func (*Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	s := &Settings{}
	err := metadata.MapToStruct(config.Settings, s, true)
	if err != nil {
		return nil, err
	}
	return &Trigger{settings: s, config: config}, nil
}

func (f *Factory) Metadata() *trigger.Metadata {
	return triggerMetadata
}

type Handler struct {
	handler  trigger.Handler
	settings *HandlerSettings
}

// GRPCTrigger is a stub for your Trigger implementation
type Trigger struct {
	config         *trigger.Config
	settings       *Settings
	handlers       map[string]*Handler
	defaultHandler *Handler
	server         *grpc.Server
	TLSConfig
	Logger log.Logger
}

// TLSConfig is to hold tls support data
type TLSConfig struct {
	enableTLS bool
	serveKey  string
	serveCert string
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

	//Check whether TLS (Transport Layer Security) is enabled for the trigger
	enableTLS := false
	serverCert := ""
	serverKey := ""
	if t.settings.EnableTLS {
		//TLS is enabled, get server certificate & key
		enableTLS = true
		serverCert = t.settings.ServerCert
		if serverCert == "" {
			panic(fmt.Sprintf("No serverCert found for trigger '%s' in settings", t.config.Id))
		}

		serverKey = t.settings.ServerKey
		if serverKey == "" {
			panic(fmt.Sprintf("No serverKey found for trigger '%s' in settings", t.config.Id))
		}

	}

	logger.Debug("enableTLS: ", enableTLS)
	if enableTLS {
		logger.Debug("serverCert: ", serverCert)
		logger.Debug("serverKey: ", serverKey)
	}
	t.enableTLS = enableTLS
	t.serveCert = serverCert
	t.serveKey = serverKey

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

	if t.enableTLS {
		creds, err := credentials.NewServerTLSFromFile(t.serveCert, t.serveKey)
		if err != nil {
			t.Logger.Error(err)
			return err
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	t.server = grpc.NewServer(opts...)

	// collect all service names used in handlers
	serviceNames := make(map[string]string)
	for _, h := range t.handlers {
		if _, ok := serviceNames[h.settings.ServiceName]; !ok {
			serviceNames[h.settings.ServiceName] = h.settings.ServiceName
		}
	}
	protoName := t.settings.ProtoName
	protoName = strings.Split(protoName, ".")[0]

	// Register each serviceName + protoName
	servRegFlag := false
	for serviceName := range serviceNames {
		if len(ServiceRegistery.ServerServices) != 0 {
			for k, service := range ServiceRegistery.ServerServices {
				if strings.Compare(k, protoName+serviceName) == 0 {
					t.Logger.Infof("Registered Proto [%v] and Service [%v]", protoName, serviceName)
					service.RunRegisterServerService(t.server, t)
					servRegFlag = true
				}
			}
			if !servRegFlag {
				t.Logger.Errorf("Proto [%s] and Service [%s] not registered", protoName, serviceName)
				return fmt.Errorf("Proto [%s] and Service [%s] not registered", protoName, serviceName)
			}
		} else {
			t.Logger.Error("gRPC server services not registered")
			return errors.New("gRPC server services not registered")
		}
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
