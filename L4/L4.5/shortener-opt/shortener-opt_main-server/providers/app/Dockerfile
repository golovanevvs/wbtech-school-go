# Build

FROM golang:1.25.4-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

ARG APP_NAME

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o $APP_NAME ./cmd/$APP_NAME

# Runtime

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder /app/$APP_NAME .

COPY --from=builder /app/providers ./providers/

ENTRYPOINT ["/root/shortener-opt_main-server"]