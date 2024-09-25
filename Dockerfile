FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git make curl bash

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o medium ./cmd/main.go
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/docs/ ./docs/
COPY --from=builder /app/migrations/ ./migrations/
COPY --from=builder /app/medium .

CMD ["./medium"]

