# ---------- Builder ----------
FROM golang:1.25.5-alpine AS builder

WORKDIR /cache

# Установим git для go get/mod
RUN apk add --no-cache git

# Настройка кэша Go
ENV GOPATH=/go
ENV GOCACHE=/go/cache

# Копируем только go.mod и go.sum для кэширования зависимостей
WORKDIR /cache/app
COPY go.mod go.sum ./
RUN go mod download

# ---------- Runner ----------
FROM golang:1.25.5-alpine

WORKDIR /app

# Кэш Go модулей из builder
COPY --from=builder /go /go

# Git нужен для go get/mod
RUN apk add --no-cache git build-base

# Включаем CGO для sqlite3
ENV CGO_ENABLED=1

# Запуск Go напрямую (stdout не буферизуется)
CMD ["go", "run", "./cmd/myapp"]