FROM golang:1.24-alpine AS builder

WORKDIR /app
# Копируем go.mod конкретного сервиса
COPY services/order-service/go.mod services/order-service/go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o order-service ./services/order-service/cmd/order-service

FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/order-service .
COPY services/order-service/config ./config

EXPOSE 50051
EXPOSE 9100

CMD ["./order-service"]
