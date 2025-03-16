package controllers

import (
	"AITodo/dto"
	"AITodo/services"
	"AITodo/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// 注:可添加手机号或者邮箱验证码登陆
func LoginHandler(c *gin.Context) {
	//获取用户账户前端传递用账号和密码
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if strings.Contains(err.Error(), "未激活") {
			c.JSON(http.StatusForbidden, gin.H{"该用户已经注销": err.Error()})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"无效的请求格式:": err.Error()})
		}
		return
	}

	//根据用户账户查询密码，判断密码是否匹配
	user, err := services.LoginService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"登陆失败:": err.Error()})
		return
	}

	//生成JWT
	tokenPair, err := util.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"token生成失败:": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.LoginResponse{
		TokenPair: tokenPair,
		User:      *user,
	})
}
