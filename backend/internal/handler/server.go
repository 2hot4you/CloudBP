package handler

import (
	"net/http"
	"strconv"

	"cloudbp-backend/internal/service"
	"cloudbp-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ServerHandler 服务器处理器
type ServerHandler struct {
	db              *gorm.DB
	rdb             *redis.Client
	serverService   *service.ServerService
	providerService *service.ProviderService
}

// NewServerHandler 创建服务器处理器
func NewServerHandler(db *gorm.DB, rdb *redis.Client) *ServerHandler {
	serverService := service.NewServerService(db, rdb)
	providerService := service.NewProviderService(db)
	
	// 初始化云厂商
	if err := providerService.InitProviders(); err != nil {
		logger.Log.Error("初始化云厂商失败", zap.Error(err))
	}

	return &ServerHandler{
		db:              db,
		rdb:             rdb,
		serverService:   serverService,
		providerService: providerService,
	}
}

// GetUserServers 获取用户服务器列表
// @Summary 获取用户服务器列表
// @Description 获取当前用户的服务器列表
// @Tags 服务器
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /user/servers [get]
func (h *ServerHandler) GetUserServers(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未找到用户信息"})
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	req := &service.GetUserServersRequest{
		UserID: userID.(uint),
		Page:   page,
		Size:   size,
	}

	resp, err := h.serverService.GetUserServers(c.Request.Context(), req)
	if err != nil {
		logger.Log.Error("获取用户服务器列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"data":    resp,
	})
}

// GetServerProducts 获取服务器产品列表
// @Summary 获取服务器产品列表
// @Description 获取可购买的服务器产品列表
// @Tags 服务器
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param provider_id query int false "厂商ID"
// @Param region query string false "地域"
// @Param type query string false "产品类型"
// @Success 200 {object} map[string]interface{}
// @Router /server/products [get]
func (h *ServerHandler) GetServerProducts(c *gin.Context) {
	providerID, _ := strconv.Atoi(c.Query("provider_id"))
	region := c.Query("region")
	productType := c.Query("type")

	req := &service.GetProductsRequest{
		ProviderID: uint(providerID),
		Region:     region,
		Type:       productType,
	}

	resp, err := h.serverService.GetProducts(c.Request.Context(), req)
	if err != nil {
		logger.Log.Error("获取产品列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"data":    resp,
	})
}

// PurchaseServer 购买服务器
// @Summary 购买服务器
// @Description 创建服务器购买订单
// @Tags 服务器
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body service.PurchaseServerRequest true "购买请求"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /server/purchase [post]
func (h *ServerHandler) PurchaseServer(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未找到用户信息"})
		return
	}

	var req service.PurchaseServerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("绑定购买请求失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	req.UserID = userID.(uint)

	resp, err := h.serverService.PurchaseServer(c.Request.Context(), &req)
	if err != nil {
		logger.Log.Error("购买服务器失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "购买成功",
		"data":    resp,
	})
}

// GetServerDetail 获取服务器详情
// @Summary 获取服务器详情
// @Description 获取指定服务器的详细信息
// @Tags 服务器
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "服务器ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /server/{id} [get]
func (h *ServerHandler) GetServerDetail(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未找到用户信息"})
		return
	}

	serverID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的服务器ID"})
		return
	}

	req := &service.GetServerDetailRequest{
		ServerID: uint(serverID),
		UserID:   userID.(uint),
	}

	resp, err := h.serverService.GetServerDetail(c.Request.Context(), req)
	if err != nil {
		logger.Log.Error("获取服务器详情失败", zap.Error(err))
		if err.Error() == "服务器不存在" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"data":    resp,
	})
}

// StartServer 启动服务器
// @Summary 启动服务器
// @Description 启动指定的服务器
// @Tags 服务器
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "服务器ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /server/{id}/start [post]
func (h *ServerHandler) StartServer(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未找到用户信息"})
		return
	}

	serverID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的服务器ID"})
		return
	}

	req := &service.StartServerRequest{
		ServerID: uint(serverID),
		UserID:   userID.(uint),
	}

	if err := h.providerService.StartInstance(c.Request.Context(), &service.StartInstanceRequest{
		ServerID: req.ServerID,
	}); err != nil {
		logger.Log.Error("启动服务器失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "启动命令已发送",
	})
}

// StopServer 停止服务器
// @Summary 停止服务器
// @Description 停止指定的服务器
// @Tags 服务器
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "服务器ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /server/{id}/stop [post]
func (h *ServerHandler) StopServer(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未找到用户信息"})
		return
	}

	serverID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的服务器ID"})
		return
	}

	req := &service.StopServerRequest{
		ServerID: uint(serverID),
		UserID:   userID.(uint),
	}

	if err := h.providerService.StopInstance(c.Request.Context(), &service.StopInstanceRequest{
		ServerID: req.ServerID,
	}); err != nil {
		logger.Log.Error("停止服务器失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "停止命令已发送",
	})
}

// RestartServer 重启服务器
// @Summary 重启服务器
// @Description 重启指定的服务器
// @Tags 服务器
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "服务器ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /server/{id}/restart [post]
func (h *ServerHandler) RestartServer(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未找到用户信息"})
		return
	}

	serverID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的服务器ID"})
		return
	}

	req := &service.RestartServerRequest{
		ServerID: uint(serverID),
		UserID:   userID.(uint),
	}

	if err := h.providerService.RestartInstance(c.Request.Context(), &service.RestartInstanceRequest{
		ServerID: req.ServerID,
	}); err != nil {
		logger.Log.Error("重启服务器失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "重启命令已发送",
	})
}