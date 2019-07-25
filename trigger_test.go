package main

import (
	"context"
	"testing"

	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/trigger"
	"github.com/project-flogo/grpc/proto/grpc2grpc"
	_ "github.com/project-flogo/grpc/trigger"
	"github.com/project-flogo/grpc/util"
	"github.com/stretchr/testify/assert"
)

type handler struct {
	handled bool
}

func (h *handler) Name() string {
	return "test"
}

func (h *handler) Settings() map[string]interface{} {
	return map[string]interface{}{
		"autoIdReply":     false,
		"useReplyHandler": false,
	}
}

func (h *handler) Handle(ctx context.Context, triggerData interface{}) (map[string]interface{}, error) {
	h.handled = true
	return map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{"message": "hello world"},
	}, nil
}

type triggerInitContext struct {
	handlers []trigger.Handler
}

func (i *triggerInitContext) Logger() log.Logger {
	return log.RootLogger()
}

func (i *triggerInitContext) GetHandlers() []trigger.Handler {
	return i.handlers
}

func TestGRPCTrigger(t *testing.T) {
	factory := trigger.GetFactory("github.com/project-flogo/grpc/trigger")
	assert.NotNil(t, factory)
	hndlrConfigs := []*trigger.HandlerConfig{}
	hc := &trigger.HandlerConfig{
		Settings: map[string]interface{}{
			"serviceName": "PetStoreService",
		},
	}
	hndlrConfigs = append(hndlrConfigs, hc)
	config := trigger.Config{
		Id: "test",
		Settings: map[string]interface{}{
			"port":      9096,
			"protoName": "petstore",
		},
		Handlers: hndlrConfigs,
	}
	instance, err := factory.New(&config)
	assert.Nil(t, err)

	h := handler{}
	context := triggerInitContext{
		handlers: []trigger.Handler{
			&h,
		},
	}
	err = instance.Initialize(&context)
	assert.Nil(t, err)

	util.Drain("9096")
	instance.Start()
	util.Pour("9096")
	defer instance.Stop()

	port, method := "9096", "pet"
	_, err = grpc2grpc.CallClient(&port, &method, "2", nil)
	assert.Nil(t, err)
	assert.True(t, h.handled)
}
