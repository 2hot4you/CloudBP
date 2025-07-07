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

// AdminHandler 管理员处理器
type AdminHandler struct {
	db           *gorm.DB
	rdb          *redis.Client
	adminService *service.AdminService
}

// NewAdminHandler 创建管理员处理器
func NewAdminHandler(db *gorm.DB, rdb *redis.Client) *AdminHandler {
	adminService := service.NewAdminService(db, rdb)

	return &AdminHandler{
		db:           db,
		rdb:          rdb,
		adminService: adminService,
	}
}

// GetDashboard 获取管理员仪表板
// @Summary 获取管理员仪表板
// @Description 获取管理员仪表板统计数据
// @Tags 管理员
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /admin/dashboard [get]
func (h *AdminHandler) GetDashboard(c *gin.Context) {
	resp, err := h.adminService.GetDashboard(c.Request.Context())
	if err != nil {
		logger.Log.Error("获取仪表板数据失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"data":    resp,
	})
}

// GetUsers 获取用户列表
// @Summary 获取用户列表
// @Description 获取用户列表（管理员）
// @Tags 管理员
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /admin/users [get]
func (h *AdminHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	keyword := c.Query("keyword")

	req := &service.GetUsersRequest{
		Page:    page,
		Size:    size,
		Keyword: keyword,
	}

	resp, err := h.adminService.GetUsers(c.Request.Context(), req)
	if err != nil {
		logger.Log.Error("获取用户列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"data":    resp,
	})
}

// GetOrders 获取订单列表
// @Summary 获取订单列表
// @Description 获取订单列表（管理员）
// @Tags 管理员
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param status query string false "订单状态"
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /admin/orders [get]
func (h *AdminHandler) GetOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	status := c.Query("status")
	keyword := c.Query("keyword")

	req := &service.GetOrdersRequest{
		Page:    page,
		Size:    size,
		Status:  status,
		Keyword: keyword,
	}

	resp, err := h.adminService.GetOrders(c.Request.Context(), req)
	if err != nil {
		logger.Log.Error("获取订单列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"data":    resp,
	})
}

// GetProducts 获取产品列表
// @Summary 获取产品列表
// @Description 获取产品列表（管理员）
// @Tags 管理员
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param provider_id query int false "厂商ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /admin/products [get]
func (h *AdminHandler) GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	providerID, _ := strconv.Atoi(c.Query("provider_id"))

	req := &service.GetAdminProductsRequest{
		Page:       page,
		Size:       size,
		ProviderID: uint(providerID),
	}

	resp, err := h.adminService.GetProducts(c.Request.Context(), req)
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