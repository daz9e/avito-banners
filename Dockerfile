FROM golang:1.18 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /banner-service .

# Запуск приложения
CMD ["/banner-service"]
