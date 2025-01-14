FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ../../. .

WORKDIR /app/cmd/

RUN go build -o gateway

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/cmd/gateway .
COPY config.yaml .

RUN chmod +x gateway

ENTRYPOINT ["./gateway"]