FROM golang:1.22.8-alpine AS builder

WORKDIR /app

COPY backend-gateway/go.mod backend-gateway/go.sum ./

RUN go mod download

COPY ../. .

WORKDIR /app/backend-gateway/cmd/

RUN go build -o gateway

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/backend-gateway/cmd/gateway .
COPY config.yaml .

RUN chmod +x gateway

ENTRYPOINT ["./gateway"]