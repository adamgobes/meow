FROM golang:latest AS builder

ENV APP_HOME /go/src/meow

RUN mkdir -p $APP_HOME

WORKDIR $APP_HOME

COPY src/ .

RUN go mod download
RUN go mod verify

CMD go test -v