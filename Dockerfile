# Build stage
FROM golang:1.24-alpine AS builder

# Install git and swag dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Install specific version of swag
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.12

# Clean existing docs
RUN rm -rf docs

# Generate swagger docs
RUN cd cmd && swag init -g main.go -o ../docs

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary and necessary files from builder
COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs

# Expose port 8080
EXPOSE 8080

# Command to run the application
CMD ["./main"] 