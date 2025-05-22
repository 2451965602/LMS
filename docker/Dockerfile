# 第一阶段：构建应用程序
FROM golang:1.24 AS builder

# 设置环境变量
ENV TZ=Asia/Shanghai
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOPROXY=https://goproxy.cn,direct
ENV GOOS=linux
ENV GOARCH=amd64
# 创建工作目录c
RUN mkdir -p /app
RUN mkdir -p /app/output

WORKDIR /app

# 复制所有文件到工作目录
COPY . .

# 下载依赖
RUN go mod tidy

# 编译应用程序
RUN bash ./build.sh

# 第二阶段：创建最终运行环境
FROM alpine

# 环境变量
ENV TZ=Asia/Shanghai
ENV SERVICE=LMS

# 换源，更新软件依赖
RUN sed -i 's#https\?://dl-cdn.alpinelinux.org/alpine#https://mirrors.tuna.tsinghua.edu.cn/alpine#g' /etc/apk/repositories
RUN apk update --no-cache && apk --no-cache add ca-certificates tzdata bash

# 创建工作目录
WORKDIR /app

COPY --from=builder /app/output /app/output
COPY --from=builder /app/config /app/config

CMD ["sh","./output/bootstrap.sh"]


