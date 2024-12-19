#docker build -t cart_purchase -f docker/cart_purchase.dockerfile .
#docker run -d -p 50053:50053 --name cart_purchase cart_purchase

# Этап сборки
FROM golang:1.23-alpine AS build

# Установка зависимостей для сборки
RUN apk add --no-cache git build-base libwebp-dev musl-dev gcc

WORKDIR /src

# Кэширование зависимостей
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

# Копируем только необходимые файлы для сборки
COPY . .

# Сборка проекта
RUN go mod tidy
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/cart_purchase cmd/cart_purchase/main.go

# Этап запуска
FROM alpine:latest

# Добавление минимально необходимых пакетов для выполнения
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Копируем только то, что нужно для выполнения
COPY --from=build /app/cart_purchase /app/
COPY config/config.yaml /app/config/config.yaml
COPY ./static_files /app/static_files/

# Устанавливаем команду запуска
CMD ["./cart_purchase"]