# Build
FROM golang:1.24-alpine AS builder
WORKDIR /app
ARG TARGETOS
ARG TARGETARCH

RUN apk add --no-cache git make

# Install goose for migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .
RUN go mod download

RUN if [ -z "$TARGETOS" ]; then \
        if [ "$(uname -s)" = "Darwin" ]; then \
            export TARGETOS="darwin"; \
        else \
            export TARGETOS="linux"; \
        fi; \
    fi && \
    if [ -z "$TARGETARCH" ]; then \
        if [ "$(uname -m)" = "arm64" ] || [ "$(uname -m)" = "aarch64" ]; then \
            export TARGETARCH="arm64"; \
        else \
            export TARGETARCH="amd64"; \
        fi; \
    fi && \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o api ./cmd/api/main.go

# App
FROM alpine:3.17

RUN apk add --no-cache ca-certificates tzdata make && \
    mkdir -p /app/migrations /app/templates

# Copy the API binaries
COPY --from=builder /app/docs /app/docs
COPY --from=builder /app/api /app/api
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY Makefile /app/Makefile

WORKDIR /app

RUN chmod +x /app/api && \
    chmod +x /usr/local/bin/goose

EXPOSE 8080

CMD ["/app/api"]