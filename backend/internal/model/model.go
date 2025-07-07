package model

import (
	"gorm.io/gorm"
)

// AutoMigrate 自动迁移数据库表
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Provider{},
		&Product{},
		&Server{},
		&Monitor{},
		&Order{},
		&Payment{},
		&Config{},
		&OperationLog{},
	)
}

// 定义常量
const (
	// 用户状态
	UserStatusActive   = 1
	UserStatusInactive = 2
	
	// 用户角色
	UserRoleUser  = "user"
	UserRoleAdmin = "admin"
	
	// 提供商状态
	ProviderStatusActive   = 1
	ProviderStatusInactive = 2
	
	// 产品状态
	ProductStatusOnline  = 1
	ProductStatusOffline = 2
	
	// 服务器状态
	ServerStatusCreating = "creating"
	ServerStatusRunning  = "running"
	ServerStatusStopped  = "stopped"
	ServerStatusExpired  = "expired"
	ServerStatusError    = "error"
	
	// 订单状态
	OrderStatusPending    = "pending"
	OrderStatusPaid       = "paid"
	OrderStatusProcessing = "processing"
	OrderStatusSuccess    = "success"
	OrderStatusFailed     = "failed"
	OrderStatusCancelled  = "cancelled"
	
	// 订单类型
	OrderTypeNew     = "new"
	OrderTypeRenew   = "renew"
	OrderTypeUpgrade = "upgrade"
	
	// 支付状态
	PaymentStatusPending = "pending"
	PaymentStatusSuccess = "success"
	PaymentStatusFailed  = "failed"
	
	// 支付方式
	PaymentMethodBalance = "balance"
	PaymentMethodWechat  = "wechat"
	PaymentMethodAlipay  = "alipay"
	
	// 监控指标类型
	MetricTypeCPU     = "cpu"
	MetricTypeMemory  = "memory"
	MetricTypeDisk    = "disk"
	MetricTypeNetwork = "network"
	
	// 腾讯云产品类型
	TencentProductTypeLighthouse = "lighthouse"
	TencentProductTypeCVM        = "cvm"
)