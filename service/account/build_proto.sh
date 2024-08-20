#!/bin/bash

#protoc --go_out=./proto  --go-grpc_out=./proto  -I ./proto  ./proto/platform/platform.proto
protoc --go_out=./proto  --go-grpc_out=./proto  -I ./proto  ./proto/account/account.proto