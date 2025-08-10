# Stage 1: Build the Go application
FROM golang:1.24-alpine as builder

# Set the working directory inside the container
WORKDIR /app

# Install necessary build tools and dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application with static linking for Alpine compatibility
RUN go build -o main ./cmd/main.go

# Stage 2: Create a minimal image to run the application
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Install necessary runtime dependencies
RUN apk add --no-cache ca-certificates

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./main"]