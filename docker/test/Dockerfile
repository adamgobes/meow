FROM golang:latest AS builder
COPY . /app
WORKDIR /app

RUN go mod download
RUN go get github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon --build="go build -o bin/meow ./src" --command="go test ./src -v"