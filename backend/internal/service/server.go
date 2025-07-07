package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloudbp-backend/internal/model"
	"gorm.io/gorm"
	"github.com/redis/go-redis/v9"
)

// ServerService 服务器服务
type ServerService struct {
	db  *gorm.DB
	rdb *redis.Client
}

// NewServerService 创建服务器服务
func NewServerService(db *gorm.DB, rdb *redis.Client) *ServerService {
	return &ServerService{
		db:  db,
		rdb: rdb,
	}
}

// GetUserServersRequest 获取用户服务器列表请求
type GetUserServersRequest struct {
	UserID uint `json:"user_id"`
	Page   int  `json:"page"`
	Size   int  `json:"size"`
}

// GetUserServersResponse 获取用户服务器列表响应
type GetUserServersResponse struct {
	Servers    []model.Server `json:"servers"`
	TotalCount int64          `json:"total_count"`
	Page       int            `json:"page"`
	Size       int            `json:"size"`
}

// GetUserServers 获取用户服务器列表
func (s *ServerService) GetUserServers(ctx context.Context, req *GetUserServersRequest) (*GetUserServersResponse, error) {
	var servers []model.Server
	var totalCount int64

	// 获取总数
	if err := s.db.Model(&model.Server{}).Where("user_id = ?", req.UserID).Count(&totalCount).Error; err != nil {
		return nil, fmt.Errorf("获取服务器总数失败: %w", err)
	}

	// 获取分页数据
	offset := (req.Page - 1) * req.Size
	if err := s.db.Preload("Provider").Preload("Product").Preload("Order").
		Where("user_id = ?", req.UserID).
		Offset(offset).Limit(req.Size).
		Order("created_at DESC").
		Find(&servers).Error; err != nil {
		return nil, fmt.Errorf("获取服务器列表失败: %w", err)
	}

	return &GetUserServersResponse{
		Servers:    servers,
		TotalCount: totalCount,
		Page:       req.Page,
		Size:       req.Size,
	}, nil
}

// GetProductsRequest 获取产品列表请求
type GetProductsRequest struct {
	ProviderID uint   `json:"provider_id"`
	Region     string `json:"region"`
	Type       string `json:"type"`
}

// GetProductsResponse 获取产品列表响应
type GetProductsResponse struct {
	Products  []model.Product  `json:"products"`
	Providers []model.Provider `json:"providers"`
}

// GetProducts 获取产品列表
func (s *ServerService) GetProducts(ctx context.Context, req *GetProductsRequest) (*GetProductsResponse, error) {
	var products []model.Product
	var providers []model.Provider

	// 获取启用的云厂商
	if err := s.db.Where("status = ?", model.ProviderStatusActive).Find(&providers).Error; err != nil {
		return nil, fmt.Errorf("获取云厂商列表失败: %w", err)
	}

	// 构建查询条件
	query := s.db.Preload("Provider").Where("status = ?", model.ProductStatusOnline)
	
	if req.ProviderID > 0 {
		query = query.Where("provider_id = ?", req.ProviderID)
	}
	if req.Region != "" {
		query = query.Where("region = ?", req.Region)
	}
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	// 获取产品列表
	if err := query.Order("provider_id, price").Find(&products).Error; err != nil {
		return nil, fmt.Errorf("获取产品列表失败: %w", err)
	}

	return &GetProductsResponse{
		Products:  products,
		Providers: providers,
	}, nil
}

// PurchaseServerRequest 购买服务器请求
type PurchaseServerRequest struct {
	UserID    uint   `json:"user_id"`
	ProductID uint   `json:"product_id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Period    int    `json:"period" binding:"required"`
	Quantity  int    `json:"quantity"`
	ImageID   string `json:"image_id"`
	Password  string `json:"password"`
	AutoRenew bool   `json:"auto_renew"`
}

// PurchaseServerResponse 购买服务器响应
type PurchaseServerResponse struct {
	OrderID   uint   `json:"order_id"`
	OrderNo   string `json:"order_no"`
	PayAmount float64 `json:"pay_amount"`
}

// PurchaseServer 购买服务器
func (s *ServerService) PurchaseServer(ctx context.Context, req *PurchaseServerRequest) (*PurchaseServerResponse, error) {
	// 获取产品信息
	var product model.Product
	if err := s.db.Preload("Provider").First(&product, req.ProductID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("产品不存在")
		}
		return nil, fmt.Errorf("获取产品信息失败: %w", err)
	}

	// 检查产品状态
	if product.Status != model.ProductStatusOnline {
		return nil, errors.New("产品已下架")
	}

	// 检查云厂商状态
	if product.Provider.Status != model.ProviderStatusActive {
		return nil, errors.New("云厂商暂不可用")
	}

	// 获取用户信息
	var user model.User
	if err := s.db.First(&user, req.UserID).Error; err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 设置默认值
	if req.Quantity == 0 {
		req.Quantity = 1
	}

	// 计算价格
	amount := product.Price * float64(req.Period) * float64(req.Quantity)
	discountAmount := 0.0
	payAmount := amount - discountAmount

	// 检查余额
	if user.Balance < payAmount {
		return nil, errors.New("余额不足")
	}

	// 生成订单号
	orderNo := fmt.Sprintf("ORD%d%d", time.Now().Unix(), req.UserID)

	// 开始事务
	var result *PurchaseServerResponse
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 创建订单
		order := model.Order{
			UserID:         req.UserID,
			OrderNo:        orderNo,
			ProviderID:     product.ProviderID,
			ProductID:      req.ProductID,
			Type:           model.OrderTypeNew,
			Status:         model.OrderStatusPending,
			Amount:         amount,
			DiscountAmount: discountAmount,
			PayAmount:      payAmount,
			PayMethod:      model.PaymentMethodBalance,
			Period:         req.Period,
			Quantity:       req.Quantity,
			Config:         fmt.Sprintf(`{"name": "%s", "image_id": "%s", "password": "%s", "auto_renew": %t}`, req.Name, req.ImageID, req.Password, req.AutoRenew),
		}

		if err := tx.Create(&order).Error; err != nil {
			return fmt.Errorf("创建订单失败: %w", err)
		}

		// 扣除余额
		if err := tx.Model(&user).Update("balance", gorm.Expr("balance - ?", payAmount)).Error; err != nil {
			return fmt.Errorf("扣除余额失败: %w", err)
		}

		// 创建支付记录
		payment := model.Payment{
			OrderID:       order.ID,
			UserID:        req.UserID,
			PaymentNo:     fmt.Sprintf("PAY%d%d", time.Now().Unix(), req.UserID),
			Method:        model.PaymentMethodBalance,
			Amount:        payAmount,
			Status:        model.PaymentStatusSuccess,
			PayTime:       &[]time.Time{time.Now()}[0],
		}

		if err := tx.Create(&payment).Error; err != nil {
			return fmt.Errorf("创建支付记录失败: %w", err)
		}

		// 更新订单状态
		now := time.Now()
		if err := tx.Model(&order).Updates(map[string]interface{}{
			"status":   model.OrderStatusPaid,
			"pay_time": &now,
		}).Error; err != nil {
			return fmt.Errorf("更新订单状态失败: %w", err)
		}

		// 创建服务器记录
		expireTime := time.Now().AddDate(0, req.Period, 0)
		server := model.Server{
			UserID:     req.UserID,
			OrderID:    order.ID,
			ProviderID: product.ProviderID,
			ProductID:  req.ProductID,
			Name:       req.Name,
			InstanceID: fmt.Sprintf("placeholder-%d", time.Now().Unix()),
			Region:     product.Region,
			Zone:       product.Zone,
			Status:     model.ServerStatusCreating,
			ExpireTime: expireTime,
			AutoRenew:  req.AutoRenew,
			Password:   req.Password,
			OSType:     product.OS,
			CPU:        product.CPU,
			Memory:     product.Memory,
			Storage:    product.Storage,
			Bandwidth:  product.Bandwidth,
			Traffic:    product.Traffic,
		}

		if err := tx.Create(&server).Error; err != nil {
			return fmt.Errorf("创建服务器记录失败: %w", err)
		}

		result = &PurchaseServerResponse{
			OrderID:   order.ID,
			OrderNo:   orderNo,
			PayAmount: payAmount,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetServerDetailRequest 获取服务器详情请求
type GetServerDetailRequest struct {
	ServerID uint `json:"server_id"`
	UserID   uint `json:"user_id"`
}

// GetServerDetail 获取服务器详情
func (s *ServerService) GetServerDetail(ctx context.Context, req *GetServerDetailRequest) (*model.Server, error) {
	var server model.Server
	if err := s.db.Preload("Provider").Preload("Product").Preload("Order").
		Where("id = ? AND user_id = ?", req.ServerID, req.UserID).
		First(&server).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("服务器不存在")
		}
		return nil, fmt.Errorf("获取服务器详情失败: %w", err)
	}

	return &server, nil
}

// StartServerRequest 启动服务器请求
type StartServerRequest struct {
	ServerID uint `json:"server_id"`
	UserID   uint `json:"user_id"`
}

// StopServerRequest 停止服务器请求
type StopServerRequest struct {
	ServerID uint `json:"server_id"`
	UserID   uint `json:"user_id"`
}

// RestartServerRequest 重启服务器请求
type RestartServerRequest struct {
	ServerID uint `json:"server_id"`
	UserID   uint `json:"user_id"`
}