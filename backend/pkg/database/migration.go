package database

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Migration 迁移记录
type Migration struct {
	ID        uint      `gorm:"primaryKey"`
	Version   string    `gorm:"unique;not null"`
	Name      string    `gorm:"not null"`
	Applied   bool      `gorm:"default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// MigrationManager 迁移管理器
type MigrationManager struct {
	db            *gorm.DB
	migrationPath string
}

// NewMigrationManager 创建迁移管理器
func NewMigrationManager(db *gorm.DB, migrationPath string) *MigrationManager {
	return &MigrationManager{
		db:            db,
		migrationPath: migrationPath,
	}
}

// Init 初始化迁移表
func (m *MigrationManager) Init() error {
	return m.db.AutoMigrate(&Migration{})
}

// GetMigrationFiles 获取迁移文件列表
func (m *MigrationManager) GetMigrationFiles() ([]string, error) {
	files, err := os.ReadDir(m.migrationPath)
	if err != nil {
		return nil, fmt.Errorf("读取迁移目录失败: %w", err)
	}

	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}

	// 按文件名排序
	sort.Strings(sqlFiles)
	return sqlFiles, nil
}

// GetAppliedMigrations 获取已应用的迁移
func (m *MigrationManager) GetAppliedMigrations() (map[string]bool, error) {
	var migrations []Migration
	if err := m.db.Find(&migrations).Error; err != nil {
		return nil, fmt.Errorf("获取已应用迁移失败: %w", err)
	}

	applied := make(map[string]bool)
	for _, migration := range migrations {
		applied[migration.Version] = migration.Applied
	}
	return applied, nil
}

// ParseMigrationVersion 解析迁移版本
func (m *MigrationManager) ParseMigrationVersion(filename string) (string, string, error) {
	parts := strings.SplitN(filename, "_", 2)
	if len(parts) < 2 {
		return "", "", fmt.Errorf("迁移文件名格式错误: %s", filename)
	}

	version := parts[0]
	name := strings.TrimSuffix(parts[1], ".sql")
	
	// 验证版本号格式
	if _, err := strconv.Atoi(version); err != nil {
		return "", "", fmt.Errorf("迁移版本号格式错误: %s", version)
	}

	return version, name, nil
}

// ExecuteMigration 执行迁移
func (m *MigrationManager) ExecuteMigration(filename string) error {
	version, name, err := m.ParseMigrationVersion(filename)
	if err != nil {
		return err
	}

	// 读取SQL文件
	sqlFile := filepath.Join(m.migrationPath, filename)
	content, err := os.ReadFile(sqlFile)
	if err != nil {
		return fmt.Errorf("读取迁移文件失败: %w", err)
	}

	// 在事务中执行迁移
	return m.db.Transaction(func(tx *gorm.DB) error {
		// 执行SQL
		if err := tx.Exec(string(content)).Error; err != nil {
			return fmt.Errorf("执行迁移SQL失败: %w", err)
		}

		// 记录迁移
		migration := Migration{
			Version: version,
			Name:    name,
			Applied: true,
		}
		if err := tx.Create(&migration).Error; err != nil {
			return fmt.Errorf("记录迁移失败: %w", err)
		}

		return nil
	})
}

// Up 执行迁移
func (m *MigrationManager) Up() error {
	// 获取迁移文件
	files, err := m.GetMigrationFiles()
	if err != nil {
		return err
	}

	// 获取已应用的迁移
	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return err
	}

	// 执行未应用的迁移
	for _, file := range files {
		version, _, err := m.ParseMigrationVersion(file)
		if err != nil {
			return err
		}

		if !applied[version] {
			fmt.Printf("正在应用迁移: %s\n", file)
			if err := m.ExecuteMigration(file); err != nil {
				return fmt.Errorf("应用迁移 %s 失败: %w", file, err)
			}
			fmt.Printf("迁移 %s 应用成功\n", file)
		}
	}

	return nil
}

// Status 查看迁移状态
func (m *MigrationManager) Status() error {
	files, err := m.GetMigrationFiles()
	if err != nil {
		return err
	}

	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return err
	}

	fmt.Println("迁移状态:")
	fmt.Println("版本\t\t状态\t\t文件名")
	fmt.Println("----\t\t----\t\t------")

	for _, file := range files {
		version, _, err := m.ParseMigrationVersion(file)
		if err != nil {
			continue
		}

		status := "待应用"
		if applied[version] {
			status = "已应用"
		}

		fmt.Printf("%s\t\t%s\t\t%s\n", version, status, file)
	}

	return nil
}

// CreateMigration 创建新的迁移文件
func (m *MigrationManager) CreateMigration(name string) error {
	// 获取现有迁移文件
	files, err := m.GetMigrationFiles()
	if err != nil {
		return err
	}

	// 生成新的版本号
	version := "001"
	if len(files) > 0 {
		lastFile := files[len(files)-1]
		lastVersion, _, err := m.ParseMigrationVersion(lastFile)
		if err != nil {
			return err
		}
		
		lastNum, err := strconv.Atoi(lastVersion)
		if err != nil {
			return err
		}
		
		version = fmt.Sprintf("%03d", lastNum+1)
	}

	// 创建迁移文件
	filename := fmt.Sprintf("%s_%s.sql", version, name)
	filePath := filepath.Join(m.migrationPath, filename)

	template := fmt.Sprintf(`-- 迁移: %s
-- 版本: %s
-- 创建时间: %s

-- 在此处编写您的迁移SQL

`, name, version, time.Now().Format("2006-01-02 15:04:05"))

	if err := os.WriteFile(filePath, []byte(template), 0644); err != nil {
		return fmt.Errorf("创建迁移文件失败: %w", err)
	}

	fmt.Printf("迁移文件已创建: %s\n", filename)
	return nil
}