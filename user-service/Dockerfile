# Build stage
FROM golang:1.23.2-alpine3.19 AS builder

WORKDIR /app

# Cache dependencies by copying go.mod and go.sum first
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the go application
# Using -ldflags="-s -w" removes debugging information from the Go binary, further reducing the size.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o userApp ./main.go

# Run stage
FROM alpine:3.17

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/userApp .

# Set a non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Expose port for the application
EXPOSE 6000

# Command to run the application
CMD [ "/app/userApp" ]