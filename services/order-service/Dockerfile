FROM golang:1.24-alpine AS builder

WORKDIR /app
# Копируем ВСЁ чтобы сохранить структуру монорепо
COPY . .

WORKDIR /app/services/order-service
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o order-service ./cmd/order-service

FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/services/order-service/order-service .
COPY services/order-service/config ./config

EXPOSE 50051
EXPOSE 9100

CMD ["./order-service"]
