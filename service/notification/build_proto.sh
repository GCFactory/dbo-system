#!/bin/bash

protoc --go_out=./gen_proto  --go-grpc_out=./gen_proto  -I ./proto  ./proto/notification/notification.proto