FROM golang:1.22.5 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY .env .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Скачиваем скрипт wait-for-it.sh
RUN curl -o /wait-for-it.sh https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh
RUN chmod +x /wait-for-it.sh

FROM alpine:latest

RUN apk --no-cache add ca-certificates bash

WORKDIR /app

COPY --from=builder /app .
COPY --from=builder /app/main .
COPY --from=builder /wait-for-it.sh /wait-for-it.sh

RUN mkdir -p pkg/logs

EXPOSE 8080

RUN chmod +x ./main

CMD ["/bin/bash", "/wait-for-it.sh", "rabbit:5672", "--", "./main"]
