package grpc

import (
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct {
	Port       int    `md:"port,required"`
	ProtoName  string `md:"protoName,required"`
	ProtoFile  string `md:"protoFile"`
	EnableTLS  bool   `md:"enableTLS"`
	ServerCert string `md:"serverCert"`
	ServerKey  string `md:"serverKey"`
}

type HandlerSettings struct {
	ServiceName string `md:"serviceName"`
	MethodName  string `md:"methodName"`
}

type Output struct {
	Params   map[string]interface{} `md:"params"`
	GrpcData map[string]interface{} `md:"grpcData"`
	Content  interface{}            `md:"content"`
}

func (o *Output) FromMap(values map[string]interface{}) error {
	var err error
	o.Params, err = coerce.ToObject(values["params"])
	if err != nil {
		return err
	}
	o.GrpcData, err = coerce.ToObject(values["grpcData"])
	if err != nil {
		return err
	}
	o.Content = values["content"]

	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"params":   o.Params,
		"grpcData": o.GrpcData,
		"content":  o.Content,
	}
}

type Reply struct {
	Code int         `md:"code"`
	Data interface{} `md:"data"`
}

func (r *Reply) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"code": r.Code,
		"data": r.Data,
	}
}

func (r *Reply) FromMap(values map[string]interface{}) error {

	var err error
	if _, ok := values["code"]; ok {
		r.Code, err = coerce.ToInt(values["code"])
		if err != nil {
			return err
		}
	}
	r.Data, _ = values["data"]

	return nil
}
