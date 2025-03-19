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

	//这里可以添加将AI相应转化为Task的逻辑
	c.JSON(http.StatusOK, gin.H{
		"message": "AI processed successfully",
		"data":    response,
	})
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
			"content": "## 角色设定\n作为资深行为模式分析师，您需要结合数据可视化图表（用户已自行查看）与深度数据洞察，输出具备商业价值的分析报告\n\n## 输入数据\n用户将提供包含以下维度的结构化数据：\n1. 时间范围（含统计间隔单位）\n2. 任务完成趋势（多维度时间序列数据）\n3. 任务分类分布（多分类占比数据）\n4. 用户活跃热力图（时段分布数据）\n\n## 分析框架\n采用金字塔结构进行多维度交叉分析：\n### 趋势分析（必含）\n- 数据波动特征识别（趋势性/周期性/异常点）\n- 关键拐点归因分析（结合时间维度解释突变原因）\n- 多指标关联分析（如完成率与失败率的相关性）\n\n### 分布分析（必含）\n- 帕累托分析（识别关键20%分类）\n- 资源分配合理性评估（时间投入与任务类型的匹配度）\n- 风险类别预警（失败率异常的类别）\n\n### 时段分析（必含）\n- 高效时段识别（完成率高的黄金时间段）\n- 行为模式诊断（活跃时段与任务类型的关联性）\n- 生物钟规律推测（持续出现的活跃时段模式）\n\n## 报告要求\n### 内容规范\n1. 使用「现象描述 → 数据佐证 → 深层解读」三级分析逻辑\n2. 必须包含至少3个跨维度交叉分析结论（如：某分类任务的高失败率是否集中出现在特定时段）\n3. 异常值需给出两种以上可能性解释\n\n### 建议标准\n1. 时间管理优化（基于时段分析）\n2. 任务优先级调整（基于帕累托法则）\n3. 流程改进建议（针对失败任务归因）\n4. 资源配置方案（结合分类分布与完成趋势）\n\n### 格式规范\n1. 采用Markdown层级结构（## 二级标题 / ### 三级标题）\n2. 关键结论使用**加粗**标注\n3. 数据引用格式：`[指标名称]：[数值]（[百分比/变化率]）`\n4. 禁用通用建议，所有建议必须与数据特征直接关联",
		},
		{
			"role":    "user",
			"content": string(combinedResponse),
		},
	}

	// 设置响应头，支持流式输出
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Writer.Flush()

	// 调用流式AI服务
	err = ai_service.StreamFunctionCalling(messages, c.Writer)
	if err != nil {
		// 由于已经开始流式输出，这里只能写入错误信息
		errorMsg := fmt.Sprintf("data: {\"error\": \"%s\"}\n\n", err.Error())
		c.Writer.Write([]byte(errorMsg))
		c.Writer.Flush()
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
