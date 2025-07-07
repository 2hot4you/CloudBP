package database

import (
	"cloudbp-backend/internal/model"
	"gorm.io/gorm"
)

// InitData 初始化基础数据
func InitData(db *gorm.DB) error {
	// 创建默认的云服务提供商
	providers := []model.Provider{
		{
			Name:        "腾讯云",
			Code:        "tencent",
			Logo:        "/assets/logos/tencent.png",
			Description: "腾讯云提供稳定可靠的云计算服务",
			Status:      model.ProviderStatusActive,
			Config:      `{"api_endpoint": "https://lighthouse.tencentcloudapi.com", "region": "ap-guangzhou"}`,
		},
		{
			Name:        "阿里云",
			Code:        "aliyun",
			Logo:        "/assets/logos/aliyun.png",
			Description: "阿里云提供全面的云计算服务",
			Status:      model.ProviderStatusInactive,
			Config:      `{"api_endpoint": "https://ecs.aliyuncs.com"}`,
		},
	}

	for _, provider := range providers {
		var existingProvider model.Provider
		if err := db.Where("code = ?", provider.Code).First(&existingProvider).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&provider).Error; err != nil {
				return err
			}
		}
	}

	// 创建腾讯云Lighthouse产品
	var tencentProvider model.Provider
	if err := db.Where("code = ?", "tencent").First(&tencentProvider).Error; err != nil {
		return err
	}

	lighthouseProducts := []model.Product{
		{
			ProviderID:    tencentProvider.ID,
			Name:          "轻量应用服务器 1核2G",
			Code:          "lighthouse-1c2g",
			Type:          model.TencentProductTypeLighthouse,
			Region:        "ap-guangzhou",
			Zone:          "ap-guangzhou-3",
			CPU:           1,
			Memory:        2,
			Storage:       50,
			StorageType:   "SSD",
			Bandwidth:     3,
			Traffic:       100,
			OS:            "Ubuntu 20.04",
			Price:         24.00,
			OriginalPrice: 24.00,
			Status:        model.ProductStatusOnline,
			Description:   "适合个人开发者和小型应用",
			Features:      `{"support_docker": true, "support_ssh": true, "backup": true}`,
		},
		{
			ProviderID:    tencentProvider.ID,
			Name:          "轻量应用服务器 2核4G",
			Code:          "lighthouse-2c4g",
			Type:          model.TencentProductTypeLighthouse,
			Region:        "ap-guangzhou",
			Zone:          "ap-guangzhou-3",
			CPU:           2,
			Memory:        4,
			Storage:       80,
			StorageType:   "SSD",
			Bandwidth:     5,
			Traffic:       200,
			OS:            "Ubuntu 20.04",
			Price:         54.00,
			OriginalPrice: 54.00,
			Status:        model.ProductStatusOnline,
			Description:   "适合中小型应用和网站",
			Features:      `{"support_docker": true, "support_ssh": true, "backup": true}`,
		},
		{
			ProviderID:    tencentProvider.ID,
			Name:          "轻量应用服务器 2核8G",
			Code:          "lighthouse-2c8g",
			Type:          model.TencentProductTypeLighthouse,
			Region:        "ap-guangzhou",
			Zone:          "ap-guangzhou-3",
			CPU:           2,
			Memory:        8,
			Storage:       100,
			StorageType:   "SSD",
			Bandwidth:     6,
			Traffic:       300,
			OS:            "Ubuntu 20.04",
			Price:         108.00,
			OriginalPrice: 108.00,
			Status:        model.ProductStatusOnline,
			Description:   "适合大型应用和数据库服务",
			Features:      `{"support_docker": true, "support_ssh": true, "backup": true}`,
		},
	}

	for _, product := range lighthouseProducts {
		var existingProduct model.Product
		if err := db.Where("provider_id = ? AND code = ?", product.ProviderID, product.Code).First(&existingProduct).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&product).Error; err != nil {
				return err
			}
		}
	}

	// 创建系统配置
	configs := []model.Config{
		{
			Key:         "site_name",
			Value:       "云服务器销售平台",
			Type:        "string",
			Group:       "system",
			Title:       "网站名称",
			Description: "网站的名称",
			Sort:        1,
			Status:      1,
		},
		{
			Key:         "site_logo",
			Value:       "/assets/logo.png",
			Type:        "string",
			Group:       "system",
			Title:       "网站Logo",
			Description: "网站的Logo路径",
			Sort:        2,
			Status:      1,
		},
		{
			Key:         "register_enabled",
			Value:       "true",
			Type:        "bool",
			Group:       "user",
			Title:       "允许注册",
			Description: "是否允许用户注册",
			Sort:        1,
			Status:      1,
		},
		{
			Key:         "email_verify_enabled",
			Value:       "false",
			Type:        "bool",
			Group:       "user",
			Title:       "邮箱验证",
			Description: "是否需要邮箱验证",
			Sort:        2,
			Status:      1,
		},
	}

	for _, config := range configs {
		var existingConfig model.Config
		if err := db.Where("key = ?", config.Key).First(&existingConfig).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&config).Error; err != nil {
				return err
			}
		}
	}

	return nil
}