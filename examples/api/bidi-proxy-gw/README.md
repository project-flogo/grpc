# bidirectional-proxy-gateway
This example designed to evaluate bidirectional proxy gateway performance using gRPC bidirectional streaming. The messages stream from gRPC client-->Bidirectional Proxy Gateway-->gRPC Echo Server-->Bidirectional Proxy Gateway-->gRPC Client.

The gRPC Client outputs the Number of messages sent and received, Average Response Time taken by all the requests, Moving Average Response Time and Transactions per second values.

## Installation
* The Go programming language version 1.11 or later should be [installed](https://golang.org/doc/install).

## Setup
```bash
git clone https://github.com/project-flogo/grpc
cd grpc/examples/api/bidi-proxy-gw
go build
```

## Performance Testing
1)Start proxy gateway:

Note: If the gRPC server is running in a different system please set the GRPCSERVERIP variable to corresponding system IP before starting the proxy gateway.
```bash
export GRPCSERVERIP=XX.XX.XX.XX
```

```bash
ulimit -n 1048576
export FLOGO_LOG_LEVEL=ERROR
FLOGO_RUNNER_TYPE=DIRECT ./bidi-proxy-gw
```

2)Start sample gRPC server in a new terminal.
```bash
ulimit -n 1048576
./bidi-proxy-gw -server
```

3)Start the gRPC client in a new terminal.

Note: If the proxy gateway is running in a different system please set the GWIP variable to corresponding system IP before starting the client.

```bash
export GWIP=XX.XX.XX.XX
```

```bash
ulimit -n 1048576
./bidi-proxy-gw -client -port 9096 -number 465
```

4)Stop the gRPC client after the required test duration.
All the outputs are printed in the client Terminal.
