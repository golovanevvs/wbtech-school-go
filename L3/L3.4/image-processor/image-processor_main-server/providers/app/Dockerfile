# Build

FROM golang:1.25.2-alpine AS builder

RUN apk add --no-cache \
    git \
    ca-certificates \
    build-base \
    pkgconfig \
    opencv-dev

WORKDIR /app

ARG APP_NAME

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o $APP_NAME ./cmd/$APP_NAME

# Runtime

FROM alpine:latest

RRUN apk add --no-cache \
    ca-certificates \
    libstdc++ \
    opencv

WORKDIR /root/

COPY --from=builder /app/$APP_NAME .

COPY --from=builder /app/providers ./providers/

RUN mkdir -p ./resources/app/filestorage
RUN chmod 755 ./resources/app/filestorage

ENTRYPOINT ["/root/image-processor_main-server"]