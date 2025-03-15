# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install development dependencies
RUN apk add --no-cache \
    git \
    tzdata \
    make \
    gcc \
    musl-dev

# Install development tools
RUN go install github.com/air-verse/air@latest && \
    go install github.com/pressly/goose/v3/cmd/goose@latest && \
    go install github.com/swaggo/swag/cmd/swag@latest

# Copy the rest of the application
COPY . .
RUN go mod download

# Generate swagger docs
RUN swag init -g cmd/api/main.go

EXPOSE 8080

# Use air for hot reload
CMD ["air", "-c", ".air.toml"]