# Используйте официальный образ Golang как базовый для компиляции
FROM golang:1.22 AS builder


WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o banner-service ./cmd/

CMD ["./banner-service"]