# 🐼 OpenPanda · 开源熊猫

个人技术交流平台 — 嵌入式 Linux / 硬件电路 / 单片机开发技术博客。

**技术栈**: Vue3 + TypeScript + Vite + Element Plus（前端） · Golang Gin + GORM + JWT（后端） · PostgreSQL + Redis（数据层）

---

## ✨ 功能

- 📝 **Markdown 撰写** — 所见即所得编辑器，支持粘贴/拖入图片自动上传
- 📂 **专栏分类** — 动态分类体系，页面交互式管理，无需改代码
- 🔍 **文章搜索** — 标题+正文模糊搜索，分类筛选
- 🌐 **中英双语** — 全站国际化，一键切换，偏好记忆
- 🌙 **夜间模式** — 日间暖橙 / 夜间深色，自动持久化
- 🔐 **JWT 认证** — 登录保护管理功能，Token 自动管理
- 🐳 **Docker 部署** — 前后端独立镜像 + compose 一键编排

---

## 🚀 快速开始

```bash
# 1. 启动数据库
docker run -d --name openpanda-db -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=openpanda -p 5432:5432 postgres:16-alpine
docker run -d --name openpanda-redis -p 6379:6379 redis:7-alpine

# 2. 启动后端
cd Backend && go run main.go

# 3. 启动前端
cd Frontend && npm install && npm run dev
```

浏览器打开 `http://localhost:3000`，默认管理员 `admin` / `!Wo3158023`。

---

## 🐳 Docker 部署

```bash
# 构建并推送镜像
npm run docker:release

# 服务器上启动
docker-compose -f docker-compose.prod.yml up -d
```

---

## 📋 npm 命令

```bash
npm run dev:frontend          # 启动前端开发
npm run dev:backend           # 启动后端（自动编译）

npm run build:backend         # 交叉编译 Linux 二进制
npm run build:frontend        # Vite 生产构建

npm run docker:build          # 构建 Docker 镜像
npm run docker:push           # 推送 DockerHub
npm run docker:release        # 一键构建 + 推送
npm run docker:up             # compose 启动
npm run docker:down           # compose 停止
```

---

## 📁 项目结构

```
openpanda/
├── Backend/                # Golang Gin 后端
│   ├── main.go             # 入口
│   ├── config/             # 配置（环境变量）
│   ├── router/             # 路由注册
│   ├── controller/         # HTTP 处理层
│   ├── service/            # 业务逻辑层
│   ├── model/              # GORM 数据模型
│   ├── middleware/         # CORS / JWT
│   └── utils/              # 统一返回格式
├── Frontend/               # Vue3 + TS 前端
│   └── src/
│       ├── api/            # Axios + API 模块
│       ├── stores/         # Pinia 状态管理
│       ├── router/         # Vue Router + 守卫
│       ├── views/          # 页面（Home/Article/Category/Login）
│       ├── components/     # 通用组件
│       ├── i18n/           # 中英文语言包
│       ├── types/          # TS 类型定义
│       └── utils/          # 工具函数
├── docs/                   # 技术文档
│   ├── architecture.md     # 架构设计
│   ├── design.md           # 功能设计 + 数据库
│   ├── build-and-deploy.md # 编译部署
│   ├── data-safety.md      # 数据安全与备份
│   ├── backend-api.md      # 后端 API 文档
│   └── frontend-api.md     # 前端模块文档
├── docker-compose.yml          # 本地开发编排
├── docker-compose.prod.yml     # 生产部署编排
└── package.json                # npm 命令
```

---

## 📖 文档

| 文档 | 链接 |
|------|------|
| 架构设计 | [docs/architecture.md](docs/architecture.md) |
| 功能设计 + 数据库 | [docs/design.md](docs/design.md) |
| 编译与部署 | [docs/build-and-deploy.md](docs/build-and-deploy.md) |
| 数据安全与备份 | [docs/data-safety.md](docs/data-safety.md) |
| 后端 API | [docs/backend-api.md](docs/backend-api.md) |
| 前端 API 模块 | [docs/frontend-api.md](docs/frontend-api.md) |
| 版本发布命名 | [/docs/versionRule.md](docs/versionRule.md) |

---

## 🔧 技术栈

| 层 | 技术 |
|----|------|
| 前端框架 | Vue 3 + TypeScript（严格模式） |
| 构建 | Vite 5 |
| UI | Element Plus |
| 状态 | Pinia |
| 路由 | Vue Router 4 |
| HTTP | Axios（泛型封装） |
| i18n | vue-i18n |
| 编辑器 | md-editor-v3 |
| 后端框架 | Gin |
| ORM | GORM v2 |
| 认证 | JWT（golang-jwt） |
| 数据库 | PostgreSQL 16 |
| 缓存 | Redis 7 |
| 部署 | Docker + Nginx |

## 📄 License

MIT
