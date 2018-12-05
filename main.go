package main

import (
	"flag"

	"github.com/project-flogo/grpc/support"
)

var (
	packageName = flag.String("package", "main", "package name")
	output      = flag.String("output", ".", "name of output directory")
	input       = flag.String("input", "", "location of the proto file")
)

func main() {
	flag.Parse()

	support.AssignValues(*output)
	err := support.GenerateSupportFiles(*packageName, *input)
	if err != nil {
		panic(err)
	}
}
