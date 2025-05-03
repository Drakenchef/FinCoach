FROM golang:1.23.3-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/main/main.go

EXPOSE 8888

CMD ["./main"]