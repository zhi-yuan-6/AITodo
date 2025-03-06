package controllers

import (
	"AITodo/services/ai_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AIRequest struct {
	Input string `json:"input"`
}

func AIAssistant(c *gin.Context) {
	var req AIRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	/*/// 异步处理AI请求
	ai_service.ProcessTaskWithAIAsync(req.Input, func(response string, err error) {
		if err != nil {
			// 通过通道或上下文通知错误（此处简化处理）
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": response})
	})
	// 立即返回202 Accepted（表示请求已接受，正在处理）
	c.JSON(http.StatusAccepted, gin.H{"message": "请求正在处理中"})*/

	response, err := ai_service.ProcessTaskWithAI(req.Input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//这里可以添加将AI相应转化为Task的逻辑
	c.JSON(http.StatusOK, gin.H{
		"message": "AI processed successfully",
		"data":    response,
	})
}
