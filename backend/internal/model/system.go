package model

import (
	"time"
	"gorm.io/gorm"
)

// Config 系统配置模型
type Config struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Key       string         `gorm:"unique;not null" json:"key"`    // 配置键
	Value     string         `gorm:"type:text" json:"value"`        // 配置值
	Type      string         `gorm:"not null" json:"type"`          // 配置类型: string、int、float、bool、json
	Group     string         `gorm:"not null" json:"group"`         // 配置分组
	Title     string         `gorm:"not null" json:"title"`         // 配置标题
	Description string       `gorm:"type:text" json:"description"`  // 配置描述
	Sort      int            `gorm:"default:0" json:"sort"`         // 排序
	Status    int            `gorm:"default:1" json:"status"`       // 1:启用 2:禁用
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Config) TableName() string {
	return "configs"
}

// OperationLog 操作日志模型
type OperationLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`                      // 用户ID
	Username  string    `json:"username"`                     // 用户名
	Module    string    `gorm:"not null" json:"module"`       // 模块
	Action    string    `gorm:"not null" json:"action"`       // 操作
	Content   string    `gorm:"type:text" json:"content"`     // 操作内容
	IP        string    `json:"ip"`                           // IP地址
	UserAgent string    `gorm:"type:text" json:"user_agent"`  // 用户代理
	Status    int       `gorm:"default:1" json:"status"`      // 1:成功 2:失败
	CreatedAt time.Time `json:"created_at"`
	
	// 关联
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (OperationLog) TableName() string {
	return "operation_logs"
}