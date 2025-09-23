FROM --platform=linux/amd64 golang:1.24 AS builder

# 设置工作目录
WORKDIR /app

# 复制整个 server 目录
COPY ./server/ ./

# 设置环境变量
ENV GOPROXY=https://goproxy.cn
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# 下载依赖并构建
RUN go mod tidy && go build -o sinking-consul main.go

# 运行镜像
FROM --platform=linux/amd64 alpine:latest

# 设置国内时区
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 设置工作目录
WORKDIR /app

# 从构建阶段复制可执行文件和相关文件
COPY --from=builder /app/sinking-consul .

# 暴露端口
EXPOSE 5678

# 启动命令
CMD ["./sinking-consul", "run"]
