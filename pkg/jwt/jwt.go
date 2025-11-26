package jwt

import (
	"errors"
	"time"

	"github.com/JokerYuan-lang/MyNoteBook/internal/config"
	"github.com/golang-jwt/jwt/v4"
)

type MyClaims struct {
	Username string `json:"username"`
	UserID   uint   `json:"user_id"`
	jwt.RegisteredClaims
}

//生成token

func GenerateToken(userID uint, username string, conf config.JwtConfig) (string, error) {
	// 设置过期时间
	expireTime := time.Now().Add(time.Duration(conf.Expire) * time.Hour)
	// 构造声明
	claims := MyClaims{
		Username: username,
		UserID:   userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()), // 签发时间
			Issuer:    "MyNoteBook",                   // 签发者
		},
	}
	// 生成 Token（HS256 算法）
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 签名（使用配置中的密钥）
	return token.SignedString([]byte(conf.Secret))
}

//解析token

// ParseToken 解析 JWT Token
func ParseToken(tokenString string, conf config.JwtConfig) (*MyClaims, error) {
	// 解析 Token（指定声明类型和签名密钥）
	token, err := jwt.ParseWithClaims(
		tokenString,
		&MyClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// 验证签名算法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("不支持的签名算法")
			}
			return []byte(conf.Secret), nil
		},
	)
	if err != nil {
		return nil, err
	}
	// 验证 Token 有效性并转换声明类型
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("token 无效")
}
