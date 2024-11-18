#docker build -t cart_purchase -f docker/cart_purchase.dockerfile .
#docker run -d -p 50053:50053 --name cart_purchase cart_purchase

# Этап сборки
FROM golang:1.22-alpine AS build

RUN apk add --no-cache gcc
RUN apk add libc-dev
WORKDIR /src
COPY cmd cmd
COPY internal internal
COPY docs docs
COPY go.mod go.mod
COPY config config
COPY pkg pkg
RUN go mod tidy
RUN go build -o cart_purchase cmd/app/main.go

# --------------------------------------------

# Этап запуска
FROM alpine:latest

WORKDIR /app
COPY --from=build /src/cart_purchase /cart_purchase
COPY config/config.yaml /app/config/config.yaml

CMD ["./cart_purchase"]