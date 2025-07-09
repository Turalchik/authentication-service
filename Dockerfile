# 1. Стадия сборки
FROM golang:1.24.2 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o authservice ./cmd/main

# 2. Финальный минимальный образ
FROM alpine:latest
WORKDIR /app

# Копируем только готовый бинарник
COPY --from=builder /app/authservice .

EXPOSE 8080
ENTRYPOINT ["./authservice"]
