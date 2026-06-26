// ============================================================
// 包名: main
// 说明: 项目入口文件
//
//	负责：加载配置 → 连接数据库/Redis → 自动迁移表结构 → 注册路由 → 启动服务
//
// 运行方式：
//
//	go run main.go
//
// 拓展方式：
//  1. 新增 model 后，在 AutoMigrate 中添加
//  2. 新增 service/controller 后，在 router.SetupRouter 中注册
//
// ============================================================
package main

import (
	"fmt"
	"log"

	"openpanda-backend/config"
	"openpanda-backend/model"
	"openpanda-backend/router"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// ============================================================
	// 1. 加载配置
	// ============================================================
	cfg := config.LoadConfig()

	// ============================================================
	// 2. 连接 PostgreSQL 数据库
	// ============================================================
	// 构建 DSN (数据源名称)
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
		cfg.Database.SSLMode,
		cfg.Database.TimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 打印SQL日志（生产环境可改为 logger.Error）
	})
	if err != nil {
		log.Fatalf("❌ 数据库连接失败: %v", err)
	}
	log.Println("✅ 数据库连接成功")

	// ============================================================
	// 3. 自动迁移表结构（AutoMigrate）
	//    GORM 会根据 model 结构体自动创建/更新表结构
	//    新增 model 后，在此处添加：
	//       db.AutoMigrate(&model.NewModel{})
	// ============================================================
	err = db.AutoMigrate(
		&model.Article{},
		&model.Category{},
		&model.Tag{},
		// 后续拓展：在此添加新的 model
		// &model.User{},
		// &model.Comment{},
	)
	if err != nil {
		log.Fatalf("❌ 数据库迁移失败: %v", err)
	}
	log.Println("✅ 数据库表结构已同步")

	// ============================================================
	// 4. 初始化种子数据（可选）
	//    首次运行时插入默认分类数据
	// ============================================================
	seedDefaultCategories(db)

	// ============================================================
	// 5. 设置路由并启动服务
	// ============================================================
	r := router.SetupRouter(db)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("🚀 服务器启动在 http://localhost%s", addr)

	// 启动HTTP服务
	if err := r.Run(addr); err != nil {
		log.Fatalf("❌ 服务器启动失败: %v", err)
	}
}

// seedDefaultCategories 初始化默认分类数据
// 只在分类表为空时插入，避免重复
func seedDefaultCategories(db *gorm.DB) {
	var count int64
	db.Model(&model.Category{}).Count(&count)
	if count > 0 {
		return // 已有数据，跳过
	}

	categories := []model.Category{
		{Name: "嵌入式Linux", Slug: "embedded-linux", Description: "Linux环境搭建、驱动开发、系统移植等实操记录", SortOrder: 1},
		{Name: "硬件电路设计", Slug: "hardware-design", Description: "原理图设计、PCB布局、硬件调试技巧", SortOrder: 2},
		{Name: "单片机开发", Slug: "mcu-development", Description: "STM32、51、ESP32等单片机开发教程与项目实战", SortOrder: 3},
	}

	// 使用 Create 批量插入（如果已存在相同 slug 则跳过）
	for _, cat := range categories {
		db.FirstOrCreate(&cat, model.Category{Slug: cat.Slug})
	}
	log.Println("✅ 默认分类数据已初始化")
}
