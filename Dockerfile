FROM golang:1.26.3-alpine AS builder

WORKDIR /app

# Copy go.mod and download dependencies
COPY go.mod ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o garnet ./cmd/server/main.go

# Use a lightweight alpine image for the final container
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/garnet .

EXPOSE 6379

# Boot in async mode by default
CMD ["./garnet", "--mode=async"]
