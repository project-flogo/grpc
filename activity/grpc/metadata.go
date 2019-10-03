package grpc

import (
	"github.com/project-flogo/core/data/coerce"
)

// Settings are the jsexec settings
type Settings struct {
	OperatingMode string `md:"operatingMode"`
	HostURL       string `md:"hosturl"`
	EnableTLS     bool   `md:"enableTLS"`
	ClientCert    string `md:"clientCert"`
}

// Input is the input into the javascript engine
type Input struct {
	GRPCMthdParamtrs map[string]interface{} `md:"grpcMthdParamtrs"`
	Header           map[string]string      `md:"header"`
	ServiceName      string                 `md:"serviceName"`
	ProtoName        string                 `md:"protoName"`
	MethodName       string                 `md:"methodName"`
	Params           map[string]string      `md:"params"`
	QueryParams      map[string]string      `md:"queryParams"`
	Content          interface{}            `md:"content"`
	PathParams       map[string]string      `md:"pathParams"`
}

// FromMap converts the values from a map into the struct Input
func (r *Input) FromMap(values map[string]interface{}) error {
	grpcMthdParamtrs, err := coerce.ToObject(values["grpcMthdParamtrs"])
	if err != nil {
		return err
	}
	if grpcMthdParamtrs == nil {
		grpcMthdParamtrs = make(map[string]interface{})
	}
	r.GRPCMthdParamtrs = grpcMthdParamtrs
	header, err := coerce.ToParams(values["header"])
	if err != nil {
		return err
	}
	r.Header = header
	serviceName, err := coerce.ToString(values["serviceName"])
	if err != nil {
		return err
	}
	r.ServiceName = serviceName
	if _, ok := r.GRPCMthdParamtrs["serviceName"]; !ok {
		r.GRPCMthdParamtrs["serviceName"] = serviceName
	}

	protoName, err := coerce.ToString(values["protoName"])
	if err != nil {
		return err
	}
	r.ProtoName = protoName
	if _, ok := r.GRPCMthdParamtrs["protoName"]; !ok {
		r.GRPCMthdParamtrs["protoName"] = protoName
	}

	methodName, err := coerce.ToString(values["methodName"])
	if err != nil {
		return err
	}
	r.MethodName = methodName
	if _, ok := r.GRPCMthdParamtrs["methodName"]; !ok {
		r.GRPCMthdParamtrs["methodName"] = methodName
	}

	params, err := coerce.ToParams(values["params"])
	if err != nil {
		return err
	}
	r.Params = params
	queryParams, err := coerce.ToParams(values["queryParams"])
	if err != nil {
		return err
	}
	r.QueryParams = queryParams
	r.Content = values["content"]
	pathParams, err := coerce.ToParams(values["pathParams"])
	if err != nil {
		return err
	}
	r.PathParams = pathParams
	return nil
}

// ToMap converts the struct Input into a map
func (r *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"grpcMthdParamtrs": r.GRPCMthdParamtrs,
		"header":           r.Header,
		"serviceName":      r.ServiceName,
		"protoName":        r.ProtoName,
		"methodName":       r.MethodName,
		"params":           r.Params,
		"queryParams":      r.QueryParams,
		"content":          r.Content,
		"pathParams":       r.PathParams,
	}
}

// Output is the ouput from the grpc request
type Output struct {
	Body interface{} `md:"body"`
}

// FromMap converts the values from a map into the struct Output
func (o *Output) FromMap(values map[string]interface{}) error {
	o.Body = values["body"]
	return nil
}

// ToMap converts the struct Output into a map
func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"body": o.Body,
	}
}
