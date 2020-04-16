FROM golang:latest AS builder

ENV APP_HOME /go/src/meow

RUN mkdir -p $APP_HOME

WORKDIR $APP_HOME

COPY src/ .

RUN go mod download
RUN go mod verify
RUN go get github.com/githubnemo/CompileDaemon

EXPOSE 8000

ENTRYPOINT CompileDaemon --build="go build -o bin/meow" --command="./bin/meow"