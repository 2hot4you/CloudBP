package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"cloudbp-backend/internal/config"
	"cloudbp-backend/pkg/database"
	"cloudbp-backend/pkg/logger"
)

func main() {
	var (
		action = flag.String("action", "status", "迁移操作: up, down, status, create")
		name   = flag.String("name", "", "迁移名称 (用于create操作)")
	)
	flag.Parse()

	// 初始化日志
	logger.Init("info")

	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("加载配置文件失败: %v\n", err)
		os.Exit(1)
	}

	// 连接数据库
	db, err := database.Init(cfg.Database)
	if err != nil {
		fmt.Printf("连接数据库失败: %v\n", err)
		os.Exit(1)
	}

	// 获取迁移管理器
	workDir, _ := os.Getwd()
	migrationPath := filepath.Join(workDir, "migrations")
	migrationManager := database.NewMigrationManager(db, migrationPath)

	// 初始化迁移表
	if err := migrationManager.Init(); err != nil {
		fmt.Printf("初始化迁移表失败: %v\n", err)
		os.Exit(1)
	}

	// 执行操作
	switch *action {
	case "up":
		if err := migrationManager.Up(); err != nil {
			fmt.Printf("执行迁移失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("所有迁移已成功应用")

	case "status":
		if err := migrationManager.Status(); err != nil {
			fmt.Printf("查看迁移状态失败: %v\n", err)
			os.Exit(1)
		}

	case "create":
		if *name == "" {
			fmt.Println("请指定迁移名称: -name <迁移名称>")
			os.Exit(1)
		}
		if err := migrationManager.CreateMigration(*name); err != nil {
			fmt.Printf("创建迁移失败: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Printf("不支持的操作: %s\n", *action)
		fmt.Println("支持的操作: up, status, create")
		os.Exit(1)
	}
}