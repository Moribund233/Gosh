# Stage 1: Build
FROM golang:1.26-alpine AS builder

# China network adaptation
RUN go env -w GOPROXY=https://goproxy.cn,direct GO111MODULE=on

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/bin/gosh-server cmd/server/main.go

# Stage 2: Runtime
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app
COPY --from=builder /app/bin/gosh-server .
COPY --from=builder /app/config config/

EXPOSE 8080

ENV GIN_MODE=release
ENV GOSH_SERVER_MODE=release

ENTRYPOINT ["./gosh-server"]
