package handler

import (
	"time"
	"cloudbp-backend/internal/service"
	"cloudbp-backend/internal/middleware"
	"cloudbp-backend/internal/config"
	"cloudbp-backend/pkg/auth"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	// 加载配置
	cfg, _ := config.Load()
	
	// 创建JWT管理器
	jwtManager := auth.NewJWTManager(auth.JWTConfig{
		SecretKey:                cfg.JWT.Secret,
		Issuer:                   "cloudbp-backend",
		ExpiresIn:                time.Duration(cfg.JWT.ExpireTime) * time.Second,  // 从配置文件读取
		RefreshTokenExpiresIn:    time.Hour * 24 * 7,                              // 刷新令牌7天有效
	})

	// 创建服务
	userService := service.NewUserService(db, rdb, jwtManager)

	// 创建处理器
	authHandler := NewAuthHandler(userService)
	serverHandler := NewServerHandler(db, rdb)
	adminHandler := NewAdminHandler(db, rdb)

	// 用户认证相关路由（不需要JWT验证）
	auth := r.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", middleware.AuthMiddleware(jwtManager), authHandler.Logout)
	}

	// 用户相关路由（需要JWT验证）
	user := r.Group("/user")
	user.Use(middleware.AuthMiddleware(jwtManager))
	{
		user.GET("/profile", authHandler.GetProfile)
		user.PUT("/profile", authHandler.UpdateProfile)
		user.POST("/change-password", authHandler.ChangePassword)
		user.GET("/servers", serverHandler.GetUserServers)
	}

	// 服务器相关路由（需要JWT验证）
	server := r.Group("/server")
	server.Use(middleware.AuthMiddleware(jwtManager))
	{
		server.GET("/products", serverHandler.GetServerProducts)
		server.POST("/purchase", serverHandler.PurchaseServer)
		server.GET("/:id", serverHandler.GetServerDetail)
		server.POST("/:id/start", serverHandler.StartServer)
		server.POST("/:id/stop", serverHandler.StopServer)
		server.POST("/:id/restart", serverHandler.RestartServer)
	}

	// 管理员相关路由（需要JWT验证和管理员权限）
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware(jwtManager), middleware.AdminMiddleware())
	{
		admin.GET("/dashboard", adminHandler.GetDashboard)
		admin.GET("/users", adminHandler.GetUsers)
		admin.GET("/orders", adminHandler.GetOrders)
		admin.GET("/products", adminHandler.GetProducts)
	}
}