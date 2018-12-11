package main

import (
	"flag"

	"github.com/project-flogo/grpc/proto/grpc2grpc"
)

func main() {
	clientMode := flag.Bool("client", false, "command to run client")
	serverMode := flag.Bool("server", false, "command to run server")
	methodName := flag.String("method", "pet", "method name")
	paramVal := flag.String("param", "user2", "method param")
	port := flag.String("port", "", "port value")
	flag.Parse()

	if *clientMode {
		grpc2grpc.CallClient(port, methodName, *paramVal, nil)
		return
	} else if *serverMode {
		server, _ := grpc2grpc.CallServer()
		server.Start()
		return
	}
}
