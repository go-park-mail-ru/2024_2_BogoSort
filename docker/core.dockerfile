# Этап сборки
FROM golang:1.23-alpine AS build

# Установка зависимостей для сборки
RUN apk add --no-cache gcc libc-dev git

WORKDIR /src

# Сначала копируем файлы go.mod и go.sum, чтобы использовать кеш для зависимостей
COPY go.mod go.sum ./

# Загрузка зависимостей (если go.mod и go.sum не менялись, этот шаг закешируется)
RUN go mod download

# Копируем остальные файлы проекта
COPY . .

# Сборка проекта
RUN go build -o /app/core cmd/app/main.go

# --------------------------------------------

# Этап запуска
FROM alpine:latest

# Добавление минимально необходимых пакетов для выполнения
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Копируем только то, что нужно для выполнения
COPY --from=build /app/core /app/
COPY config/config.yaml /app/config/config.yaml
COPY static_files /app/static_files

# Устанавливаем команду запуска
CMD ["./core"]