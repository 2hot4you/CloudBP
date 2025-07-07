package service

import (
	"context"
	"errors"
	"fmt"

	"cloudbp-backend/internal/model"
	"cloudbp-backend/pkg/provider"
	"gorm.io/gorm"
)

// ProviderService 厂商服务
type ProviderService struct {
	db              *gorm.DB
	providerManager *provider.ProviderManager
}

// NewProviderService 创建厂商服务
func NewProviderService(db *gorm.DB) *ProviderService {
	providerManager := provider.NewProviderManager()
	return &ProviderService{
		db:              db,
		providerManager: providerManager,
	}
}

// InitProviders 初始化厂商
func (s *ProviderService) InitProviders() error {
	// 从数据库获取启用的厂商配置
	var providers []model.Provider
	if err := s.db.Where("status = ?", model.ProviderStatusActive).Find(&providers).Error; err != nil {
		return fmt.Errorf("获取厂商配置失败: %w", err)
	}

	for _, p := range providers {
		if err := s.registerProvider(&p); err != nil {
			return fmt.Errorf("注册厂商 %s 失败: %w", p.Name, err)
		}
	}

	return nil
}

// registerProvider 注册厂商
func (s *ProviderService) registerProvider(p *model.Provider) error {
	switch p.Code {
	case "tencent":
		config, err := provider.ParseTencentCloudConfig(p.Config)
		if err != nil {
			return err
		}

		tencentProvider, err := provider.NewTencentCloudProvider(config)
		if err != nil {
			return err
		}

		s.providerManager.RegisterProvider(tencentProvider)
		return nil

	default:
		return fmt.Errorf("不支持的厂商: %s", p.Code)
	}
}

// GetProvider 获取厂商适配器
func (s *ProviderService) GetProvider(code string) (provider.CloudProvider, error) {
	p, exists := s.providerManager.GetProvider(code)
	if !exists {
		return nil, fmt.Errorf("厂商 %s 不存在或未启用", code)
	}
	return p, nil
}

// CreateInstance 创建实例
func (s *ProviderService) CreateInstance(ctx context.Context, req *CreateInstanceRequest) (*CreateInstanceResponse, error) {
	// 获取厂商信息
	var providerModel model.Provider
	if err := s.db.First(&providerModel, req.ProviderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("厂商不存在")
		}
		return nil, err
	}

	// 获取产品信息
	var product model.Product
	if err := s.db.First(&product, req.ProductID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("产品不存在")
		}
		return nil, err
	}

	// 获取厂商适配器
	cloudProvider, err := s.GetProvider(providerModel.Code)
	if err != nil {
		return nil, err
	}

	// 调用厂商API创建实例
	createReq := &provider.CreateInstanceRequest{
		Name:         req.Name,
		Region:       product.Region,
		Zone:         product.Zone,
		ImageID:      req.ImageID,
		InstanceType: product.Code,
		Password:     req.Password,
		Period:       req.Period,
		AutoRenew:    req.AutoRenew,
		UserData:     req.UserData,
	}

	resp, err := cloudProvider.CreateInstance(ctx, createReq)
	if err != nil {
		return nil, fmt.Errorf("创建实例失败: %w", err)
	}

	return &CreateInstanceResponse{
		InstanceID: resp.InstanceID,
		OrderID:    resp.OrderID,
	}, nil
}

// DeleteInstance 删除实例
func (s *ProviderService) DeleteInstance(ctx context.Context, req *DeleteInstanceRequest) error {
	// 获取服务器信息
	var server model.Server
	if err := s.db.Preload("Provider").First(&server, req.ServerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("服务器不存在")
		}
		return err
	}

	// 获取厂商适配器
	cloudProvider, err := s.GetProvider(server.Provider.Code)
	if err != nil {
		return err
	}

	// 调用厂商API删除实例
	deleteReq := &provider.DeleteInstanceRequest{
		InstanceID: server.InstanceID,
	}

	return cloudProvider.DeleteInstance(ctx, deleteReq)
}

// StartInstance 启动实例
func (s *ProviderService) StartInstance(ctx context.Context, req *StartInstanceRequest) error {
	return s.operateInstance(ctx, req.ServerID, "start")
}

// StopInstance 停止实例
func (s *ProviderService) StopInstance(ctx context.Context, req *StopInstanceRequest) error {
	return s.operateInstance(ctx, req.ServerID, "stop")
}

// RestartInstance 重启实例
func (s *ProviderService) RestartInstance(ctx context.Context, req *RestartInstanceRequest) error {
	return s.operateInstance(ctx, req.ServerID, "restart")
}

// operateInstance 操作实例
func (s *ProviderService) operateInstance(ctx context.Context, serverID uint, operation string) error {
	// 获取服务器信息
	var server model.Server
	if err := s.db.Preload("Provider").First(&server, serverID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("服务器不存在")
		}
		return err
	}

	// 获取厂商适配器
	cloudProvider, err := s.GetProvider(server.Provider.Code)
	if err != nil {
		return err
	}

	// 根据操作类型调用对应的API
	switch operation {
	case "start":
		startReq := &provider.StartInstanceRequest{InstanceID: server.InstanceID}
		return cloudProvider.StartInstance(ctx, startReq)
	case "stop":
		stopReq := &provider.StopInstanceRequest{InstanceID: server.InstanceID}
		return cloudProvider.StopInstance(ctx, stopReq)
	case "restart":
		restartReq := &provider.RestartInstanceRequest{InstanceID: server.InstanceID}
		return cloudProvider.RestartInstance(ctx, restartReq)
	default:
		return fmt.Errorf("不支持的操作: %s", operation)
	}
}

// GetInstanceDetail 获取实例详情
func (s *ProviderService) GetInstanceDetail(ctx context.Context, req *GetInstanceDetailRequest) (*GetInstanceDetailResponse, error) {
	// 获取服务器信息
	var server model.Server
	if err := s.db.Preload("Provider").First(&server, req.ServerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("服务器不存在")
		}
		return nil, err
	}

	// 获取厂商适配器
	cloudProvider, err := s.GetProvider(server.Provider.Code)
	if err != nil {
		return nil, err
	}

	// 调用厂商API获取实例详情
	detailReq := &provider.GetInstanceDetailRequest{
		InstanceID: server.InstanceID,
	}

	resp, err := cloudProvider.GetInstanceDetail(ctx, detailReq)
	if err != nil {
		return nil, fmt.Errorf("获取实例详情失败: %w", err)
	}

	return &GetInstanceDetailResponse{
		InstanceID:   resp.InstanceID,
		Name:         resp.Name,
		State:        resp.State,
		Region:       resp.Region,
		Zone:         resp.Zone,
		PublicIP:     resp.PublicIP,
		PrivateIP:    resp.PrivateIP,
		OSType:       resp.OSType,
		OSName:       resp.OSName,
		CPU:          resp.CPU,
		Memory:       resp.Memory,
		Storage:      resp.Storage,
		Bandwidth:    resp.Bandwidth,
		Traffic:      resp.Traffic,
		UsedTraffic:  resp.UsedTraffic,
		CreatedTime:  resp.CreatedTime,
		ExpiredTime:  resp.ExpiredTime,
		InstanceType: resp.InstanceType,
		ImageID:      resp.ImageID,
	}, nil
}

// SyncInstanceStatus 同步实例状态
func (s *ProviderService) SyncInstanceStatus(ctx context.Context, serverID uint) error {
	// 获取服务器信息
	var server model.Server
	if err := s.db.Preload("Provider").First(&server, serverID).Error; err != nil {
		return err
	}

	// 获取厂商适配器
	cloudProvider, err := s.GetProvider(server.Provider.Code)
	if err != nil {
		return err
	}

	// 获取实例详情
	detailReq := &provider.GetInstanceDetailRequest{
		InstanceID: server.InstanceID,
	}

	resp, err := cloudProvider.GetInstanceDetail(ctx, detailReq)
	if err != nil {
		return err
	}

	// 更新数据库中的服务器信息
	updates := map[string]interface{}{
		"status":       resp.State,
		"public_ip":    resp.PublicIP,
		"private_ip":   resp.PrivateIP,
		"os_type":      resp.OSType,
		"os_name":      resp.OSName,
		"cpu":          resp.CPU,
		"memory":       resp.Memory,
		"storage":      resp.Storage,
		"bandwidth":    resp.Bandwidth,
		"traffic":      resp.Traffic,
		"used_traffic": resp.UsedTraffic,
		"expire_time":  resp.ExpiredTime,
	}

	return s.db.Model(&server).Updates(updates).Error
}

// GetRegions 获取可用区域列表
func (s *ProviderService) GetRegions(ctx context.Context, providerCode string) (*GetRegionsResponse, error) {
	cloudProvider, err := s.GetProvider(providerCode)
	if err != nil {
		return nil, err
	}

	resp, err := cloudProvider.GetRegions(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取区域列表失败: %w", err)
	}

	return &GetRegionsResponse{
		Regions: resp.Regions,
	}, nil
}

// GetImages 获取可用镜像列表
func (s *ProviderService) GetImages(ctx context.Context, req *GetImagesRequest) (*GetImagesResponse, error) {
	cloudProvider, err := s.GetProvider(req.ProviderCode)
	if err != nil {
		return nil, err
	}

	imagesReq := &provider.GetImagesRequest{
		Region:    req.Region,
		OSType:    req.OSType,
		ImageType: req.ImageType,
	}

	resp, err := cloudProvider.GetImages(ctx, imagesReq)
	if err != nil {
		return nil, fmt.Errorf("获取镜像列表失败: %w", err)
	}

	return &GetImagesResponse{
		Images: resp.Images,
	}, nil
}

// GetInstanceTypes 获取可用规格列表
func (s *ProviderService) GetInstanceTypes(ctx context.Context, req *GetInstanceTypesRequest) (*GetInstanceTypesResponse, error) {
	cloudProvider, err := s.GetProvider(req.ProviderCode)
	if err != nil {
		return nil, err
	}

	typesReq := &provider.GetInstanceTypesRequest{
		Region: req.Region,
	}

	resp, err := cloudProvider.GetInstanceTypes(ctx, typesReq)
	if err != nil {
		return nil, fmt.Errorf("获取规格列表失败: %w", err)
	}

	return &GetInstanceTypesResponse{
		InstanceTypes: resp.InstanceTypes,
	}, nil
}

// 请求和响应结构体
type CreateInstanceRequest struct {
	ProviderID uint   `json:"provider_id"`
	ProductID  uint   `json:"product_id"`
	Name       string `json:"name"`
	ImageID    string `json:"image_id"`
	Password   string `json:"password"`
	Period     int    `json:"period"`
	AutoRenew  bool   `json:"auto_renew"`
	UserData   string `json:"user_data"`
}

type CreateInstanceResponse struct {
	InstanceID string `json:"instance_id"`
	OrderID    string `json:"order_id"`
}

type DeleteInstanceRequest struct {
	ServerID uint `json:"server_id"`
}

type StartInstanceRequest struct {
	ServerID uint `json:"server_id"`
}

type StopInstanceRequest struct {
	ServerID uint `json:"server_id"`
}

type RestartInstanceRequest struct {
	ServerID uint `json:"server_id"`
}

type GetInstanceDetailRequest struct {
	ServerID uint `json:"server_id"`
}

type GetInstanceDetailResponse struct {
	InstanceID   string `json:"instance_id"`
	Name         string `json:"name"`
	State        string `json:"state"`
	Region       string `json:"region"`
	Zone         string `json:"zone"`
	PublicIP     string `json:"public_ip"`
	PrivateIP    string `json:"private_ip"`
	OSType       string `json:"os_type"`
	OSName       string `json:"os_name"`
	CPU          int    `json:"cpu"`
	Memory       int    `json:"memory"`
	Storage      int    `json:"storage"`
	Bandwidth    int    `json:"bandwidth"`
	Traffic      int    `json:"traffic"`
	UsedTraffic  int    `json:"used_traffic"`
	CreatedTime  any    `json:"created_time"`
	ExpiredTime  any    `json:"expired_time"`
	InstanceType string `json:"instance_type"`
	ImageID      string `json:"image_id"`
}

type GetRegionsResponse struct {
	Regions []provider.Region `json:"regions"`
}

type GetImagesRequest struct {
	ProviderCode string `json:"provider_code"`
	Region       string `json:"region"`
	OSType       string `json:"os_type"`
	ImageType    string `json:"image_type"`
}

type GetImagesResponse struct {
	Images []provider.Image `json:"images"`
}

type GetInstanceTypesRequest struct {
	ProviderCode string `json:"provider_code"`
	Region       string `json:"region"`
}

type GetInstanceTypesResponse struct {
	InstanceTypes []provider.InstanceType `json:"instance_types"`
}