package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// 密码工具
type PasswordUtils struct{}

// 创建密码工具实例
func NewPasswordUtils() *PasswordUtils {
	return &PasswordUtils{}
}

// 生成密码哈希
func (p *PasswordUtils) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// 验证密码
func (p *PasswordUtils) CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// 生成随机密码
func (p *PasswordUtils) GenerateRandomPassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	
	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}
	
	return string(bytes), nil
}

// 验证密码强度
func (p *PasswordUtils) ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("密码长度不能少于8个字符")
	}
	
	if len(password) > 50 {
		return fmt.Errorf("密码长度不能超过50个字符")
	}
	
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	
	if !hasUpper || !hasLower || !hasDigit {
		return fmt.Errorf("密码必须包含大写字母、小写字母和数字")
	}
	
	return nil
}

// 验证邮箱格式
func ValidateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("邮箱格式不正确")
	}
	return nil
}

// 验证用户名格式
func ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 20 {
		return fmt.Errorf("用户名长度必须在3-20个字符之间")
	}
	
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !usernameRegex.MatchString(username) {
		return fmt.Errorf("用户名只能包含字母、数字、下划线和连字符")
	}
	
	return nil
}

// 验证手机号格式
func ValidatePhone(phone string) error {
	phoneRegex := regexp.MustCompile(`^1[3456789]\d{9}$`)
	if !phoneRegex.MatchString(phone) {
		return fmt.Errorf("手机号格式不正确")
	}
	return nil
}

// 生成API密钥
func GenerateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	
	hasher := sha256.New()
	hasher.Write(bytes)
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// 脱敏处理
func MaskSensitiveData(data, dataType string) string {
	switch strings.ToLower(dataType) {
	case "email":
		parts := strings.Split(data, "@")
		if len(parts) != 2 {
			return data
		}
		username := parts[0]
		domain := parts[1]
		if len(username) <= 3 {
			return data
		}
		return username[:2] + "***" + username[len(username)-1:] + "@" + domain
	case "phone":
		if len(data) != 11 {
			return data
		}
		return data[:3] + "****" + data[7:]
	case "realname":
		if len(data) <= 2 {
			return data
		}
		return data[:1] + "**" + data[len(data)-1:]
	default:
		return data
	}
}