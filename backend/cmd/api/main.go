package main

import (
	"log"
	"cloudbp-backend/internal/config"
	"cloudbp-backend/internal/handler"
	"cloudbp-backend/internal/middleware"
	"cloudbp-backend/pkg/database"
	"cloudbp-backend/pkg/cache"
	"cloudbp-backend/pkg/logger"
	
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title 云服务器销售平台API
// @version 1.0
// @description 多厂商云服务器销售平台后端API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// 初始化配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("加载配置失败:", err)
	}

	// 初始化日志
	logger.Init(cfg.Log.Level)

	// 初始化数据库
	db, err := database.Init(cfg.Database)
	if err != nil {
		log.Fatal("数据库初始化失败:", err)
	}

	// 初始化缓存
	rdb, err := cache.Init(cfg.Redis)
	if err != nil {
		log.Fatal("Redis初始化失败:", err)
	}

	// 设置Gin模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建路由
	r := gin.New()
	
	// 全局中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Cors())
	r.Use(gin.Recovery())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API路由组
	apiV1 := r.Group("/api/v1")
	{
		// 注册路由
		handler.RegisterRoutes(apiV1, db, rdb)
	}

	// 启动服务器
	log.Printf("服务器启动在端口: %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}