FROM golang

RUN apt-get update && apt-get install -y graphviz

WORKDIR /usr/src/app

COPY ../go.mod go.sum ./

RUN go mod download && go mod verify
COPY ../cmd ./cmd
COPY ../config ./config
COPY ../internal ./internal

RUN go build -v -o /usr/local/bin/app ./cmd/api/main.go

ENV config=config/config

CMD ["app"]