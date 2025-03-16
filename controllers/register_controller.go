package controllers

import (
	"AITodo/config"
	"AITodo/dto"
	"AITodo/services"
	"AITodo/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// 基本注册
func RegisterHandler(c *gin.Context) {
	//解析用户传递信息
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求格式"})
		return
	}

	//调用服务处创建用户
	user, err := services.RegisterService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("注册失败,失败原因:%s", err.Error())})
		return
	}

	tokenPair, err := util.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"生成token失败:": err.Error()})
	}

	c.JSON(http.StatusOK, dto.RegisterResponse{
		TokenPair: tokenPair,
		User:      *user,
	})
}

// 发送验证码api
func SendVerifyCodeHandler(c *gin.Context) {
	//解析参数
	phone := c.Query("phone")
	if !util.IsPhone(phone) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "手机号格式不正确"})
		return
	}

	// 初始化短信服务
	smsService, err := util.NewSMSService(&config.Cfg.SMS)
	if err != nil {
		log.Fatalf("SMS service initialization failed: %v", err)
	}

	//生成code
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	if err := util.StoreCode(phone, code, 60*time.Minute); err != nil {
		log.Fatalf("验证码保存失败: %v", err)
	}

	//发送验证码
	err = smsService.SendVerificationCode(phone, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("验证码发送失败:%s", err.Error())})
		return
	}

	// 立即返回成功
	c.JSON(http.StatusOK, gin.H{"message": "验证码已发送"})

}
