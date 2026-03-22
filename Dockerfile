FROM golang:1.21-alpine

WORKDIR /app

# 安装必要的依赖
RUN apk add --no-cache git

# 复制源代码
COPY task-management/ .

# 安装依赖并编译
RUN go mod download && \
    go build -o task-skill .

# 默认命令
ENTRYPOINT ["./task-skill", "--server"]
