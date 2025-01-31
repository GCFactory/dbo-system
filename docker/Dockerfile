# Initial stage: download modules
FROM golang:1.21-alpine as builder

ENV config=docker

WORKDIR /app

COPY ./ /app

RUN go mod download

RUN go build -o service cmd/api/main.go


# Intermediate stage: Build the binary
FROM golang:1.21-alpine

ENV config=docker

WORKDIR /app

RUN addgroup -g 3001 gogroup \
    && adduser -H -D gouser -u 3002 -G gogroup

COPY --from=builder --chown=gouser:gogroup --chmod=0755 /app/service /app