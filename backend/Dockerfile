# Multi-stage Dockerfile for Google Cloud Run (zero idle baseline cost)

# Stage 1: Build binary
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install ca-certificates for secure outbound TLS
RUN apk add --no-cache ca-certificates git

# Cache dependencies
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o /app/server ./cmd/server/main.go

# Stage 2: Minimal runtime image
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/server /server

# Cloud Run default port is 8080
ENV PORT=8080
ENV APP_ENV=production

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/server"]
