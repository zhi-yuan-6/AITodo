package controllers

import (
	"AITodo/dto"
	"AITodo/services/ai_service"
	"AITodo/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AIAssistant(c *gin.Context) {
	var req dto.AIRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	/*//这里可以添加将AI相应转化为Task的逻辑
	c.JSON(http.StatusOK, gin.H{
		"message": "AI processed successfully",
		"data":    response,
	})*/

	// 设置响应头，支持流式输出
	setStreamHeaders(c)

	// 调用流式AI服务
	err = ai_service.StreamFunctionCalling(response, c.Writer)
	if err != nil {
		writeStreamError(c, err)
		return
	}

	// 流式输出完成后，发送结束信号
	c.Writer.Write([]byte("data: [DONE]\n\n"))
	c.Writer.Flush()
}

// 生成数据分析报告
func AIAnalytics(c *gin.Context) {
	var req dto.Analytics
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
			"content": "角色设定：你是一位经验丰富的资深行为模式事项管理分析师，我将提供给你一些最近事项相关数据，你需结合这些数据为我编写一篇对用户有意义的分析报告和建议，且要求你不能将报告划分为任务完成趋势分析、类别分布分析、时间热力图分析等部分，而要从整体上进行综合分析，给出一篇完整的分析报告，报告内容要包括针对用户未来行为的可行性和实用性兼具的建议，以帮助用户更好地理解和优化其行为模式。",
		},
		{
			"role":    "user",
			"content": string(combinedResponse),
		},
	}

	// 设置响应头，支持流式输出
	setStreamHeaders(c)

	// 调用流式AI服务
	err = ai_service.StreamFunctionCalling(messages, c.Writer)
	if err != nil {
		writeStreamError(c, err)
		return
	}

	// 流式输出完成后，发送结束信号
	c.Writer.Write([]byte("data: [DONE]\n\n"))
	c.Writer.Flush()

	/*//非流式输出
	completion, err := ai_service.FunctionCalling(messages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(completion)
	content := completion.Get("choices.0.message.content").Raw

	//这里可以添加将AI相应转化为Task的逻辑
	c.JSON(http.StatusOK, gin.H{
		"message": "AI processed successfully",
		"data":    content,
	})*/
}

// setStreamHeaders configures the response headers for SSE streaming
func setStreamHeaders(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Writer.Flush()
}

// writeStreamError writes an error message to the stream
func writeStreamError(c *gin.Context, err error) {
	errorMsg := fmt.Sprintf("data: {\"error\": \"%s\"}\n\n", err.Error())
	c.Writer.Write([]byte(errorMsg))
	c.Writer.Flush()
}
