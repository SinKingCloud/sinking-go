FROM golang as sinking-consul-builder

#作者
MAINTAINER sinkingcloud 1178710004@qq.com

#修改时区
RUN echo "Asia/Shanghai" > /etc/timezone && rm /etc/localtime && dpkg-reconfigure -f noninteractive tzdata

#修改环境变量
ENV GO111MODULE on
ENV CGO_ENABLED 0
ENV GOPROXY=https://goproxy.cn

#复制代码
COPY . /data/sinking-consul

#指定目录
WORKDIR /data/sinking-consul

#安装依赖
RUN go mod tidy

#编译
RUN go build

FROM alpine as sinking-consul-run

# 设置工作目录
WORKDIR /data/sinking-consul

# 复制文件
COPY --from=sinking-consul-builder /data/sinking-consul/sinking-consul .
COPY --from=sinking-consul-builder /data/sinking-consul/config/ ./config/

#映射端口
EXPOSE 817

#启动
ENTRYPOINT  ["./sinking-consul"]

