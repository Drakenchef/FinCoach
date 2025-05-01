FROM golang:1.23.3-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/main/main.go

EXPOSE ${APP_PORT}

CMD ["./main"]