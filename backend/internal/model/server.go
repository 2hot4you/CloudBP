package model

import (
	"time"
	"gorm.io/gorm"
)

// Server 服务器实例模型
type Server struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `gorm:"not null" json:"user_id"`        // 用户ID
	OrderID      uint           `gorm:"not null" json:"order_id"`       // 订单ID
	ProviderID   uint           `gorm:"not null" json:"provider_id"`    // 提供商ID
	ProductID    uint           `gorm:"not null" json:"product_id"`     // 产品ID
	Name         string         `gorm:"not null" json:"name"`           // 服务器名称
	InstanceID   string         `gorm:"unique;not null" json:"instance_id"` // 云厂商实例ID
	Region       string         `gorm:"not null" json:"region"`         // 地域
	Zone         string         `json:"zone"`                           // 可用区
	PublicIP     string         `json:"public_ip"`                      // 公网IP
	PrivateIP    string         `json:"private_ip"`                     // 私网IP
	Status       string         `gorm:"default:creating" json:"status"` // creating、running、stopped、expired等
	ExpireTime   time.Time      `json:"expire_time"`                    // 到期时间
	AutoRenew    bool           `gorm:"default:false" json:"auto_renew"`// 自动续费
	Password     string         `json:"-"`                              // 登录密码
	OSType       string         `json:"os_type"`                        // 操作系统类型
	OSName       string         `json:"os_name"`                        // 操作系统名称
	CPU          int            `json:"cpu"`                            // CPU核心数
	Memory       int            `json:"memory"`                         // 内存GB
	Storage      int            `json:"storage"`                        // 存储GB
	Bandwidth    int            `json:"bandwidth"`                      // 带宽Mbps
	Traffic      int            `json:"traffic"`                        // 流量包GB
	UsedTraffic  int            `gorm:"default:0" json:"used_traffic"`  // 已用流量GB
	Config       string         `gorm:"type:text" json:"config"`        // JSON格式的配置信息
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	
	// 关联
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Order     Order     `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Provider  Provider  `gorm:"foreignKey:ProviderID" json:"provider,omitempty"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Monitors  []Monitor `gorm:"foreignKey:ServerID" json:"monitors,omitempty"`
}

// TableName 指定表名
func (Server) TableName() string {
	return "servers"
}

// Monitor 监控数据模型
type Monitor struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ServerID   uint      `gorm:"not null" json:"server_id"`   // 服务器ID
	MetricType string    `gorm:"not null" json:"metric_type"` // cpu、memory、disk、network等
	Value      float64   `gorm:"not null" json:"value"`       // 监控值
	Unit       string    `json:"unit"`                        // 单位
	Timestamp  time.Time `gorm:"not null" json:"timestamp"`   // 时间戳
	CreatedAt  time.Time `json:"created_at"`
	
	// 关联
	Server Server `gorm:"foreignKey:ServerID" json:"server,omitempty"`
}

// TableName 指定表名
func (Monitor) TableName() string {
	return "monitors"
}