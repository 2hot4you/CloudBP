package database

import (
	"fmt"
	"path/filepath"
	"cloudbp-backend/internal/config"
	"cloudbp-backend/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Shanghai",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 自动迁移数据库表
	if err := model.AutoMigrate(db); err != nil {
		return nil, fmt.Errorf("数据库迁移失败: %w", err)
	}

	// 初始化基础数据
	if err := InitData(db); err != nil {
		return nil, fmt.Errorf("初始化数据失败: %w", err)
	}
	
	return db, nil
}

// InitWithMigration 使用迁移管理器初始化数据库
func InitWithMigration(cfg config.DatabaseConfig, migrationPath string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Shanghai",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 使用迁移管理器
	if migrationPath != "" {
		migrationManager := NewMigrationManager(db, migrationPath)
		
		// 初始化迁移表
		if err := migrationManager.Init(); err != nil {
			return nil, fmt.Errorf("初始化迁移表失败: %w", err)
		}
		
		// 执行迁移
		if err := migrationManager.Up(); err != nil {
			return nil, fmt.Errorf("执行迁移失败: %w", err)
		}
	} else {
		// 自动迁移数据库表
		if err := model.AutoMigrate(db); err != nil {
			return nil, fmt.Errorf("数据库迁移失败: %w", err)
		}
	}

	// 初始化基础数据
	if err := InitData(db); err != nil {
		return nil, fmt.Errorf("初始化数据失败: %w", err)
	}
	
	return db, nil
}

// GetMigrationManager 获取迁移管理器
func GetMigrationManager(db *gorm.DB, basePath string) *MigrationManager {
	migrationPath := filepath.Join(basePath, "migrations")
	return NewMigrationManager(db, migrationPath)
}