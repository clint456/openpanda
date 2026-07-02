// ============================================================
// 包名: service
// 说明: 业务逻辑层，处理具体业务规则
//
//	controller 只负责接收请求和返回响应，具体逻辑放在 service 中
//	拓展方式：
//	  1. 新建 service 文件（如 comment_service.go）
//	  2. 定义 Service 结构体，包含 *gorm.DB 依赖
//	  3. 在 controller 中注入使用
//
// ============================================================
package service

import (
	"openpanda-backend/model"
	"openpanda-backend/utils"

	"gorm.io/gorm"
)

// ArticleService 文章业务服务
// 结构体中持有 *gorm.DB，通过依赖注入方式使用
type ArticleService struct {
	DB *gorm.DB
}

// NewArticleService 构造函数：创建 ArticleService 实例
func NewArticleService(db *gorm.DB) *ArticleService {
	return &ArticleService{DB: db}
}

// GetList 获取文章列表（支持分页、分类筛选、标签筛选）
// 参数说明：
//
//	page:     页码（从1开始）
//	pageSize: 每页条数
//	categoryID: 分类ID（0表示不过滤）
//	tagID:    标签ID（0表示不过滤）
//
// 返回值：文章列表、总数、错误
func (s *ArticleService) GetList(page, pageSize int, categoryID, tagID uint) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64

	// 构建查询条件
	query := s.DB.Model(&model.Article{}).Where("is_published = ?", true)

	// 按分类筛选
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	// 按标签筛选（多对多关联查询）
	if tagID > 0 {
		query = query.Joins("JOIN article_tags ON article_tags.article_id = articles.id").
			Where("article_tags.tag_id = ?", tagID)
	}

	// 查询总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询，预加载关联的 Category 和 Tags
	offset := (page - 1) * pageSize
	if err := query.
		Preload("Category"). // 预加载分类信息
		Preload("Tags").     // 预加载标签信息
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// GetByID 根据ID获取文章详情
func (s *ArticleService) GetByID(id uint) (*model.Article, error) {
	var article model.Article
	if err := s.DB.
		Preload("Category").
		Preload("Tags").
		First(&article, id).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

// Create 创建文章（自动生成 slug）
func (s *ArticleService) Create(article *model.Article) error {
	if article.Slug == "" {
		article.Slug = utils.GenerateSlug(article.Title)
	}
	return s.DB.Create(article).Error
}

// Update 更新文章（slug 为空则重新生成）
func (s *ArticleService) Update(article *model.Article) error {
	if article.Slug == "" {
		article.Slug = utils.GenerateSlug(article.Title)
	}
	return s.DB.Save(article).Error
}

// Delete 删除文章（软删除建议用 GORM 的 DeletedAt，这里做硬删除示例）
func (s *ArticleService) Delete(id uint) error {
	return s.DB.Delete(&model.Article{}, id).Error
}

// IncrementViewCount 增加文章阅读量
// 后续可改为 Redis 缓存计数 + 定时同步到数据库
func (s *ArticleService) IncrementViewCount(id uint) error {
	return s.DB.Model(&model.Article{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// GetHotArticles 获取热门文章（按阅读量排序）
func (s *ArticleService) GetHotArticles(limit int) ([]model.Article, error) {
	var articles []model.Article
	if err := s.DB.
		Where("is_published = ?", true).
		Order("view_count DESC").
		Limit(limit).
		Preload("Category").
		Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

// GetAllPublished 获取所有已发布文章（无分页，用于 sitemap 等）
func (s *ArticleService) GetAllPublished() ([]model.Article, error) {
	var articles []model.Article
	if err := s.DB.
		Where("is_published = ?", true).
		Order("updated_at DESC").
		Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

// Search 文章搜索（标题模糊匹配）
func (s *ArticleService) Search(keyword string, page, pageSize int) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64

	query := s.DB.Model(&model.Article{}).
		Where("is_published = ?", true).
		Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.
		Preload("Category").
		Preload("Tags").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// ============================================================
// CategoryService 分类业务服务
// ============================================================

// CategoryService 分类业务服务
type CategoryService struct {
	DB *gorm.DB
}

// NewCategoryService 构造函数
func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{DB: db}
}

// GetAll 获取所有分类
func (s *CategoryService) GetAll() ([]model.Category, error) {
	var categories []model.Category
	if err := s.DB.Order("sort_order ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetBySlug 根据 slug 获取分类
func (s *CategoryService) GetBySlug(slug string) (*model.Category, error) {
	var category model.Category
	if err := s.DB.Where("slug = ?", slug).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// Create 创建分类
func (s *CategoryService) Create(category *model.Category) error {
	return s.DB.Create(category).Error
}

// Update 更新分类
func (s *CategoryService) Update(category *model.Category) error {
	return s.DB.Save(category).Error
}

// Delete 删除分类
func (s *CategoryService) Delete(id uint) error {
	return s.DB.Delete(&model.Category{}, id).Error
}
