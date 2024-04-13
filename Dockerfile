FROM golang:1.21
WORKDIR /app
COPY go.* .
RUN go mod download
COPY . .
RUN go build -o banner_service ./cmd/main.go
CMD ["./banner_service"]
