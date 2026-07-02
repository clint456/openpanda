// ============================================================
// 文件: utils/index.ts
// 说明: 通用工具函数
//       项目中可复用的工具函数集中管理
//       拓展方式：新增函数直接在此添加，或按模块拆分文件
// ============================================================

/**
 * 格式化日期
 * @param dateStr ISO日期字符串
 * @param locale 语言（zh-CN / en-US）
 * @returns 格式化后的日期字符串
 *
 * 示例：
 *   formatDate('2024-01-15T10:30:00Z', 'zh-CN') => '2024年1月15日'
 */
export function formatDate(dateStr: string, locale: string = 'zh-CN'): string {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return d.toLocaleDateString(locale, {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

/**
 * 截断文本（超过指定长度加省略号）
 * @param text 原始文本
 * @param maxLength 最大长度
 * @returns 截断后的文本
 *
 * 示例：
 *   truncate('这是一段很长的文本内容', 5) => '这是一段很...'
 */
export function truncate(text: string, maxLength: number): string {
  if (!text || text.length <= maxLength) return text
  return text.slice(0, maxLength) + '...'
}

/**
 * 去除 HTML 标签，获取纯文本
 * @param html HTML字符串
 * @returns 纯文本
 *
 * 示例：
 *   stripHtml('<p>Hello <b>World</b></p>') => 'Hello World'
 */
export function stripHtml(html: string): string {
  if (!html) return ''
  return html.replace(/<[^>]*>/g, '')
}

/**
 * 生成 SEO 友好的文章 URL
 * 格式：/articles/{id}-{slug}
 * 如果 slug 为空则退化为 /articles/{id}
 *
 * 示例：getArticleUrl({ id: 1, slug: 'hello-world' }) => '/articles/1-hello-world'
 */
export function getArticleUrl(article: { id: number; slug?: string }): string {
  const slug = article.slug || ''
  return slug ? `/articles/${article.id}-${slug}` : `/articles/${article.id}`
}

/**
 * 从文章 URL slug 参数中提取数字 ID
 * "/articles/1-hello-world" → 1
 * "/articles/123" → 123
 */
export function parseArticleId(slug: string): number {
  const match = slug.match(/^(\d+)/)
  return match ? Number(match[1]) : 0
}

/**
 * 防抖函数（限制函数调用频率）
 * 用途：搜索输入框、窗口 resize 等高频触发场景
 *
 * @param fn 需要防抖的函数
 * @param delay 延迟时间（毫秒）
 * @returns 防抖后的函数
 *
 * 示例：
 *   const debouncedSearch = debounce(searchApi, 300)
 *   input.addEventListener('input', debouncedSearch)
 */
export function debounce<T extends (...args: unknown[]) => unknown>(
  fn: T,
  delay: number = 300
): (...args: Parameters<T>) => void {
  let timer: ReturnType<typeof setTimeout> | null = null
  return (...args: Parameters<T>) => {
    if (timer) clearTimeout(timer)
    timer = setTimeout(() => {
      fn(...args)
    }, delay)
  }
}

/**
 * 节流函数（固定时间间隔内只执行一次）
 * 用途：滚动加载、按钮防重复点击等
 *
 * @param fn 需要节流的函数
 * @param interval 时间间隔（毫秒）
 * @returns 节流后的函数
 */
export function throttle<T extends (...args: unknown[]) => unknown>(
  fn: T,
  interval: number = 300
): (...args: Parameters<T>) => void {
  let lastTime = 0
  return (...args: Parameters<T>) => {
    const now = Date.now()
    if (now - lastTime >= interval) {
      lastTime = now
      fn(...args)
    }
  }
}
