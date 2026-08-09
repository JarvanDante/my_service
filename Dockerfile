# 构建：本地一体化入口（front+backend+manage 同进程，便于开发）
FROM golang:1.23-bookworm AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/myservice .

# 运行
FROM debian:bookworm-slim
RUN apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates curl \
  && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /out/myservice /app/myservice
COPY manifest/config/config.docker.yaml /app/manifest/config/config.yaml
EXPOSE 8000
CMD ["/app/myservice"]
