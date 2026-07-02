// ============================================================
// 包名: router
// 说明: 路由注册中心
//
//	所有API路由在此文件中集中注册和管理
//	拓展方式：
//	  1. 新增 controller 后，在此文件中注入依赖并注册路由
//	  2. 复杂项目可按模块拆分路由文件（如 router_article.go）
//	  3. 需要认证的路由组使用 .Use(middleware.JWTAuthMiddleware())
//
// ============================================================
package router

import (
	"openpanda-backend/controller"
	"openpanda-backend/middleware"
	"openpanda-backend/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter 初始化路由
// 参数 db: GORM 数据库连接实例
// 返回: 配置好的 Gin Engine
//
// 路由结构说明：
//
//	/api/v1/
//	  ├── /articles          (公开) 文章列表
//	  ├── /articles/:id      (公开) 文章详情
//	  ├── /articles/hot      (公开) 热门文章
//	  ├── /articles/search   (公开) 文章搜索
//	  ├── /categories        (公开) 分类列表
//	  └── /admin/articles    (需认证) 文章管理
func SetupRouter(db *gorm.DB, aiController *controller.AIController, settingController *controller.SettingController) *gin.Engine {
	// 创建 Gin 引擎（默认带 Logger 和 Recovery 中间件）
	r := gin.Default()

	// ============================================================
	// 全局中间件注册
	// ============================================================
	r.Use(middleware.CORSMiddleware()) // 跨域处理

	// ============================================================
	// 依赖注入：创建 Service 和 Controller 实例
	// ============================================================
	articleService := service.NewArticleService(db)
	categoryService := service.NewCategoryService(db)
	articleController := controller.NewArticleController(articleService, categoryService)
	authController := controller.NewAuthController()
	uploadController := controller.NewUploadController()
	sitemapController := controller.NewSitemapController(articleService, categoryService)

	// ============================================================
	// 静态文件服务（上传的图片通过此路由访问）
	// ============================================================
	r.Static("/uploads", "./uploads")

	// ============================================================
	// API v1 路由组（公开接口，无需认证）
	// ============================================================
	v1 := r.Group("/api/v1")
	{
		// --- 认证相关（公开） ---
		v1.POST("/login", authController.Login) // 登录

		// --- 文章相关 ---
		v1.GET("/articles", articleController.GetArticleList)        // 文章列表
		v1.GET("/articles/hot", articleController.GetHotArticles)    // 热门文章（注意：/hot 必须在 /:id 之前注册，否则会被 :id 匹配）
		v1.GET("/articles/search", articleController.SearchArticles) // 文章搜索（同理，在 /:id 之前）
		v1.GET("/articles/:id", articleController.GetArticleDetail)  // 文章详情

		// --- 分类相关 ---
		v1.GET("/categories", articleController.GetCategories) // 分类列表

		// 后续拓展示例：
		// v1.GET("/tags", tagController.GetTags)               // 标签列表
		// v1.GET("/comments", commentController.GetComments)   // 评论列表
	}

	// ============================================================
	// Admin 路由组（需 JWT 认证）
	// ============================================================
	admin := r.Group("/api/v1/admin")
	admin.Use(middleware.JWTAuthMiddleware())
	{
		// --- 文章管理 ---
		admin.POST("/articles", articleController.CreateArticle)
		admin.PUT("/articles/:id", articleController.UpdateArticle)
		admin.DELETE("/articles/:id", articleController.DeleteArticle)

		// --- 图片上传 ---
		admin.POST("/upload/image", uploadController.UploadImage)

		// --- 分类管理 ---
		admin.POST("/categories", articleController.CreateCategory)
		admin.PUT("/categories/:id", articleController.UpdateCategory)
		admin.DELETE("/categories/:id", articleController.DeleteCategory)

		// --- AI 聊天 ---
		admin.POST("/ai/chat", aiController.Chat)           // 聊天（SSE流式）
		admin.GET("/ai/providers", aiController.ListProviders) // 可用供应商列表

		// --- AI 设置 ---
		admin.GET("/settings/ai", settingController.GetAISettings)    // 获取 AI 配置
		admin.PUT("/settings/ai", settingController.UpdateAISettings) // 更新 AI 配置
	}

	// ============================================================
	// 认证路由组（需 JWT 认证，用于验证 Token 有效性等）
	// ============================================================
	auth := r.Group("/api/v1/auth")
	auth.Use(middleware.JWTAuthMiddleware())
	{
		auth.GET("/me", authController.GetMe) // 获取当前用户信息
	}

	// ============================================================
	// Sitemap（SEO，搜索引擎爬虫）
	// ============================================================
	r.GET("/sitemap.xml", sitemapController.Generate)

	// ============================================================
	// 健康检查接口（用于监控和负载均衡探测）
	// ============================================================
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	return r
}
