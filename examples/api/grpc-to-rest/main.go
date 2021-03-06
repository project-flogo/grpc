package main

import (
	"flag"

	"github.com/project-flogo/core/engine"
	"github.com/project-flogo/grpc/examples"
	"github.com/project-flogo/grpc/proto/grpc2rest"
)

func main() {
	clientMode := flag.Bool("client", false, "command to run client")
	methodName := flag.String("method", "pet", "method name")
	paramVal := flag.String("param", "user2", "method param")
	port := flag.String("port", "", "port value")
	flag.Parse()

	if *clientMode {
		grpc2rest.CallClient(port, methodName, *paramVal)
		return
	}

	e, err := examples.GRPC2RestExample()
	if err != nil {
		panic(err)
	}
	engine.RunEngine(e)
}
