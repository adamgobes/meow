FROM golang:latest AS builder
COPY . /app
WORKDIR /app

RUN go mod download
RUN go get github.com/githubnemo/CompileDaemon

EXPOSE 8000

ENTRYPOINT CompileDaemon --build="go build -o bin/meow ./src" --command="./bin/meow"