# syntax=docker/dockerfile:1

FROM golang:1.23.3 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN  go build -o /app/server ./cmd/main/main.go
RUN  go build -o /app/migrate ./cmd/migrate/main.go

FROM debian:bookworm-slim

WORKDIR /app

# Копируем бинарники с абсолютных путей
COPY --from=builder /app/server /app/server
COPY --from=builder /app/migrate /app/migrate
COPY .env /app/.env
COPY . .
RUN chmod +x /app/server /app/migrate

EXPOSE 8888

CMD ["/app/server"]