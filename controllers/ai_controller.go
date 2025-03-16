package controllers

import (
	"AITodo/dto"
	"AITodo/services/ai_service"
	"AITodo/util"
	"encoding/json"
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

	uid, err := util.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	response, err := ai_service.ProcessTaskWithAI(uid, req.Input)
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

type Analytics struct {
	TimeRange struct {
		Start    string `json:"start"`
		End      string `json:"end"`
		Interval string `json:"interval"`
	} `json:"timeRange"`
	AnalyticsData dto.CombinedResponse `json:"analyticsData"`
}

// 生成数据分析报告
func AIAnalytics(c *gin.Context) {
	var req Analytics
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	combinedResponse, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal req" + err.Error()})
	}

	messages := []map[string]interface{}{
		{
			"role":    "system",
			"content": "你是一数据分析师，接下来用户会提供给你他的事项在一段时间内的情况，包括任务完成趋势、任务分类分布、用户活跃时段的信息，你需要根据这些信息生成一篇分析报告以及未来建议，要求：1.报告要有意义，2.以markdown格式返回结果",
		},
		{
			"role":    "user",
			"content": string(combinedResponse),
		},
	}
	completion, err := ai_service.FunctionCalling(messages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	content := completion.Get("choices.0.message.content").Raw
	//fmt.Println(content)
	//这里可以添加将AI相应转化为Task的逻辑
	c.JSON(http.StatusOK, gin.H{
		"message": "AI processed successfully",
		"data":    content,
	})
}
