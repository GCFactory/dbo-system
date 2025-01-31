FROM golang

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download && go mod verify
COPY cmd ./cmd
COPY config ./config
COPY gen_proto ./gen_proto
COPY internal ./internal
COPY migration ./migration
COPY pkg ./pkg
COPY proto ./proto

RUN go build -v -o /usr/local/bin/app ./cmd/api/main.go

ENV config=config/config

CMD ["app"]