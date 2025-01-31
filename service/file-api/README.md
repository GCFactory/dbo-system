# File api service


## Commands
### Gen protobuf \
**Windows**
```shell
protoc --go_out=.\proto  --go-grpc_out=.\proto  -I .\proto  .\proto\*  --proto_path=googleapis
```
**Linux**
```shell
protoc --go_out=./proto  --go-grpc_out=./proto  -I ./proto  ./proto/*  --proto_path=googleapis
```

### Gen certificates \
```shell
openssl ecparam -out private.pem -name secp384r1  -genkey
openssl req -new -x509 -nodes -days 3650 -key private.pem -out cert.crt
```

```shell
-----
Country Name (2 letter code) [AU]:UA
State or Province Name (full name) [Some-State]:Kiev
Locality Name (eg, city) []:Kiev
Organization Name (eg, company) [Internet Widgits Pty Ltd]:GCFactory Co
Organizational Unit Name (eg, section) []:
Common Name (e.g. server FQDN or YOUR name) []:file-api
Email Address []:file-api.dbo@gcfactory.space
```