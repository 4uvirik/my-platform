FROM golang:1.24-alpine AS builder

WORKDIR /app
# Копируем go.mod конкретного сервиса
COPY services/notification-service/go.mod services/notification-service/go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o notification-service ./services/notification-service/cmd/notification-service

FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/notification-service .
COPY services/notification-service/config ./config

EXPOSE 9100

CMD ["./notification-service"]
