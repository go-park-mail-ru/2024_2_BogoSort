# docker build -t survey_service -f docker/survey.dockerfile .
# docker run -d -p 8082:8082 --name survey survey_service

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
COPY db/migrations db/migrations
RUN go mod tidy
RUN go build -o survey cmd/survey/main.go

# --------------------------------------------
# Этап запуска
FROM alpine:latest

WORKDIR /app
COPY --from=build /src/survey /app
COPY config/config.yaml /app/config/config.yaml
COPY ./static_files /src/static_files/

CMD ["./survey"]