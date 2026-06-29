# 编译与部署文档

## 本地开发

### 前置条件

- Node.js 18+
- Go 1.21+
- PostgreSQL 16（或 Docker）
- Redis 7（或 Docker）

### 启动数据库（Docker 方式）

```bash
docker run -d --name openpanda-db \
  -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=openpanda \
  -p 5432:5432 postgres:16-alpine

docker run -d --name openpanda-redis \
  -p 6379:6379 redis:7-alpine
```

### 启动后端

```bash
cd Backend
go run main.go
# 启动在 http://localhost:8080
```

首次启动会自动：
- 连接 PostgreSQL 并创建表结构
- 插入默认三大分类（嵌入式Linux / 硬件电路设计 / 单片机开发）

### 启动前端

```bash
cd Frontend
npm install
npm run dev
# 启动在 http://localhost:3000，自动代理 /api → localhost:8080
```

## npm 命令速查（根目录）

```bash
# 开发
npm run dev:frontend     # 启动前端 Vite HMR
npm run dev:backend      # 编译 Windows 二进制并启动后端

# 编译
npm run build:frontend          # 前端生产构建
npm run build:backend           # 交叉编译 Linux 二进制（Docker 用）
npm run build:backend:local     # 编译 Windows 二进制（本地调试）

# Docker 本地启动
npm run docker:up        # docker-compose up -d（本地构建）
npm run docker:down      # 停止
npm run docker:logs      # 查看日志

# Docker 构建与推送
npm run docker:build              # 构建前后端镜像
npm run docker:build:backend      # 构建后端镜像
npm run docker:build:frontend     # 构建前端镜像
npm run docker:push               # 推送镜像到 DockerHub
npm run docker:release            # 一键构建 + 推送
```

## Docker 镜像构建流程

```
npm run docker:release

  ├── npm run docker:build:backend
  │   ├── npm run build:backend       ← 交叉编译 Linux 二进制
  │   │   CGO_ENABLED=0 GOOS=linux go build → Backend/server
  │   └── docker build -t clintonluo/openpanda-backend ./Backend
  │       └── COPY server → alpine → 最终镜像 ~15MB
  │
  ├── npm run docker:build:frontend
  │   ├── npm run build:frontend      ← vite build → dist/
  │   └── docker build -t clintonluo/openpanda-frontend ./Frontend
  │       ├── Stage1: node → npm install → vite build
  │       └── Stage2: nginx + dist/ + nginx.conf
  │
  └── npm run docker:push
      ├── docker push clintonluo/openpanda-backend
      └── docker push clintonluo/openpanda-frontend
```

## 两个 Compose 文件的区别

| | `docker-compose.yml` | `docker-compose.prod.yml` |
|---|---|---|
| **用途** | 本地开发调试 | 服务器生产部署 |
| **后端** | `build: ./Backend` 从源码编译 Go 并打包 | `image: clintonluo/openpanda-backend` 拉 DockerHub 成品 |
| **前端** | `build: ./Frontend` 从源码 npm build 并打包 | `image: clintonluo/openpanda-frontend` 拉 DockerHub 成品 |
| **需要环境** | Go + Node.js + npm | 仅 Docker |
| **启动速度** | 慢（每次编译） | 快（直接拉镜像） |

## 服务器部署

### 第一次部署

```bash
# 1. 在服务器上创建持久化目录
mkdir -p /data/openpanda/{pgdata,uploads,backups}

# 2. 将部署文件传到服务器
scp docker-compose.prod.yml user@server:/opt/openpanda/

# 3. SSH 到服务器
ssh user@server
cd /opt/openpanda

# 4. 启动
docker-compose -f docker-compose.prod.yml up -d
```

> 所有数据存储在宿主机 `/data/openpanda/` 下，容器删除或软件更新都不会丢失。

### 更新部署

```bash
cd /opt/openpanda
docker-compose -f docker-compose.prod.yml pull
docker-compose -f docker-compose.prod.yml up -d
```

### 数据备份

```bash
# 备份数据库（推荐每日 crontab）
docker exec openpanda-db pg_dump -U postgres openpanda \
  > /data/openpanda/backups/backup-$(date +%Y%m%d).sql

# 恢复数据库
docker exec -i openpanda-db psql -U postgres openpanda \
  < /data/openpanda/backups/backup-20260626.sql

# 备份上传的图片（rsync 到远程或云存储）
rsync -av /data/openpanda/uploads/ user@backup-server:/backups/openpanda-uploads/
```

### 定时备份（crontab）

```bash
# 每天凌晨 2 点自动备份数据库，保留最近 30 天
0 2 * * * docker exec openpanda-db pg_dump -U postgres openpanda > /data/openpanda/backups/backup-$(date +\%Y\%m\%d).sql && find /data/openpanda/backups/ -name '*.sql' -mtime +30 -delete
```

### 端口说明

| 服务 | 端口 | 用途 |
|------|------|------|
| Nginx (Frontend) | 80 | 用户访问入口 |
| Gin (Backend) | 8080 | REST API |
| PostgreSQL | 5432 | 数据库 |
| Redis | 6379 | 缓存 |

### 环境变量（docker-compose.prod.yml 中配置）

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `JWT_SECRET` | 需修改 | JWT 签名密钥，生产务必更换 |
| `ADMIN_USERNAME` | admin | 管理员用户名 |
| `ADMIN_PASSWORD` | !Wo3158023 | 管理员密码，生产务必更换 |
| `DB_PASSWORD` | postgres | 数据库密码 |
| `REDIS_PASSWORD` | (空) | Redis 密码 |
