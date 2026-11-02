package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	tokenTypeAccess  = "access"
	tokenTypeRefresh = "refresh"
)

type JwtClaims struct {
	UserID    string `json:"user_id"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"` // token版本
	jwt.RegisteredClaims
}

func CreateAccessToken(signingKey []byte, claims JwtClaims, expiresIn int) (string, error) {
	return CreateToken(signingKey, claims, tokenTypeAccess, expiresIn)
}
func CreateRefreshToken(signingKey []byte, claims JwtClaims, expiresIn int) (string, error) {
	return CreateToken(signingKey, claims, tokenTypeRefresh, expiresIn)
}
func CreateToken(signingKey []byte, claims JwtClaims, tokenType string, expiresIn int) (string, error) {
	now := time.Now()
	claims.TokenType = tokenType
	claims.IssuedAt = jwt.NewNumericDate(now)
	claims.NotBefore = jwt.NewNumericDate(now)
	claims.ExpiresAt = jwt.NewNumericDate(now.Add(time.Second * time.Duration(expiresIn)))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		return "", fmt.Errorf("sign token err: %w", err)
	}
	return signedToken, nil
}

// -------------------------------

func GetClaims(tokenString string, signingKey []byte) (*JwtClaims, error) {
	token, err := ParseToken(tokenString, signingKey)
	if err != nil {
		return nli, err
	}
	claims, ok := token.Claims.(*JwtClaims)
	if !ok {
		return nil, fmt.Errorf("parse claims err: %w", err)
	}

	return claims, nil
}
func ParseToken(tokenString string, signingKey []byte) (*jwt.Token, error) {
	// 校验token是否为空
	if tokenString == "" {
		return nil, fmt.Errorf("token is empty")
	}

	// 使用jwt库解析token
	token, err := jwt.ParseWithClaims(
		tokenString,  // 待解析的jwt字符串
		&JwtClaims{}, // 自定义声明结构体，用于存储token中的用户信息
		func(token *jwt.Token) (any, error) { // 密钥回调函数 用于验证签名
			// 验证签名是否为 HS256（确保token使用的是预期的加密算法）
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
			}
			// 返回密钥签名 用于验证token的签名是否有效
			return signingKey, nil
		},
		// 配置选项：
		// 1. 限制只接受HS256 签名算法的 token（与上面回调中的检查冗余 但增加了一层安全保障）
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		// 2. 要求 token 必须包含过期时间（exp声明），如果缺失或已过期则验证失败
		jwt.WithExpirationRequired(),
	)

	// 处理解析错误（包括签名错误、过期错误、格式错误等等）
	if err != nil {
		return nil, fmt.Errorf("parse token err: %w", err)
	}
	return token, nil
}
