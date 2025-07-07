package model

import (
	"time"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"unique;not null" json:"username"`
	Email     string         `gorm:"unique;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	Phone     string         `gorm:"unique" json:"phone"`
	RealName  string         `json:"real_name"`
	Avatar    string         `json:"avatar"`
	Status    int            `gorm:"default:1" json:"status"` // 1:正常 2:禁用
	Role      string         `gorm:"default:user" json:"role"` // user, admin
	Balance   float64        `gorm:"default:0" json:"balance"` // 账户余额
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	// 关联
	Orders    []Order    `gorm:"foreignKey:UserID" json:"orders,omitempty"`
	Servers   []Server   `gorm:"foreignKey:UserID" json:"servers,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}