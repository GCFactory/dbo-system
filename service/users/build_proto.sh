#!/bin/bash

protoc --go_out=./gen_proto  --go-grpc_out=./gen_proto  -I ./proto  ./proto/platform/platform.proto
protoc --go_out=./gen_proto  --go-grpc_out=./gen_proto  -I ./proto  ./proto/users/users.proto