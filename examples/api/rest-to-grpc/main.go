package main

import (
	"flag"

	"github.com/project-flogo/core/engine"
	"github.com/project-flogo/grpc/examples"
	"github.com/project-flogo/grpc/proto/rest2grpc"
)

func main() {
	clientMode := flag.Bool("client", false, "command to run client")
	serverMode := flag.Bool("server", false, "command to run server")
	methodName := flag.String("method", "pet", "method name")
	paramVal := flag.String("param", "user2", "method param")
	port := flag.String("port", "", "port value")
	flag.Parse()

	if *clientMode {
		rest2grpc.CallClient(port, methodName, *paramVal)
		return
	} else if *serverMode {
		server, _ := rest2grpc.CallServer()
		server.Start()
		return
	}

	e, err := examples.Rest2GRPCExample()
	if err != nil {
		panic(err)
	}
	engine.RunEngine(e)
}
