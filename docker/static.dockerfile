# docker build -t static_service -f docker/static.dockerfile .
# docker run -d -p 8081:8081 --name static static_service

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
COPY db/migrations db/migrations
RUN go mod tidy
RUN go build -o static cmd/static/main.go

# --------------------------------------------
# Этап запуска
FROM alpine:latest

WORKDIR /app
COPY --from=build /src/static /app
COPY config/config.yaml /app/config/config.yaml
COPY ./static_files /src/static_files/

CMD ["./static"]