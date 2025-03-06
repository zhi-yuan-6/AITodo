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

		if isWhitelist(c) {
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

		//将用户信息存入上下文
		c.Set("usrID", claims.ID)

		c.Next() //继续后续处理
	}
}

/*func isWhitelist(path string) bool {
	whitelist := map[string]bool{
		"/user/login":    true,
		"/user/register": true,
	}
	return whitelist[path]
}*/

// 支持 HTTP 方法 + 路径组合
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
}
