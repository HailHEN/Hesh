#  Builder stage. Named or provided with alias as a builder 
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Cache dependencies first 
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build a static, lightweight binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /hesh-core-service ./main.go

# minimal runtime to run core service
FROM alpine:3.20

WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /hesh-core-service .

# Expose application port
EXPOSE 8080

# Execute the binary
CMD ["./hesh-core-service"]