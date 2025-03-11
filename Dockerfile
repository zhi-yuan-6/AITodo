#使用Go官方镜像作为构建阶段
FROM golang:1.23.4 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o app main.go

#使用轻量级镜像作为运行阶段
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/app .
RUN mkdir "config"
COPY config/config.yaml ./config/config.yaml
COPY private_key.pem .
EXPOSE 8080
CMD ["./app"]

LABEL authors="纸鸢"




