package main

import (
	"flag"

	"github.com/project-flogo/grpc/support"
)

var (
	name      = flag.String("name", "", "name of application")
	protoPath = flag.String("protoPath", "", "location of the proto file")
)

func main() {
	flag.Parse()

	support.AssignValues(*name)
	err := support.GenerateSupportFiles(*protoPath)
	if err != nil {
		panic(err)
	}
}
