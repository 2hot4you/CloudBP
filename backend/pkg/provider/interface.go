package provider

import (
	"context"
	"time"
)

// CloudProvider 云厂商接口
type CloudProvider interface {
	// 获取厂商名称
	GetName() string
	
	// 获取厂商代码
	GetCode() string
	
	// 创建实例
	CreateInstance(ctx context.Context, req *CreateInstanceRequest) (*CreateInstanceResponse, error)
	
	// 删除实例
	DeleteInstance(ctx context.Context, req *DeleteInstanceRequest) error
	
	// 启动实例
	StartInstance(ctx context.Context, req *StartInstanceRequest) error
	
	// 停止实例
	StopInstance(ctx context.Context, req *StopInstanceRequest) error
	
	// 重启实例
	RestartInstance(ctx context.Context, req *RestartInstanceRequest) error
	
	// 获取实例详情
	GetInstanceDetail(ctx context.Context, req *GetInstanceDetailRequest) (*GetInstanceDetailResponse, error)
	
	// 获取实例列表
	GetInstanceList(ctx context.Context, req *GetInstanceListRequest) (*GetInstanceListResponse, error)
	
	// 获取实例监控数据
	GetInstanceMonitor(ctx context.Context, req *GetInstanceMonitorRequest) (*GetInstanceMonitorResponse, error)
	
	// 重置实例密码
	ResetInstancePassword(ctx context.Context, req *ResetInstancePasswordRequest) error
	
	// 重装实例系统
	RebuildInstance(ctx context.Context, req *RebuildInstanceRequest) error
	
	// 获取可用区域列表
	GetRegions(ctx context.Context) (*GetRegionsResponse, error)
	
	// 获取可用镜像列表
	GetImages(ctx context.Context, req *GetImagesRequest) (*GetImagesResponse, error)
	
	// 获取可用规格列表
	GetInstanceTypes(ctx context.Context, req *GetInstanceTypesRequest) (*GetInstanceTypesResponse, error)
}

// CreateInstanceRequest 创建实例请求
type CreateInstanceRequest struct {
	Name         string `json:"name"`          // 实例名称
	Region       string `json:"region"`        // 地域
	Zone         string `json:"zone"`          // 可用区
	ImageID      string `json:"image_id"`      // 镜像ID
	InstanceType string `json:"instance_type"` // 实例规格
	Password     string `json:"password"`      // 登录密码
	Period       int    `json:"period"`        // 购买时长（月）
	AutoRenew    bool   `json:"auto_renew"`    // 是否自动续费
	UserData     string `json:"user_data"`     // 用户数据
}

// CreateInstanceResponse 创建实例响应
type CreateInstanceResponse struct {
	InstanceID string `json:"instance_id"` // 实例ID
	OrderID    string `json:"order_id"`    // 订单ID
}

// DeleteInstanceRequest 删除实例请求
type DeleteInstanceRequest struct {
	InstanceID string `json:"instance_id"` // 实例ID
}

// StartInstanceRequest 启动实例请求
type StartInstanceRequest struct {
	InstanceID string `json:"instance_id"` // 实例ID
}

// StopInstanceRequest 停止实例请求
type StopInstanceRequest struct {
	InstanceID string `json:"instance_id"` // 实例ID
}

// RestartInstanceRequest 重启实例请求
type RestartInstanceRequest struct {
	InstanceID string `json:"instance_id"` // 实例ID
}

// GetInstanceDetailRequest 获取实例详情请求
type GetInstanceDetailRequest struct {
	InstanceID string `json:"instance_id"` // 实例ID
}

// GetInstanceDetailResponse 获取实例详情响应
type GetInstanceDetailResponse struct {
	InstanceID   string    `json:"instance_id"`   // 实例ID
	Name         string    `json:"name"`          // 实例名称
	State        string    `json:"state"`         // 实例状态
	Region       string    `json:"region"`        // 地域
	Zone         string    `json:"zone"`          // 可用区
	PublicIP     string    `json:"public_ip"`     // 公网IP
	PrivateIP    string    `json:"private_ip"`    // 私网IP
	OSType       string    `json:"os_type"`       // 操作系统类型
	OSName       string    `json:"os_name"`       // 操作系统名称
	CPU          int       `json:"cpu"`           // CPU核数
	Memory       int       `json:"memory"`        // 内存GB
	Storage      int       `json:"storage"`       // 存储GB
	Bandwidth    int       `json:"bandwidth"`     // 带宽Mbps
	Traffic      int       `json:"traffic"`       // 流量包GB
	UsedTraffic  int       `json:"used_traffic"`  // 已用流量GB
	CreatedTime  time.Time `json:"created_time"`  // 创建时间
	ExpiredTime  time.Time `json:"expired_time"`  // 到期时间
	InstanceType string    `json:"instance_type"` // 实例规格
	ImageID      string    `json:"image_id"`      // 镜像ID
}

// GetInstanceListRequest 获取实例列表请求
type GetInstanceListRequest struct {
	Region string `json:"region"` // 地域
	Offset int    `json:"offset"` // 偏移量
	Limit  int    `json:"limit"`  // 限制数量
}

// GetInstanceListResponse 获取实例列表响应
type GetInstanceListResponse struct {
	Instances  []GetInstanceDetailResponse `json:"instances"`   // 实例列表
	TotalCount int                         `json:"total_count"` // 总数量
}

// GetInstanceMonitorRequest 获取实例监控数据请求
type GetInstanceMonitorRequest struct {
	InstanceID  string    `json:"instance_id"`  // 实例ID
	MetricName  string    `json:"metric_name"`  // 监控指标名称
	StartTime   time.Time `json:"start_time"`   // 开始时间
	EndTime     time.Time `json:"end_time"`     // 结束时间
	Period      int       `json:"period"`       // 数据粒度（秒）
}

// GetInstanceMonitorResponse 获取实例监控数据响应
type GetInstanceMonitorResponse struct {
	MetricName string                `json:"metric_name"` // 监控指标名称
	DataPoints []MonitorDataPoint    `json:"data_points"` // 监控数据点
}

// MonitorDataPoint 监控数据点
type MonitorDataPoint struct {
	Timestamp time.Time `json:"timestamp"` // 时间戳
	Value     float64   `json:"value"`     // 值
	Unit      string    `json:"unit"`      // 单位
}

// ResetInstancePasswordRequest 重置实例密码请求
type ResetInstancePasswordRequest struct {
	InstanceID  string `json:"instance_id"`  // 实例ID
	NewPassword string `json:"new_password"` // 新密码
}

// RebuildInstanceRequest 重装实例系统请求
type RebuildInstanceRequest struct {
	InstanceID string `json:"instance_id"` // 实例ID
	ImageID    string `json:"image_id"`    // 镜像ID
	Password   string `json:"password"`    // 登录密码
}

// GetRegionsResponse 获取可用区域列表响应
type GetRegionsResponse struct {
	Regions []Region `json:"regions"` // 区域列表
}

// Region 区域信息
type Region struct {
	RegionID   string `json:"region_id"`   // 区域ID
	RegionName string `json:"region_name"` // 区域名称
	Zones      []Zone `json:"zones"`       // 可用区列表
}

// Zone 可用区信息
type Zone struct {
	ZoneID   string `json:"zone_id"`   // 可用区ID
	ZoneName string `json:"zone_name"` // 可用区名称
}

// GetImagesRequest 获取可用镜像列表请求
type GetImagesRequest struct {
	Region   string `json:"region"`    // 地域
	OSType   string `json:"os_type"`   // 操作系统类型
	ImageType string `json:"image_type"` // 镜像类型
}

// GetImagesResponse 获取可用镜像列表响应
type GetImagesResponse struct {
	Images []Image `json:"images"` // 镜像列表
}

// Image 镜像信息
type Image struct {
	ImageID     string `json:"image_id"`     // 镜像ID
	ImageName   string `json:"image_name"`   // 镜像名称
	OSType      string `json:"os_type"`      // 操作系统类型
	OSName      string `json:"os_name"`      // 操作系统名称
	ImageSize   int    `json:"image_size"`   // 镜像大小
	ImageType   string `json:"image_type"`   // 镜像类型
	Description string `json:"description"`  // 描述
}

// GetInstanceTypesRequest 获取可用规格列表请求
type GetInstanceTypesRequest struct {
	Region string `json:"region"` // 地域
}

// GetInstanceTypesResponse 获取可用规格列表响应
type GetInstanceTypesResponse struct {
	InstanceTypes []InstanceType `json:"instance_types"` // 实例规格列表
}

// InstanceType 实例规格信息
type InstanceType struct {
	InstanceType string  `json:"instance_type"` // 实例规格
	CPU          int     `json:"cpu"`           // CPU核数
	Memory       int     `json:"memory"`        // 内存GB
	Storage      int     `json:"storage"`       // 存储GB
	Bandwidth    int     `json:"bandwidth"`     // 带宽Mbps
	Traffic      int     `json:"traffic"`       // 流量包GB
	Price        float64 `json:"price"`         // 价格
	Description  string  `json:"description"`   // 描述
}

// ProviderConfig 厂商配置
type ProviderConfig struct {
	SecretID  string `json:"secret_id"`  // 密钥ID
	SecretKey string `json:"secret_key"` // 密钥
	Region    string `json:"region"`     // 默认地域
	Endpoint  string `json:"endpoint"`   // API端点
}

// ProviderManager 厂商管理器
type ProviderManager struct {
	providers map[string]CloudProvider
}

// NewProviderManager 创建厂商管理器
func NewProviderManager() *ProviderManager {
	return &ProviderManager{
		providers: make(map[string]CloudProvider),
	}
}

// RegisterProvider 注册厂商
func (pm *ProviderManager) RegisterProvider(provider CloudProvider) {
	pm.providers[provider.GetCode()] = provider
}

// GetProvider 获取厂商
func (pm *ProviderManager) GetProvider(code string) (CloudProvider, bool) {
	provider, exists := pm.providers[code]
	return provider, exists
}

// GetAllProviders 获取所有厂商
func (pm *ProviderManager) GetAllProviders() map[string]CloudProvider {
	return pm.providers
}