# ============================================================
# Backend Dockerfile - Go + Gin 后端服务
# 多阶段构建：编译阶段 → 运行阶段（最终镜像极小）
# ============================================================

# -------------------- 阶段1: 编译 --------------------
FROM golang:1.21-alpine AS builder

# 设置 Go 环境（国内加速）
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

# 设置工作目录
WORKDIR /app

# 先复制依赖文件（利用 Docker 缓存层）
COPY go.mod go.sum ./
RUN go mod download

# 复制所有源代码
COPY . .

# 编译为静态链接的二进制文件
# CGO_ENABLED=0: 禁用 CGO，生成纯静态二进制
# -ldflags="-w -s": 去除调试信息，减小体积
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server main.go

# -------------------- 阶段2: 运行 --------------------
FROM alpine:3.19

# 安装 CA 证书（HTTPS 请求需要）和时区数据
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

WORKDIR /app

# 从编译阶段复制二进制文件
COPY --from=builder /app/server .

# 暴露端口
EXPOSE 8080

# 运行服务
CMD ["./server"]
