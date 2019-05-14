package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/project-flogo/core/engine"
	"github.com/project-flogo/grpc/examples"
	"github.com/project-flogo/grpc/proto/bidiproxygw"
)

var ()

func main() {
	//flag.Parse()
	cpuProfile := flag.String("cpuprofile", "", "Writes CPU profile to the specified file")
	memProfile := flag.String("memprofile", "", "Writes memory profile to the specified file")
	clientMode := flag.Bool("client", false, "command to run client")
	serverMode := flag.Bool("server", false, "command to run server")
	methodName := flag.String("method", "pet", "method name")
	number := flag.Int("number", 1, "number of times to run client")
	paramVal := flag.String("param", "user2", "method param")
	port := flag.String("port", "", "port value")
	flag.Parse()

	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			fmt.Println(fmt.Sprintf("Failed to create CPU profiling file due to error - %s", err.Error()))
			os.Exit(1)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	if *memProfile != "" {
		f, err := os.Create(*memProfile)
		if err != nil {
			fmt.Println(fmt.Sprintf("Failed to create memory profiling file due to error - %s", err.Error()))
			os.Exit(1)
		}

		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			fmt.Println(fmt.Sprintf("Failed to write memory profiling data to file due to error - %s", err.Error()))
			os.Exit(1)
		}
		f.Close()
	}

	if *clientMode {
		bidiproxygw.CallClient(port, methodName, *paramVal, nil, *number)
	} else if *serverMode {
		server, _ := bidiproxygw.CallServer()
		server.Start()
		return
	} else {

		e, err := examples.GRPC2GRPCBidiProxyExample()
		if err != nil {
			panic(err)
		}
		engine.RunEngine(e)
	}
}
