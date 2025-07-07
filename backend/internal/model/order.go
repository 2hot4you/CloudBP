package model

import (
	"time"
	"gorm.io/gorm"
)

// Order 订单模型
type Order struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"not null" json:"user_id"`         // 用户ID
	OrderNo       string         `gorm:"unique;not null" json:"order_no"` // 订单号
	ProviderID    uint           `gorm:"not null" json:"provider_id"`     // 提供商ID
	ProductID     uint           `gorm:"not null" json:"product_id"`      // 产品ID
	Type          string         `gorm:"not null" json:"type"`            // new、renew、upgrade等
	Status        string         `gorm:"default:pending" json:"status"`   // pending、paid、processing、success、failed、cancelled
	Amount        float64        `gorm:"not null" json:"amount"`          // 订单金额
	DiscountAmount float64       `gorm:"default:0" json:"discount_amount"`// 优惠金额
	PayAmount     float64        `gorm:"not null" json:"pay_amount"`      // 实付金额
	PayMethod     string         `json:"pay_method"`                      // 支付方式
	PayTime       *time.Time     `json:"pay_time"`                        // 支付时间
	Period        int            `gorm:"not null" json:"period"`          // 购买周期(月)
	Quantity      int            `gorm:"default:1" json:"quantity"`       // 数量
	Config        string         `gorm:"type:text" json:"config"`         // JSON格式的配置信息
	Remark        string         `gorm:"type:text" json:"remark"`         // 备注
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	
	// 关联
	User     User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Provider Provider `gorm:"foreignKey:ProviderID" json:"provider,omitempty"`
	Product  Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Servers  []Server `gorm:"foreignKey:OrderID" json:"servers,omitempty"`
}

// TableName 指定表名
func (Order) TableName() string {
	return "orders"
}

// Payment 支付记录模型
type Payment struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	OrderID     uint           `gorm:"not null" json:"order_id"`      // 订单ID
	UserID      uint           `gorm:"not null" json:"user_id"`       // 用户ID
	PaymentNo   string         `gorm:"unique;not null" json:"payment_no"` // 支付单号
	Method      string         `gorm:"not null" json:"method"`        // 支付方式
	Amount      float64        `gorm:"not null" json:"amount"`        // 支付金额
	Status      string         `gorm:"default:pending" json:"status"` // pending、success、failed
	TransactionID string       `json:"transaction_id"`               // 第三方交易号
	PayTime     *time.Time     `json:"pay_time"`                     // 支付时间
	Remark      string         `gorm:"type:text" json:"remark"`      // 备注
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	
	// 关联
	Order Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	User  User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (Payment) TableName() string {
	return "payments"
}