#docker build -t cart_purchase -f docker/cart_purchase.dockerfile .
#docker run -d -p 50053:50053 --name cart_purchase cart_purchase

# Этап сборки
FROM golang:1.23-alpine AS build

RUN apk add --no-cache gcc libc-dev git
WORKDIR /src
COPY cmd cmd
COPY internal internal
COPY docs docs
COPY go.mod go.mod
COPY config config
COPY static_files static_files
COPY pkg pkg
RUN go mod tidy
RUN go build -o cart_purchase cmd/cart_purchase/main.go

# --------------------------------------------

# Этап запуска
FROM alpine:latest

WORKDIR /app
COPY --from=build /src/cart_purchase /app
COPY config/config.yaml /app/config/config.yaml
COPY static_files /app/static_files
CMD ["./cart_purchase"]