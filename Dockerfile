FROM golang:1.25.5-alpine AS builder

RUN apk add --no-cache git

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

# Generate OpenAPI spec
RUN swag init -o . --outputTypes json && \
    mv ./swagger.json ./openapi.json

# Find and build the main package
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o converter-routine .

FROM alpine:3.20

RUN apk --no-cache add python3 ca-certificates tzdata

LABEL authors="valeriovinciarelli"

RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /opt/converter

COPY --from=builder /build/converter-routine converter-routine

RUN chown -R appuser:appgroup /opt/converter

USER appuser:appgroup

CMD ["./converter-routine"]