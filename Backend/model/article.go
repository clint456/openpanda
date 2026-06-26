// ============================================================
// 包名: model
// 说明: 数据库模型定义（GORM），一张表对应一个结构体文件
//
//	拓展方式：
//	  1. 仿照此文件新建 model 文件（如 comment.go, user.go）
//	  2. 定义结构体并实现 TableName() 方法
//	  3. 在 main.go 的 AutoMigrate 中注册新模型
//
// ============================================================
package model

import (
	"time"
)

// Article 文章模型（核心业务表）
// GORM 标签说明：
//
//	column: 数据库列名
//	type: 数据库列类型
//	not null: 非空约束
//	default: 默认值
//	index: 创建索引（加速查询）
type Article struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`            // 主键ID，自增
	Title       string    `gorm:"type:varchar(255);not null;index" json:"title"` // 文章标题，建索引方便搜索
	Slug        string    `gorm:"type:varchar(255);uniqueIndex" json:"slug"`     // URL友好标识，唯一索引
	Content     string    `gorm:"type:text;not null" json:"content"`             // 文章正文（支持富文本HTML）
	Summary     string    `gorm:"type:varchar(500)" json:"summary"`              // 文章摘要
	CoverImage  string    `gorm:"type:varchar(500)" json:"cover_image"`          // 封面图URL
	CategoryID  uint      `gorm:"index" json:"category_id"`                      // 所属分类ID，建索引
	ViewCount   int       `gorm:"default:0" json:"view_count"`                   // 阅读量
	IsPublished bool      `gorm:"default:false;index" json:"is_published"`       // 是否已发布
	Language    string    `gorm:"type:varchar(10);default:'zh'" json:"language"` // 语言: zh=中文, en=英文, both=双语
	CreatedAt   time.Time `json:"created_at"`                                    // 创建时间（GORM自动管理）
	UpdatedAt   time.Time `json:"updated_at"`                                    // 更新时间（GORM自动管理）

	// 关联关系（GORM 外键）
	// 拓展方式：新增关联只需添加字段，例如：
	//   Comments []Comment `gorm:"foreignKey:ArticleID"` // 一对多关联评论
	Category Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"` // 所属分类
	Tags     []Tag    `gorm:"many2many:article_tags;" json:"tags,omitempty"`   // 多对多关联标签
}

// TableName 指定表名（GORM 默认用结构体名复数形式，这里显式指定）
func (Article) TableName() string {
	return "articles"
}

// Category 文章分类模型
type Category struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"` // 分类名称，唯一
	Slug        string    `gorm:"type:varchar(100);uniqueIndex" json:"slug"`          // URL友好标识
	Description string    `gorm:"type:varchar(500)" json:"description"`               // 分类描述
	SortOrder   int       `gorm:"default:0" json:"sort_order"`                        // 排序序号
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Category) TableName() string {
	return "categories"
}

// Tag 标签模型
type Tag struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"` // 标签名，唯一
	Slug      string    `gorm:"type:varchar(100);uniqueIndex" json:"slug"`          // URL友好标识
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Tag) TableName() string {
	return "tags"
}
