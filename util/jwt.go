package util

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"sync"
	"time"
)

const (
	TokenExpiration   = 72 * time.Hour
	RefreshExpiration = 7 * 24 * time.Hour
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// CustomClaims 自定义声明结构体
type CustomClaims struct {
	UserID               uint `json:"user_id"` // 用户ID，JSON序列化时字段名为user_id
	jwt.RegisteredClaims      // 嵌入jwt库中已注册的声明
}

// jwtSecret 密钥 (ECDSA 私钥)
var (
	JwtSecret interface{}
	PublicKey interface{}
	jwtMutex  = &sync.RWMutex{}
)

// GenerateJWT 生成令牌和刷新令牌
func GenerateJWT(userID uint) (*TokenPair, error) {
	jwtMutex.RLock()
	defer jwtMutex.RUnlock()
	// 创建自定义声明实例，设置用户ID和过期时间
	accessClaims := CustomClaims{
		UserID: userID, // 设置用户ID
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpiration)), // 设置令牌的过期时间为当前时间加72小时 NewNumericDate:时间转换为JWT标准的数值日期格式（自1970年1月1日UTC时间以来的秒数）
			Issuer:    "AITodo",                                            // 设置令牌的发布者
		},
	}
	// 创建新的JWT令牌，使用ES256签名算法和自定义声明
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodES256, accessClaims).SignedString(JwtSecret)
	if err != nil {
		return nil, err
	}

	//生成RefreshToken
	refreshClims := CustomClaims{
		UserID: userID, // 设置用户ID
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshExpiration)),
			Issuer:    "AITodo", // 设置令牌的发布者
		},
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodES256, refreshClims).SignedString(JwtSecret)
	if err != nil {
		return nil, err
	}

	//存储到Redis
	if err := StoreToken(userID, refreshToken, TokenExpiration); err != nil {
		return nil, fmt.Errorf("刷新令牌存储失败:%w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// 解析并验证令牌有效性
func ParseAndValidateJWT(tokenString string) (*CustomClaims, error) {
	claims, err := ParseJWT(tokenString)
	if err != nil {
		return nil, err
	}
	//检查redis中的令牌有效性
	StoredToken, err := GetToken(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("令牌不存在")
	}
	if StoredToken != tokenString {
		return nil, fmt.Errorf("令牌已失效")
	}
	return claims, nil
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

	// 处理解析错误
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, fmt.Errorf("令牌格式错误:%v", err)
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return nil, fmt.Errorf("令牌已过期或尚未生效")
			} else {
				return nil, fmt.Errorf("令牌验证错误: %s", err.Error())
			}
		}
		return nil, fmt.Errorf("解析令牌失败: %s", err.Error())
	}

	// 验证令牌是否有效
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("令牌无效")
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
