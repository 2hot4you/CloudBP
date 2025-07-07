package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWT配置
type JWTConfig struct {
	SecretKey     string
	Issuer        string
	ExpiresIn     time.Duration
	RefreshTokenExpiresIn time.Duration
}

// JWT管理器
type JWTManager struct {
	config JWTConfig
}

// 自定义Claims
type CustomClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// 创建JWT管理器
func NewJWTManager(config JWTConfig) *JWTManager {
	return &JWTManager{config: config}
}

// 生成访问令牌
func (j *JWTManager) GenerateAccessToken(userID uint, username, role string) (string, error) {
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.config.Issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.config.ExpiresIn)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.config.SecretKey))
}

// 生成刷新令牌
func (j *JWTManager) GenerateRefreshToken(userID uint, username, role string) (string, error) {
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.config.Issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.config.RefreshTokenExpiresIn)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.config.SecretKey))
}

// 验证令牌
func (j *JWTManager) ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("意外的签名方法")
		}
		return []byte(j.config.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("无效的令牌")
	}

	return claims, nil
}

// 刷新令牌
func (j *JWTManager) RefreshToken(tokenString string) (string, string, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", "", err
	}

	// 生成新的访问令牌和刷新令牌
	accessToken, err := j.GenerateAccessToken(claims.UserID, claims.Username, claims.Role)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := j.GenerateRefreshToken(claims.UserID, claims.Username, claims.Role)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}