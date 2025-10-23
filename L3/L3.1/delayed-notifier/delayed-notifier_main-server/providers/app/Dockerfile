# ===== Stage 1: Build =====

FROM golang:1.25.2-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

ARG APP_NAME
ARG APP_TRANSPORT_HTTP_PORT

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o $APP_NAME ./cmd/$APP_NAME

# ===== Stage 2: Runtime =====

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder /app/$APP_NAME .

EXPOSE $APP_TRANSPORT_HTTP_PORT

ENTRYPOINT ["/root/delayed-notifier_main-server"]