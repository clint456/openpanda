# 前端 API 模块文档

## API 层架构

```
src/api/
├── index.ts          # Axios 实例（拦截器/Token/错误处理）
├── types/
│   └── common.ts     # 通用响应类型
└── modules/
    ├── article.ts    # 文章 + 分类 API
    ├── auth.ts       # 认证 API
    └── upload.ts     # 图片上传 API
```

## Axios 实例 (`api/index.ts`)

### 请求拦截器
- 自动从 `localStorage` 读取 Token 附加到请求头
- 自动附加语言偏好 `Accept-Language`

### 响应拦截器
- 统一处理 `code !== 200` 业务错误
- 401 自动清除 Token
- HTTP 错误分类提示（400/401/403/404/500/超时/网络）

### 使用方式

```typescript
import http from '@/api'
import type { ApiResponse, Article } from '@/types'

// GET
const { data } = await http.get<ApiResponse<Article>>('/articles/1')

// POST
await http.post<ApiResponse<Article>>('/admin/articles', formData)

// Upload
const formData = new FormData()
formData.append('file', file)
await http.post('/admin/upload/image', formData, {
  headers: { 'Content-Type': 'multipart/form-data' }
})
```

---

## auth.ts — 认证模块

| 函数 | 方法 | 路径 | 认证 |
|------|------|------|------|
| `login(params)` | POST | `/login` | 否 |
| `getMe()` | GET | `/auth/me` | 是 |

```typescript
import { login, getMe } from '@/api/modules/auth'

// 登录
const { data } = await login({ username: 'admin', password: '!Wo3158023' })
// data.data.token → JWT Token

// 验证
const { data } = await getMe()
// data.data → { user_id, username, role }
```

---

## article.ts — 文章 + 分类模块

| 函数 | 方法 | 路径 | 认证 |
|------|------|------|------|
| `getArticles(params)` | GET | `/articles` | 否 |
| `getArticleById(id)` | GET | `/articles/:id` | 否 |
| `getHotArticles(limit)` | GET | `/articles/hot` | 否 |
| `searchArticles(params)` | GET | `/articles/search` | 否 |
| `createArticle(data)` | POST | `/admin/articles` | 是 |
| `updateArticle(id, data)` | PUT | `/admin/articles/:id` | 是 |
| `deleteArticle(id)` | DELETE | `/admin/articles/:id` | 是 |
| `getCategories()` | GET | `/categories` | 否 |
| `createCategory(data)` | POST | `/admin/categories` | 是 |
| `updateCategory(id, data)` | PUT | `/admin/categories/:id` | 是 |
| `deleteCategory(id)` | DELETE | `/admin/categories/:id` | 是 |

```typescript
import { getArticles, createArticle, getCategories } from '@/api/modules/article'

// 文章列表
const { data } = await getArticles({ page: 1, page_size: 10, category_id: 1 })

// 创建文章
await createArticle({
  title: 'STM32入门',
  content: '# Hello\n正文...',
  category_id: 3,
  language: 'zh'
})

// 获取分类
const { data } = await getCategories()
```

---

## upload.ts — 上传模块

| 函数 | 方法 | 路径 | 认证 |
|------|------|------|------|
| `uploadImage(file)` | POST | `/admin/upload/image` | 是 |

```typescript
import { uploadImage } from '@/api/modules/upload'

// 编辑器自动调用，返回 URL 后插入 Markdown
const { data } = await uploadImage(file)
// data.data.url → "http://xxx/uploads/2026/06/abc.png"
```

---

## 类型定义 (`types/index.ts`)

```typescript
interface Article {
  id: number
  title: string
  slug: string
  content: string     // Markdown 源码
  summary: string
  cover_image: string
  category_id: number
  view_count: number
  is_published: boolean
  language: 'zh' | 'en' | 'both'
  category?: Category
  tags?: Tag[]
  created_at: string
  updated_at: string
}

interface Category {
  id: number
  name: string
  slug: string
  description: string
  sort_order: number
}

interface PaginatedData<T> {
  list: T[]
  total: number
  page: number
  page_size: number
}

interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}
```

---

## 新增 API 模块的步骤

1. 在 `api/modules/` 下新建文件（如 `comment.ts`）
2. 导入 `http` 实例和类型
3. 编写导出函数
4. 在页面组件中 `import` 使用

```typescript
// api/modules/comment.ts
import http from '@/api'
import type { ApiResponse, PaginatedData } from '@/types'

export function getComments(articleId: number, page = 1) {
  return http.get<ApiResponse<PaginatedData<Comment>>>('/comments', {
    params: { article_id: articleId, page }
  })
}
```
