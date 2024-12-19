# docker build -t auth_service -f docker/auth.dockerfile .
# docker run -d -p 50051:50051 --name auth auth_service

# Этап сборки
FROM golang:1.23-alpine AS build

# Установка зависимостей для сборки
RUN apk add --no-cache git build-base libwebp-dev musl-dev gcc

WORKDIR /src

# Сначала копируем файлы go.mod и go.sum, чтобы использовать кеш для зависимостей
COPY go.mod go.mod
COPY go.sum go.sum

# Загрузка зависимостей (если go.mod и go.sum не менялись, этот шаг закешируется)
RUN go mod download

# Копируем только необходимые файлы для сборки
COPY . .

# Сборка проекта
RUN go mod tidy
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/auth cmd/auth/main.go

# --------------------------------------------

# Этап запуска
FROM alpine:latest
WORKDIR /app

# Копируем только то, что нужно для выполнения
COPY --from=build /app/auth /app/
COPY config/config.yaml /app/config/config.yaml
COPY ./static_files /app/static_files/

# Устанавливаем команду запуска
CMD ["./auth"]