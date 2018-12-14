# gRPC Trigger, Activity, and Examples
This repo contains [gRPC](https://en.wikipedia.org/wiki/GRPC) trigger, activity and examples. It is intended to work with the [microgateway](https://github.com/project-flogo/microgateway).

## gRPC utility
The gRPC utility generates the files needed by the trigger and activity from a proto file.

### Install
Run the following in the root of the repo:
```bash
go install
```

### Example
Generates the needed files from petstore.proto and places them in src as the main package:
```bash
grpc -input petstore.proto -output src/ -package main
```
