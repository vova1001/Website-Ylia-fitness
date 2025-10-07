# Dockerfile
FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем из cmd директории
RUN go build -o main ./cmd

EXPOSE 8080

CMD ["./main"]