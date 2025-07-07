package handler

import (
	"net/http"
	"cloudbp-backend/internal/service"
	"cloudbp-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 认证处理器
type AuthHandler struct {
	userService *service.UserService
}

// 创建认证处理器
func NewAuthHandler(userService *service.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

// 用户登录
// @Summary 用户登录
// @Description 用户登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param body body service.LoginRequest true "登录请求"
// @Success 200 {object} service.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("绑定登录请求失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	resp, err := h.userService.Login(&req)
	if err != nil {
		logger.Log.Error("用户登录失败", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	logger.Log.Info("用户登录成功", zap.String("username", req.Username))
	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"data":    resp,
	})
}

// 用户注册
// @Summary 用户注册
// @Description 用户注册
// @Tags 认证
// @Accept json
// @Produce json
// @Param body body service.RegisterRequest true "注册请求"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("绑定注册请求失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	if err := h.userService.Register(&req); err != nil {
		logger.Log.Error("用户注册失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Log.Info("用户注册成功", zap.String("username", req.Username))
	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
	})
}

// 刷新令牌
// @Summary 刷新令牌
// @Description 刷新访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param body body map[string]string true "刷新令牌请求"
// @Success 200 {object} service.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("绑定刷新令牌请求失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	resp, err := h.userService.RefreshToken(req.RefreshToken)
	if err != nil {
		logger.Log.Error("刷新令牌失败", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	logger.Log.Info("刷新令牌成功")
	c.JSON(http.StatusOK, gin.H{
		"message": "刷新成功",
		"data":    resp,
	})
}

// 用户登出
// @Summary 用户登出
// @Description 用户登出
// @Tags 认证
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]string
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// 获取用户信息
	username, _ := c.Get("username")
	
	// 这里可以添加将令牌加入黑名单的逻辑
	// 目前仅返回成功消息
	
	logger.Log.Info("用户登出", zap.Any("username", username))
	c.JSON(http.StatusOK, gin.H{
		"message": "登出成功",
	})
}

// 获取用户信息
// @Summary 获取用户信息
// @Description 获取当前用户的详细信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} model.User
// @Failure 401 {object} map[string]string
// @Router /user/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未找到用户信息"})
		return
	}

	user, err := h.userService.GetProfile(userID.(uint))
	if err != nil {
		logger.Log.Error("获取用户信息失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"data":    user,
	})
}

// 更新用户信息
// @Summary 更新用户信息
// @Description 更新当前用户的信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body service.UpdateProfileRequest true "更新用户信息请求"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /user/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未找到用户信息"})
		return
	}

	var req service.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("绑定更新用户信息请求失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	if err := h.userService.UpdateProfile(userID.(uint), &req); err != nil {
		logger.Log.Error("更新用户信息失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Log.Info("更新用户信息成功", zap.Any("user_id", userID))
	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
	})
}

// 修改密码
// @Summary 修改密码
// @Description 修改当前用户的密码
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body service.ChangePasswordRequest true "修改密码请求"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /user/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未找到用户信息"})
		return
	}

	var req service.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("绑定修改密码请求失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	if err := h.userService.ChangePassword(userID.(uint), &req); err != nil {
		logger.Log.Error("修改密码失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Log.Info("修改密码成功", zap.Any("user_id", userID))
	c.JSON(http.StatusOK, gin.H{
		"message": "修改密码成功",
	})
}