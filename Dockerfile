# Build stage
FROM golang:1.25.4-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the Go application binary
RUN CGO_ENABLED=0 GOOS=linux go build -o agromart-server .

# Run stage
FROM alpine:3.21

WORKDIR /app

# Install certificates for AWS/HTTPS requests
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/agromart-server .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./agromart-server"]
