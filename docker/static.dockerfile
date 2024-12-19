# docker build -t static_service -f docker/static.dockerfile .
# docker run -d -p 8081:8081 --name static static_service

# Этап сборки
FROM golang:1.23-alpine AS build

# Установка зависимостей для сборки
RUN apk add --no-cache git build-base musl-dev gcc

WORKDIR /src

# Сначала копируем файлы go.mod и go.sum, чтобы использовать кеш для зависимостей
COPY go.mod go.sum ./

# Загрузка зависимостей (если go.mod и go.sum не менялись, этот шаг закешируется)
RUN go mod download

# Копируем остальные файлы проекта
COPY . .

# Сборка проекта
RUN go mod tidy
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/static cmd/static/main.go

# Этап запуска
FROM alpine:latest

# Добавление минимально необходимых пакетов для выполнения
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Копируем только то, что нужно для выполнения
COPY --from=build /app/static /app/
COPY config/config.yaml /app/config/config.yaml
COPY static_files /app/static_files/

# Устанавливаем команду запуска
CMD ["./static"]