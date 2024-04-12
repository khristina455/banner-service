FROM golang:1.21
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . ./
RUN go build -o banner_service
CMD ["./app/banner_service"]
