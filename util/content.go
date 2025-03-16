package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) (uint, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, fmt.Errorf("用户未认证")
	}

	uid, ok := userID.(uint)
	if !ok {
		return 0, fmt.Errorf("用户 ID 格式错误")
	}
	return uid, nil
}
