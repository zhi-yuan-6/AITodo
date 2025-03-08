LABEL authors="纸鸢"
#输出系统中当前的进程信息。
ENTRYPOINT ["top", "-b"]

#使用Go官方镜像作为构建阶段
FROM golang:1.23.4 AS builder
WORKDIR /ai_todo
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o aitodo

#使用轻量级镜像作为运行阶段
FROM alipine:latest
RUN apk -no-cache add ca-certificates
WORKDIR /ai_todo
COPY --from=builder /ai_todo/aitodo /app/aitodo
RUN mkdir "config"
COPY config/config.yaml ./config/config.yaml
COPY private_key.pem .
EXPOSE 8080
CMD ["./aitodo"]


