package jwt

import (
	"Personal-blog/configs"
	"Personal-blog/internal/model"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims 自定义 JWT 载荷（包含用户核心信息）
type CustomClaims struct {
	UserID               uint   `json:"user_id"`  // 用户 ID
	Username             string `json:"username"` // 用户名
	jwt.RegisteredClaims        // 内置标准载荷（过期时间、签发者等）
}

// GenerateToken 生成 JWT Token（保持不变）
func GenerateToken(user *model.User) (string, error) {
	claims := CustomClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(configs.Config.App.JWT.ExpiresAt) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    configs.Config.App.JWT.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(configs.Config.App.JWT.Secret))
}

// ParseToken 解析 Token 并验证有效性（修复错误判断逻辑）
func ParseToken(tokenStr string) (*CustomClaims, error) {
	// 1. 定义验证函数（验证签名算法和签发者）
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// 验证签名算法是否为 HMAC-SHA256
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("不支持的签名算法")
			}
			return []byte(configs.Config.App.JWT.Secret), nil
		},
	)

	// 2. 处理解析错误（核心修复部分）
	if err != nil {
		// 使用 errors.Is 判断 v5 的错误类型
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, errors.New("token 已过期")
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, errors.New("token 格式错误")
		case errors.Is(err, jwt.ErrTokenUnverifiable):
			return nil, errors.New("token 无法验证")
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, errors.New("token 签名无效")
		default:
			return nil, errors.New("token 解析失败：" + err.Error())
		}
	}

	// 3. 验证 Token 有效性并返回载荷
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		// 额外验证签发者（可选，增强安全性）
		if claims.Issuer != configs.Config.App.JWT.Issuer {
			return nil, errors.New("token 签发者不匹配")
		}
		return claims, nil
	}

	return nil, errors.New("token 无效")
}
