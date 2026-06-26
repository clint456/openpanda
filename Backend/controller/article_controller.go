// ============================================================
// 包名: controller
// 说明: 控制器层，处理HTTP请求和响应
//
//	每个 controller 对应一组相关的API接口
//	拓展方式：
//	  1. 新建 controller 文件（如 comment_controller.go）
//	  2. 定义 Controller 结构体，持有 *service.XxxService
//	  3. 在 router 中注册路由并注入依赖
//
// ============================================================
package controller

import (
	"strconv"

	"openpanda-backend/model"
	"openpanda-backend/service"
	"openpanda-backend/utils"

	"github.com/gin-gonic/gin"
)

// ArticleController 文章控制器
// 持有 ArticleService 用于处理业务逻辑
type ArticleController struct {
	ArticleService  *service.ArticleService
	CategoryService *service.CategoryService
}

// NewArticleController 构造函数
func NewArticleController(articleService *service.ArticleService, categoryService *service.CategoryService) *ArticleController {
	return &ArticleController{
		ArticleService:  articleService,
		CategoryService: categoryService,
	}
}

// ============================================================
// 文章相关接口
// ============================================================

// GetArticleList 获取文章列表
// GET /api/v1/articles?page=1&page_size=10&category_id=1&tag_id=2
func (ctrl *ArticleController) GetArticleList(c *gin.Context) {
	// 从查询参数获取分页信息，提供默认值
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	categoryID, _ := strconv.ParseUint(c.DefaultQuery("category_id", "0"), 10, 64)
	tagID, _ := strconv.ParseUint(c.DefaultQuery("tag_id", "0"), 10, 64)

	// 参数校验
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}

	// 调用 service 层获取数据
	articles, total, err := ctrl.ArticleService.GetList(page, pageSize, uint(categoryID), uint(tagID))
	if err != nil {
		utils.InternalError(c, "获取文章列表失败")
		return
	}

	// 返回带分页的成功响应
	utils.SuccessWithPage(c, utils.PageData{
		List:     articles,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// GetArticleDetail 获取文章详情
// GET /api/v1/articles/:id
func (ctrl *ArticleController) GetArticleDetail(c *gin.Context) {
	// 从路径参数获取文章ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "文章ID格式错误")
		return
	}

	article, err := ctrl.ArticleService.GetByID(uint(id))
	if err != nil {
		utils.NotFound(c, "文章不存在")
		return
	}

	// 异步增加阅读量（不影响接口响应速度）
	go func() {
		_ = ctrl.ArticleService.IncrementViewCount(uint(id))
	}()

	utils.Success(c, article)
}

// CreateArticle 创建文章
// POST /api/v1/articles
// 请求体 JSON 示例：
//
//	{
//	  "title": "STM32时钟树详解",
//	  "content": "<p>详细内容...</p>",
//	  "category_id": 1,
//	  "tag_ids": [1, 2]
//	}
func (ctrl *ArticleController) CreateArticle(c *gin.Context) {
	var input struct {
		Title      string `json:"title" binding:"required"`       // binding:"required" 表示必填
		Content    string `json:"content" binding:"required"`     // 正文必填
		Summary    string `json:"summary"`                        // 摘要可选
		CoverImage string `json:"cover_image"`                    // 封面图可选
		CategoryID uint   `json:"category_id" binding:"required"` // 分类必填
		TagIDs     []uint `json:"tag_ids"`                        // 标签ID列表（可选）
		Language   string `json:"language"`                       // 语言: zh/en/both
	}

	// 绑定并校验 JSON 请求体
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 构建文章模型
	article := &model.Article{
		Title:       input.Title,
		Content:     input.Content,
		Summary:     input.Summary,
		CoverImage:  input.CoverImage,
		CategoryID:  input.CategoryID,
		Language:    input.Language,
		IsPublished: true,
	}

	// 处理标签关联（如果传了标签ID列表）
	if len(input.TagIDs) > 0 {
		// 创建 Tag 切片并附加到 article
		tags := make([]model.Tag, len(input.TagIDs))
		for i, id := range input.TagIDs {
			tags[i] = model.Tag{ID: id}
		}
		article.Tags = tags
	}

	if err := ctrl.ArticleService.Create(article); err != nil {
		utils.InternalError(c, "创建文章失败")
		return
	}

	utils.SuccessWithMessage(c, "创建成功", article)
}

// UpdateArticle 更新文章
// PUT /api/v1/articles/:id
func (ctrl *ArticleController) UpdateArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "文章ID格式错误")
		return
	}

	// 先检查文章是否存在
	existing, err := ctrl.ArticleService.GetByID(uint(id))
	if err != nil {
		utils.NotFound(c, "文章不存在")
		return
	}

	var input struct {
		Title      string `json:"title"`
		Content    string `json:"content"`
		Summary    string `json:"summary"`
		CoverImage string `json:"cover_image"`
		CategoryID uint   `json:"category_id"`
		Language   string `json:"language"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	// 只更新非零值字段
	if input.Title != "" {
		existing.Title = input.Title
	}
	if input.Content != "" {
		existing.Content = input.Content
	}
	if input.Summary != "" {
		existing.Summary = input.Summary
	}
	if input.CoverImage != "" {
		existing.CoverImage = input.CoverImage
	}
	if input.CategoryID > 0 {
		existing.CategoryID = input.CategoryID
	}
	if input.Language != "" {
		existing.Language = input.Language
	}

	if err := ctrl.ArticleService.Update(existing); err != nil {
		utils.InternalError(c, "更新文章失败")
		return
	}

	utils.Success(c, existing)
}

// DeleteArticle 删除文章
// DELETE /api/v1/articles/:id
func (ctrl *ArticleController) DeleteArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "文章ID格式错误")
		return
	}

	if err := ctrl.ArticleService.Delete(uint(id)); err != nil {
		utils.InternalError(c, "删除文章失败")
		return
	}

	utils.Success(c, nil)
}

// GetHotArticles 获取热门文章
// GET /api/v1/articles/hot
func (ctrl *ArticleController) GetHotArticles(c *gin.Context) {
	articles, err := ctrl.ArticleService.GetHotArticles(10)
	if err != nil {
		utils.InternalError(c, "获取热门文章失败")
		return
	}
	utils.Success(c, articles)
}

// SearchArticles 搜索文章
// GET /api/v1/articles/search?keyword=STM32&page=1&page_size=10
func (ctrl *ArticleController) SearchArticles(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		utils.BadRequest(c, "请输入搜索关键词")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	articles, total, err := ctrl.ArticleService.Search(keyword, page, pageSize)
	if err != nil {
		utils.InternalError(c, "搜索失败")
		return
	}

	utils.SuccessWithPage(c, utils.PageData{
		List:     articles,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// ============================================================
// 分类相关接口
// ============================================================

// GetCategories 获取所有分类
// GET /api/v1/categories
func (ctrl *ArticleController) GetCategories(c *gin.Context) {
	categories, err := ctrl.CategoryService.GetAll()
	if err != nil {
		utils.InternalError(c, "获取分类失败")
		return
	}
	utils.Success(c, categories)
}
