// ============================================================
// 包名: controller
// 说明: Sitemap 控制器，动态生成 XML Sitemap
//       自动包含：首页、文章列表、分类页、所有已发布文章详情页
//       生产环境通过 SITE_URL 环境变量设置域名
// ============================================================
package controller

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"openpanda-backend/config"
	"openpanda-backend/service"

	"github.com/gin-gonic/gin"
)

// SitemapController Sitemap 控制器
type SitemapController struct {
	articleSvc  *service.ArticleService
	categorySvc *service.CategoryService
	siteURL     string
}

// NewSitemapController 构造函数
func NewSitemapController(articleSvc *service.ArticleService, categorySvc *service.CategoryService) *SitemapController {
	siteURL := config.GetEnv("SITE_URL", "http://localhost:3000")
	siteURL = strings.TrimRight(siteURL, "/")
	return &SitemapController{
		articleSvc:  articleSvc,
		categorySvc: categorySvc,
		siteURL:     siteURL,
	}
}

// sitemapEntry sitemap 条目
type sitemapEntry struct {
	Loc        string
	LastMod    string
	ChangeFreq string
	Priority   string
}

// Generate 生成并返回 XML Sitemap
// GET /sitemap.xml
func (ctrl *SitemapController) Generate(c *gin.Context) {
	now := time.Now().Format(time.RFC3339)
	entries := []sitemapEntry{
		// 静态页面
		{Loc: ctrl.siteURL + "/", ChangeFreq: "daily", Priority: "1.0", LastMod: now},
		{Loc: ctrl.siteURL + "/articles", ChangeFreq: "daily", Priority: "0.9", LastMod: now},
	}

	// 分类页面
	categories, err := ctrl.categorySvc.GetAll()
	if err == nil {
		for _, cat := range categories {
			entries = append(entries, sitemapEntry{
				Loc:        fmt.Sprintf("%s/category/%s", ctrl.siteURL, cat.Slug),
				ChangeFreq: "weekly",
				Priority:   "0.7",
			})
		}
	}

	// 所有已发布文章
	articles, err := ctrl.articleSvc.GetAllPublished()
	if err == nil {
		for _, a := range articles {
			slug := a.Slug
			if slug == "" {
				slug = fmt.Sprintf("%d", a.ID)
			}
			entries = append(entries, sitemapEntry{
				Loc:        fmt.Sprintf("%s/articles/%d-%s", ctrl.siteURL, a.ID, slug),
				LastMod:    a.UpdatedAt.Format(time.RFC3339),
				ChangeFreq: "monthly",
				Priority:   "0.8",
			})
		}
	}

	// 构建 XML
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	sb.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")

	for _, e := range entries {
		sb.WriteString("  <url>\n")
		sb.WriteString(fmt.Sprintf("    <loc>%s</loc>\n", xmlEscape(e.Loc)))
		if e.LastMod != "" {
			sb.WriteString(fmt.Sprintf("    <lastmod>%s</lastmod>\n", xmlEscape(e.LastMod)))
		}
		if e.ChangeFreq != "" {
			sb.WriteString(fmt.Sprintf("    <changefreq>%s</changefreq>\n", xmlEscape(e.ChangeFreq)))
		}
		if e.Priority != "" {
			sb.WriteString(fmt.Sprintf("    <priority>%s</priority>\n", xmlEscape(e.Priority)))
		}
		sb.WriteString("  </url>\n")
	}

	sb.WriteString("</urlset>\n")

	c.Header("Content-Type", "application/xml; charset=utf-8")
	c.String(http.StatusOK, sb.String())
}

// xmlEscape 转义 XML 特殊字符
func xmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}
