// ============================================================
// 包名: middleware
// 说明: 中间件集合
//
//	拓展方式：新建 .go 文件，实现 gin.HandlerFunc 类型的函数，
//	然后在 router/router.go 中注册使用
//
// ============================================================
package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware 跨域中间件
// 允许前端（不同端口）访问后端API
// 生产环境应限制 AllowOrigins 为具体域名
func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		// AllowOrigins: 允许的前端域名列表
		// 开发环境用 * 允许所有，生产环境改为具体域名如 ["https://openpanda.com"]
		AllowOrigins: []string{"*"},
		// AllowMethods: 允许的HTTP方法
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		// AllowHeaders: 允许的请求头
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization", // JWT Token 放在此请求头中
			"Accept-Language",
		},
		// ExposeHeaders: 允许前端读取的响应头
		ExposeHeaders: []string{"Content-Length"},
		// AllowCredentials: 是否允许携带Cookie
		AllowCredentials: true,
		// MaxAge: 预检请求缓存时间
		MaxAge: 12 * time.Hour,
	})
}
