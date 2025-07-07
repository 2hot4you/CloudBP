package service

import (
	"context"
	"fmt"

	"cloudbp-backend/internal/model"
	"gorm.io/gorm"
	"github.com/redis/go-redis/v9"
)

// AdminService 管理员服务
type AdminService struct {
	db  *gorm.DB
	rdb *redis.Client
}

// NewAdminService 创建管理员服务
func NewAdminService(db *gorm.DB, rdb *redis.Client) *AdminService {
	return &AdminService{
		db:  db,
		rdb: rdb,
	}
}

// DashboardData 仪表板数据
type DashboardData struct {
	UserCount       int64   `json:"user_count"`       // 用户总数
	ServerCount     int64   `json:"server_count"`     // 服务器总数
	OrderCount      int64   `json:"order_count"`      // 订单总数
	TodayOrderCount int64   `json:"today_order_count"`// 今日订单数
	TotalRevenue    float64 `json:"total_revenue"`    // 总收入
	TodayRevenue    float64 `json:"today_revenue"`    // 今日收入
	ActiveServers   int64   `json:"active_servers"`   // 运行中服务器数
	PendingOrders   int64   `json:"pending_orders"`   // 待处理订单数
}

// GetDashboard 获取仪表板数据
func (s *AdminService) GetDashboard(ctx context.Context) (*DashboardData, error) {
	data := &DashboardData{}

	// 获取用户总数
	if err := s.db.Model(&model.User{}).Count(&data.UserCount).Error; err != nil {
		return nil, fmt.Errorf("获取用户总数失败: %w", err)
	}

	// 获取服务器总数
	if err := s.db.Model(&model.Server{}).Count(&data.ServerCount).Error; err != nil {
		return nil, fmt.Errorf("获取服务器总数失败: %w", err)
	}

	// 获取订单总数
	if err := s.db.Model(&model.Order{}).Count(&data.OrderCount).Error; err != nil {
		return nil, fmt.Errorf("获取订单总数失败: %w", err)
	}

	// 获取今日订单数
	if err := s.db.Model(&model.Order{}).
		Where("DATE(created_at) = CURRENT_DATE").
		Count(&data.TodayOrderCount).Error; err != nil {
		return nil, fmt.Errorf("获取今日订单数失败: %w", err)
	}

	// 获取总收入
	var totalRevenue float64
	if err := s.db.Model(&model.Order{}).
		Where("status = ?", model.OrderStatusSuccess).
		Select("SUM(pay_amount)").
		Scan(&totalRevenue).Error; err != nil {
		return nil, fmt.Errorf("获取总收入失败: %w", err)
	}
	data.TotalRevenue = totalRevenue

	// 获取今日收入
	var todayRevenue float64
	if err := s.db.Model(&model.Order{}).
		Where("status = ? AND DATE(created_at) = CURRENT_DATE", model.OrderStatusSuccess).
		Select("SUM(pay_amount)").
		Scan(&todayRevenue).Error; err != nil {
		return nil, fmt.Errorf("获取今日收入失败: %w", err)
	}
	data.TodayRevenue = todayRevenue

	// 获取运行中服务器数
	if err := s.db.Model(&model.Server{}).
		Where("status = ?", model.ServerStatusRunning).
		Count(&data.ActiveServers).Error; err != nil {
		return nil, fmt.Errorf("获取运行中服务器数失败: %w", err)
	}

	// 获取待处理订单数
	if err := s.db.Model(&model.Order{}).
		Where("status IN ?", []string{model.OrderStatusPending, model.OrderStatusProcessing}).
		Count(&data.PendingOrders).Error; err != nil {
		return nil, fmt.Errorf("获取待处理订单数失败: %w", err)
	}

	return data, nil
}

// GetUsersRequest 获取用户列表请求
type GetUsersRequest struct {
	Page    int    `json:"page"`
	Size    int    `json:"size"`
	Keyword string `json:"keyword"`
}

// GetUsersResponse 获取用户列表响应
type GetUsersResponse struct {
	Users      []model.User `json:"users"`
	TotalCount int64        `json:"total_count"`
	Page       int          `json:"page"`
	Size       int          `json:"size"`
}

// GetUsers 获取用户列表
func (s *AdminService) GetUsers(ctx context.Context, req *GetUsersRequest) (*GetUsersResponse, error) {
	var users []model.User
	var totalCount int64

	// 构建查询
	query := s.db.Model(&model.User{})
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("username LIKE ? OR email LIKE ? OR real_name LIKE ?", keyword, keyword, keyword)
	}

	// 获取总数
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, fmt.Errorf("获取用户总数失败: %w", err)
	}

	// 获取分页数据
	offset := (req.Page - 1) * req.Size
	if err := query.Offset(offset).Limit(req.Size).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, fmt.Errorf("获取用户列表失败: %w", err)
	}

	// 脱敏处理
	for i := range users {
		users[i].Password = ""
	}

	return &GetUsersResponse{
		Users:      users,
		TotalCount: totalCount,
		Page:       req.Page,
		Size:       req.Size,
	}, nil
}

// GetOrdersRequest 获取订单列表请求
type GetOrdersRequest struct {
	Page    int    `json:"page"`
	Size    int    `json:"size"`
	Status  string `json:"status"`
	Keyword string `json:"keyword"`
}

// GetOrdersResponse 获取订单列表响应
type GetOrdersResponse struct {
	Orders     []model.Order `json:"orders"`
	TotalCount int64         `json:"total_count"`
	Page       int           `json:"page"`
	Size       int           `json:"size"`
}

// GetOrders 获取订单列表
func (s *AdminService) GetOrders(ctx context.Context, req *GetOrdersRequest) (*GetOrdersResponse, error) {
	var orders []model.Order
	var totalCount int64

	// 构建查询
	query := s.db.Model(&model.Order{}).Preload("User").Preload("Provider").Preload("Product")
	
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("order_no LIKE ?", keyword)
	}

	// 获取总数
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, fmt.Errorf("获取订单总数失败: %w", err)
	}

	// 获取分页数据
	offset := (req.Page - 1) * req.Size
	if err := query.Offset(offset).Limit(req.Size).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("获取订单列表失败: %w", err)
	}

	return &GetOrdersResponse{
		Orders:     orders,
		TotalCount: totalCount,
		Page:       req.Page,
		Size:       req.Size,
	}, nil
}

// GetAdminProductsRequest 获取产品列表请求
type GetAdminProductsRequest struct {
	Page       int  `json:"page"`
	Size       int  `json:"size"`
	ProviderID uint `json:"provider_id"`
}

// GetAdminProductsResponse 获取产品列表响应
type GetAdminProductsResponse struct {
	Products   []model.Product  `json:"products"`
	Providers  []model.Provider `json:"providers"`
	TotalCount int64            `json:"total_count"`
	Page       int              `json:"page"`
	Size       int              `json:"size"`
}

// GetProducts 获取产品列表
func (s *AdminService) GetProducts(ctx context.Context, req *GetAdminProductsRequest) (*GetAdminProductsResponse, error) {
	var products []model.Product
	var providers []model.Provider
	var totalCount int64

	// 获取所有云厂商
	if err := s.db.Find(&providers).Error; err != nil {
		return nil, fmt.Errorf("获取云厂商列表失败: %w", err)
	}

	// 构建查询
	query := s.db.Model(&model.Product{}).Preload("Provider")
	
	if req.ProviderID > 0 {
		query = query.Where("provider_id = ?", req.ProviderID)
	}

	// 获取总数
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, fmt.Errorf("获取产品总数失败: %w", err)
	}

	// 获取分页数据
	offset := (req.Page - 1) * req.Size
	if err := query.Offset(offset).Limit(req.Size).Order("provider_id, created_at DESC").Find(&products).Error; err != nil {
		return nil, fmt.Errorf("获取产品列表失败: %w", err)
	}

	return &GetAdminProductsResponse{
		Products:   products,
		Providers:  providers,
		TotalCount: totalCount,
		Page:       req.Page,
		Size:       req.Size,
	}, nil
}