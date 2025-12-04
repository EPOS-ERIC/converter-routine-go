FROM golang:1.25.5-alpine AS builder

# Install required tools
RUN apk add --no-cache git

WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Install swag for OpenAPI generation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy source code
COPY . .

# Generate OpenAPI spec
RUN swag init -o ./server --outputTypes json && \
    mv ./server/swagger.json ./server/openapi.json

# Build - use ./... to find main packages automatically
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o converter-routine ./...

FROM alpine:3.20

RUN apk --no-cache add openjdk21 python3 ca-certificates tzdata

LABEL authors="marcosalvi,valeriovinciarelli"

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /opt/converter

COPY --from=builder /build/converter-routine converter-routine

# Set ownership
RUN chown -R appuser:appgroup /opt/converter

USER appuser:appgroup

CMD ["./converter-routine"]