# docker build -t core_service -f docker/core.dockerfile .
# docker run -d -p 8080:8080 --name core core_service

# Этап сборки
FROM golang:1.23-alpine AS build

RUN apk add --no-cache gcc
RUN apk add libc-dev
WORKDIR /src
COPY cmd cmd
COPY internal internal
COPY docs docs
COPY go.mod go.mod
COPY config config
COPY static_files static_files
COPY pkg pkg
RUN go mod tidy
RUN go build -o core cmd/app/main.go

# --------------------------------------------

# Этап запуска
FROM alpine:latest

WORKDIR /app
COPY --from=build /src/core /app
COPY config/config.yaml /app/config/config.yaml

CMD ["./core"]