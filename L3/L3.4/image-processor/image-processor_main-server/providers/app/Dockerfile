# Build
FROM golang:1.25.2-bullseye AS builder

RUN apt-get update && apt-get install -y \
    pkg-config \
    libopencv-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

ARG APP_NAME

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o $APP_NAME ./cmd/$APP_NAME

# Runtime
FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y \
    libopencv-core4.5 \
    libopencv-imgproc4.5 \
    libopencv-imgcodecs4.5 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /root/

COPY --from=builder /app/$APP_NAME .

COPY --from=builder /app/providers ./providers/

RUN mkdir -p ./resources/app/filestorage
RUN chmod 755 ./resources/app/filestorage

ENTRYPOINT ["/root/image-processor_main-server"]