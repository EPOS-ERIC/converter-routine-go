FROM golang:1.25-alpine AS builder

WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o converter-routine ./cmd/main.go

FROM alpine:3.20

RUN apk --no-cache add openjdk21 python3

LABEL authors="valeriovinciarelli"

WORKDIR /opt/converter

COPY --from=builder /build/converter-routine converter-routine

CMD ["./converter-routine"]