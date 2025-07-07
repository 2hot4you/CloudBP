package service

import (
	"cloudbp-backend/internal/model"
	"cloudbp-backend/pkg/auth"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// 用户服务
type UserService struct {
	db          *gorm.DB
	redis       *redis.Client
	jwtManager  *auth.JWTManager
	passwordUtils *auth.PasswordUtils
}

// 创建用户服务
func NewUserService(db *gorm.DB, redis *redis.Client, jwtManager *auth.JWTManager) *UserService {
	return &UserService{
		db:          db,
		redis:       redis,
		jwtManager:  jwtManager,
		passwordUtils: auth.NewPasswordUtils(),
	}
}

// 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Phone    string `json:"phone"`
	RealName string `json:"real_name"`
}

// 登录响应
type LoginResponse struct {
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	ExpiresIn    int64      `json:"expires_in"`
	User         *model.User `json:"user"`
}

// 用户登录
func (s *UserService) Login(req *LoginRequest) (*LoginResponse, error) {
	var user model.User
	
	// 查找用户（支持用户名、邮箱、手机号登录）
	if err := s.db.Where("username = ? OR email = ? OR phone = ?", req.Username, req.Username, req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	
	// 检查用户状态
	if user.Status != model.UserStatusActive {
		return nil, errors.New("用户已被禁用")
	}
	
	// 验证密码
	if err := s.passwordUtils.CheckPassword(req.Password, user.Password); err != nil {
		return nil, errors.New("密码错误")
	}
	
	// 生成JWT令牌
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}
	
	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}
	
	// 脱敏处理
	user.Password = ""
	user.Email = auth.MaskSensitiveData(user.Email, "email")
	user.Phone = auth.MaskSensitiveData(user.Phone, "phone")
	user.RealName = auth.MaskSensitiveData(user.RealName, "realname")
	
	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(time.Hour.Seconds()),
		User:         &user,
	}, nil
}

// 用户注册
func (s *UserService) Register(req *RegisterRequest) error {
	// 验证用户名格式
	if err := auth.ValidateUsername(req.Username); err != nil {
		return err
	}
	
	// 验证邮箱格式
	if err := auth.ValidateEmail(req.Email); err != nil {
		return err
	}
	
	// 验证手机号格式（如果提供）
	if req.Phone != "" {
		if err := auth.ValidatePhone(req.Phone); err != nil {
			return err
		}
	}
	
	// 验证密码强度
	if err := s.passwordUtils.ValidatePasswordStrength(req.Password); err != nil {
		return err
	}
	
	// 检查用户名是否存在
	var count int64
	if err := s.db.Model(&model.User{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("用户名已存在")
	}
	
	// 检查邮箱是否存在
	if err := s.db.Model(&model.User{}).Where("email = ?", req.Email).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("邮箱已存在")
	}
	
	// 检查手机号是否存在（如果提供）
	if req.Phone != "" {
		if err := s.db.Model(&model.User{}).Where("phone = ?", req.Phone).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("手机号已存在")
		}
	}
	
	// 加密密码
	hashedPassword, err := s.passwordUtils.HashPassword(req.Password)
	if err != nil {
		return err
	}
	
	// 创建用户
	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Phone:    req.Phone,
		RealName: req.RealName,
		Status:   model.UserStatusActive,
		Role:     model.UserRoleUser,
		Balance:  0,
	}
	
	if err := s.db.Create(&user).Error; err != nil {
		return err
	}
	
	return nil
}

// 刷新令牌
func (s *UserService) RefreshToken(refreshToken string) (*LoginResponse, error) {
	accessToken, newRefreshToken, err := s.jwtManager.RefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	
	// 验证令牌并获取用户信息
	claims, err := s.jwtManager.ValidateToken(accessToken)
	if err != nil {
		return nil, err
	}
	
	var user model.User
	if err := s.db.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		return nil, err
	}
	
	// 脱敏处理
	user.Password = ""
	user.Email = auth.MaskSensitiveData(user.Email, "email")
	user.Phone = auth.MaskSensitiveData(user.Phone, "phone")
	user.RealName = auth.MaskSensitiveData(user.RealName, "realname")
	
	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(time.Hour.Seconds()),
		User:         &user,
	}, nil
}

// 获取用户信息
func (s *UserService) GetProfile(userID uint) (*model.User, error) {
	var user model.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	
	// 脱敏处理
	user.Password = ""
	user.Email = auth.MaskSensitiveData(user.Email, "email")
	user.Phone = auth.MaskSensitiveData(user.Phone, "phone")
	user.RealName = auth.MaskSensitiveData(user.RealName, "realname")
	
	return &user, nil
}

// 更新用户信息
type UpdateProfileRequest struct {
	Avatar   string `json:"avatar"`
	RealName string `json:"real_name"`
	Phone    string `json:"phone"`
}

func (s *UserService) UpdateProfile(userID uint, req *UpdateProfileRequest) error {
	// 验证手机号格式（如果提供）
	if req.Phone != "" {
		if err := auth.ValidatePhone(req.Phone); err != nil {
			return err
		}
		
		// 检查手机号是否被其他用户使用
		var count int64
		if err := s.db.Model(&model.User{}).Where("phone = ? AND id != ?", req.Phone, userID).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("手机号已被其他用户使用")
		}
	}
	
	// 更新用户信息
	updates := map[string]interface{}{
		"avatar":    req.Avatar,
		"real_name": req.RealName,
		"phone":     req.Phone,
	}
	
	if err := s.db.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		return err
	}
	
	return nil
}

// 修改密码
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func (s *UserService) ChangePassword(userID uint, req *ChangePasswordRequest) error {
	var user model.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return err
	}
	
	// 验证旧密码
	if err := s.passwordUtils.CheckPassword(req.OldPassword, user.Password); err != nil {
		return errors.New("原密码错误")
	}
	
	// 验证新密码强度
	if err := s.passwordUtils.ValidatePasswordStrength(req.NewPassword); err != nil {
		return err
	}
	
	// 加密新密码
	hashedPassword, err := s.passwordUtils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}
	
	// 更新密码
	if err := s.db.Model(&user).Update("password", hashedPassword).Error; err != nil {
		return err
	}
	
	return nil
}