package jwt

import (
	"errors"
	"time"
	"vetrue-vben-api/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	RoleID   int64  `json:"role_id"`
	jwt.RegisteredClaims
}

// GenerateJwtToken 生成JWT Token
func GenerateJwtToken(data map[string]interface{}) (string, error) {
	expire := config.AppConfig.JWT.Expire
	if expire == 0 {
		expire = 86400 // 默认24小时
	}

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expire) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	if id, ok := data["id"].(int64); ok {
		claims.ID = id
	}
	if username, ok := data["username"].(string); ok {
		claims.Username = username
	}
	if roleID, ok := data["role_id"].(int64); ok {
		claims.RoleID = roleID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := config.AppConfig.JWT.SecretKey
	if secretKey == "" {
		secretKey = "default_secret_key"
	}

	return token.SignedString([]byte(secretKey))
}

// ParseToken 解析JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	secretKey := config.AppConfig.JWT.SecretKey
	if secretKey == "" {
		secretKey = "default_secret_key"
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
