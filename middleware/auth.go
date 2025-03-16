package middleware

import (
	"AITodo/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		//非白名单路由 执行jwt验证逻辑
		//从handler中获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "未提供认证令牌",
			})
			return
		}

		//解析token
		//对于一个标准的 JWT 认证头，其格式通常是： Authorization: Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "令牌格式错误"})
			return
		}

		//验证jwt
		claims, err := util.ParseJWT(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		valid, err := util.IsTokenValid(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if !valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "令牌已经失效"})
		}

		//将用户信息存入上下文
		c.Set("user_id", claims.UserID)

		c.Next() //继续后续处理
	}
}

// middleware/auth.go
func RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			RefreshToken string `json:"refresh_token"`
		}

		// 1. 参数校验
		if err := c.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "无效请求格式"})
			return
		}

		// 2. 解析令牌
		claims, err := util.ParseAndValidateJWT(req.RefreshToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// 3. 生成新令牌
		newTokens, err := util.GenerateJWT(claims.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "令牌刷新失败"})
			return
		}

		c.JSON(200, newTokens)
	}
}

/*// 支持 HTTP 方法 + 路径组合
var whitelist = map[string]map[string]bool{
	"/user/send_code": {
		"GET":     true,
		"OPTIONS": true, // 添加OPTIONS支持
	},
	"/user/register": {
		"POST":    true,
		"OPTIONS": true,
	},
	"/user/login": {
		"POST":    true,
		"OPTIONS": true,
	},
}

func isWhitelist(c *gin.Context) bool {
	methods, ok := whitelist[c.Request.URL.Path]
	if !ok {
		return false
	}
	return methods[c.Request.Method]
}*/
