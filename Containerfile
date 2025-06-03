# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/app

# Build the queue system
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o queue ./cmd/queue

# Server stage
FROM alpine:latest AS server

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/server .

# Copy views directory for HTML templates
COPY --from=builder /app/views ./views

# Copy .env file from builder stage
COPY --from=builder /app/.env .

# Expose port 8080
EXPOSE 8080

# Command to run the application
CMD ["./server"]

# Queue stage
FROM alpine:latest AS queue

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/queue .

# Copy .env file from builder stage
COPY --from=builder /app/.env .

# Expose port 8081
EXPOSE 8081

# Command to run the application
CMD ["./queue"]
