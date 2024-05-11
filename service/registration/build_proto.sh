#!/bin/bash

protoc --go_out=./proto  --go-grpc_out=./proto  -I ./proto  ./proto/registration/registration.proto
protoc --go_out=./proto  --go-grpc_out=./proto  -I ./proto  ./proto/platform/platform.proto