FROM golang:1.21

WORKDIR /app

RUN go get github.com/githubnemo/CompileDaemon && \
  go get github.com/go-delve/delve/cmd/dlv

ENV config=docker

COPY ./ /app

RUN go mod download
EXPOSE 5000 40000

ENTRYPOINT CompileDaemon --build="go build cmd/api/main.go" --command="dlv debug --headless --listen=:40000 --api-version=2 --accept-multiclient cmd/api/main.go"