FROM golang:latest AS builder

ENV APP_HOME /go/src/meow

RUN mkdir -p $APP_HOME

WORKDIR $APP_HOME

COPY src/ .

RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -a -o /meow

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /meow ./
RUN chmod +x ./meow
EXPOSE 8000
CMD ./meow