FROM golang:latest AS builder
COPY . /app
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -a -o /meow ./src

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /meow ./
RUN chmod +x ./meow
EXPOSE 8000
CMD ./meow