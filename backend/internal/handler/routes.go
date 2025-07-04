package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	// 用户认证相关路由
	auth := r.Group("/auth")
	{
		auth.POST("/login", login)
		auth.POST("/register", register)
		auth.POST("/logout", logout)
		auth.POST("/refresh", refresh)
	}

	// 用户相关路由
	user := r.Group("/user")
	{
		user.GET("/profile", getUserProfile)
		user.PUT("/profile", updateUserProfile)
		user.GET("/servers", getUserServers)
	}

	// 服务器相关路由
	server := r.Group("/server")
	{
		server.GET("/products", getServerProducts)
		server.POST("/purchase", purchaseServer)
		server.GET("/:id", getServerDetail)
		server.POST("/:id/start", startServer)
		server.POST("/:id/stop", stopServer)
		server.POST("/:id/restart", restartServer)
	}

	// 管理员相关路由
	admin := r.Group("/admin")
	{
		admin.GET("/dashboard", getDashboard)
		admin.GET("/users", getUsers)
		admin.GET("/orders", getOrders)
		admin.GET("/products", getProducts)
	}
}

// 临时处理函数，后续会移到具体的handler文件中
func login(c *gin.Context) {
	c.JSON(200, gin.H{"message": "登录接口"})
}

func register(c *gin.Context) {
	c.JSON(200, gin.H{"message": "注册接口"})
}

func logout(c *gin.Context) {
	c.JSON(200, gin.H{"message": "登出接口"})
}

func refresh(c *gin.Context) {
	c.JSON(200, gin.H{"message": "刷新Token接口"})
}

func getUserProfile(c *gin.Context) {
	c.JSON(200, gin.H{"message": "获取用户资料接口"})
}

func updateUserProfile(c *gin.Context) {
	c.JSON(200, gin.H{"message": "更新用户资料接口"})
}

func getUserServers(c *gin.Context) {
	c.JSON(200, gin.H{"message": "获取用户服务器列表接口"})
}

func getServerProducts(c *gin.Context) {
	c.JSON(200, gin.H{"message": "获取服务器产品列表接口"})
}

func purchaseServer(c *gin.Context) {
	c.JSON(200, gin.H{"message": "购买服务器接口"})
}

func getServerDetail(c *gin.Context) {
	c.JSON(200, gin.H{"message": "获取服务器详情接口"})
}

func startServer(c *gin.Context) {
	c.JSON(200, gin.H{"message": "启动服务器接口"})
}

func stopServer(c *gin.Context) {
	c.JSON(200, gin.H{"message": "停止服务器接口"})
}

func restartServer(c *gin.Context) {
	c.JSON(200, gin.H{"message": "重启服务器接口"})
}

func getDashboard(c *gin.Context) {
	c.JSON(200, gin.H{"message": "管理员仪表板接口"})
}

func getUsers(c *gin.Context) {
	c.JSON(200, gin.H{"message": "获取用户列表接口"})
}

func getOrders(c *gin.Context) {
	c.JSON(200, gin.H{"message": "获取订单列表接口"})
}

func getProducts(c *gin.Context) {
	c.JSON(200, gin.H{"message": "获取产品列表接口"})
}