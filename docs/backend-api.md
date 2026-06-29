# 后端 API 接口文档

Base URL: `/api/v1`

## 通用返回格式

```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

| code | 含义 |
|------|------|
| 200 | 成功 |
| 400 | 参数错误 |
| 401 | 未授权 |
| 404 | 不存在 |
| 500 | 服务器错误 |

---

## 认证

### POST /login — 登录

> 公开接口

**Request:**
```json
{ "username": "admin", "password": "!Wo3158023" }
```

**Response:**
```json
{
  "code": 200,
  "data": { "token": "eyJhbG...", "username": "admin" }
}
```

### GET /auth/me — 获取当前用户

> 需认证 `Authorization: Bearer <token>`

---

## 文章

### GET /articles — 文章列表

> 公开接口

**Query:**

| 参数 | 类型 | 默认 | 说明 |
|------|------|------|------|
| page | int | 1 | 页码 |
| page_size | int | 10 | 每页条数 |
| category_id | int | 0 | 按分类筛选 |
| tag_id | int | 0 | 按标签筛选 |

**Response:**
```json
{
  "code": 200,
  "data": {
    "list": [{ "id": 1, "title": "...", "category": {...}, "tags": [...] }],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

### GET /articles/:id — 文章详情

### GET /articles/hot — 热门文章

> 按阅读量倒序，默认返回 10 条

### GET /articles/search — 搜索文章

| 参数 | 说明 |
|------|------|
| keyword | 关键词（标题+正文模糊匹配） |
| page | 页码 |
| page_size | 每页条数 |

---

## 分类

### GET /categories — 所有分类

> 公开接口

```json
{
  "code": 200,
  "data": [
    { "id": 1, "name": "嵌入式Linux", "slug": "embedded-linux", "description": "...", "sort_order": 1 }
  ]
}
```

---

## 管理接口（需 JWT 认证）

> 请求头: `Authorization: Bearer <token>`

### POST /admin/articles — 创建文章

**Request:**
```json
{
  "title": "STM32时钟树详解",
  "content": "# 标题\n正文内容（Markdown）",
  "summary": "摘要",
  "cover_image": "https://...",
  "category_id": 3,
  "tag_ids": [1, 2],
  "language": "zh"
}
```

### PUT /admin/articles/:id — 更新文章

> 支持部分更新，只传需要修改的字段

### DELETE /admin/articles/:id — 删除文章

### POST /admin/upload/image — 上传图片

> Content-Type: multipart/form-data  
> 字段名: `file`

**Response:**
```json
{ "code": 200, "data": { "url": "http://localhost:8080/uploads/2026/06/xxx.png" } }
```

### POST /admin/categories — 创建分类

```json
{ "name": "FPGA开发", "slug": "fpga-dev", "description": "...", "sort_order": 4 }
```

### PUT /admin/categories/:id — 更新分类

### DELETE /admin/categories/:id — 删除分类

---

## 健康检查

### GET /health

```json
{ "status": "ok" }
```
