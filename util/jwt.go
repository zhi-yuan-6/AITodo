package util

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"sync"
	"time"
)

// CustomClaims 自定义声明结构体
type CustomClaims struct {
	UserID               int `json:"user_id"` // 用户ID，JSON序列化时字段名为user_id
	jwt.RegisteredClaims     // 嵌入jwt库中已注册的声明
}

// jwtSecret 密钥 (ECDSA 私钥)
var (
	JwtSecret interface{}
	PublicKey interface{}
	jwtMutex  = &sync.RWMutex{}
)

// GenerateJWT 生成JWT令牌函数
func GenerateJWT(userID int) (string, error) {
	jwtMutex.RLock()
	defer jwtMutex.RUnlock()
	// 创建自定义声明实例，设置用户ID和过期时间
	claims := CustomClaims{
		UserID: userID, // 设置用户ID
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)), // 设置令牌的过期时间为当前时间加72小时
			Issuer:    "AITodo",                                           // 设置令牌的发布者
		},
	}
	// 创建新的JWT令牌，使用ES256签名算法和自定义声明
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	// 使用秘密密钥对令牌进行签名，并返回生成的令牌字符串
	return token.SignedString(JwtSecret)
}

func ParseJWT(tokenString string) (*CustomClaims, error) {
	// 加锁保护
	jwtMutex.RLock()
	defer jwtMutex.RUnlock()

	//解析JWT令牌
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return PublicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("解析令牌失败：%s", err.Error())
	}

	//验证令牌 是否有效
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("令牌失效：%s", err.Error())
	}
}

func SetJwtSecret(secret interface{}) {
	jwtMutex.Lock()
	defer jwtMutex.Unlock()
	JwtSecret = secret
}

func GetJwtSecret() interface{} {
	jwtMutex.RLock()
	defer jwtMutex.RUnlock()
	return JwtSecret
}

func SetPublicKey(pulicKey interface{}) {
	jwtMutex.Lock()
	defer jwtMutex.Unlock()
	PublicKey = pulicKey
}

func GetPublicKey() interface{} {
	jwtMutex.RLock()
	defer jwtMutex.RUnlock()
	return PublicKey
}
