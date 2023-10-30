# File api service


### Commands
Gen protobuf \
**Windows**
```shell
protoc --go_out=.\proto  --go-grpc_out=.\proto  -I .\proto  .\proto\*  --proto_path=googleapis
```
**Linux**
```shell
protoc --go_out=./proto  --go-grpc_out=./proto  -I ./proto  ./proto/*  --proto_path=googleapis
```