package model

import (
	"time"
	"gorm.io/gorm"
)

// Provider 云服务提供商模型
type Provider struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`           // 腾讯云、阿里云等
	Code        string         `gorm:"unique;not null" json:"code"`    // tencent、aliyun等
	Logo        string         `json:"logo"`                           // 提供商logo
	Description string         `json:"description"`                    // 描述
	Status      int            `gorm:"default:1" json:"status"`        // 1:启用 2:禁用
	Config      string         `gorm:"type:text" json:"config"`        // JSON格式的配置信息
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	
	// 关联
	Products []Product `gorm:"foreignKey:ProviderID" json:"products,omitempty"`
}

// TableName 指定表名
func (Provider) TableName() string {
	return "providers"
}

// Product 产品模型（针对腾讯云Lighthouse）
type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ProviderID  uint           `gorm:"not null" json:"provider_id"`    // 提供商ID
	Name        string         `gorm:"not null" json:"name"`           // 产品名称
	Code        string         `gorm:"not null" json:"code"`           // 产品代码
	Type        string         `gorm:"not null" json:"type"`           // lighthouse、cvm等
	Region      string         `gorm:"not null" json:"region"`         // 地域
	Zone        string         `json:"zone"`                           // 可用区
	CPU         int            `gorm:"not null" json:"cpu"`            // CPU核心数
	Memory      int            `gorm:"not null" json:"memory"`         // 内存GB
	Storage     int            `gorm:"not null" json:"storage"`        // 存储GB
	StorageType string         `json:"storage_type"`                   // 存储类型
	Bandwidth   int            `json:"bandwidth"`                      // 带宽Mbps
	Traffic     int            `json:"traffic"`                        // 流量包GB
	OS          string         `json:"os"`                             // 操作系统
	Price       float64        `gorm:"not null" json:"price"`          // 价格/月
	OriginalPrice float64      `json:"original_price"`                 // 原价
	Status      int            `gorm:"default:1" json:"status"`        // 1:上架 2:下架
	Description string         `gorm:"type:text" json:"description"`   // 产品描述
	Features    string         `gorm:"type:text" json:"features"`      // 特性JSON
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	
	// 关联
	Provider Provider `gorm:"foreignKey:ProviderID" json:"provider,omitempty"`
}

// TableName 指定表名
func (Product) TableName() string {
	return "products"
}