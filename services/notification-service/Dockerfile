FROM golang:1.24-alpine AS builder

WORKDIR /app
# Копируем ВСЁ чтобы сохранить структуру монорепо
COPY . .

WORKDIR /app/services/notification-service

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o notification-service ./cmd/notification-service

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/services/notification-service/notification-service .
COPY services/notification-service/config ./config

EXPOSE 9100

CMD ["./notification-service"]
