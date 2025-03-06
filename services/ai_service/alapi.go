package ai_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"os"
	"time"
)

// HTTP Client 复用，避免每次创建新的 client
var client = &http.Client{
	Timeout: 10 * time.Second, // 设置超时
}

// ToolCalls 结构体
type ToolCalls struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Index    int    `json:"index"`
	Function struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"` // 动态解析 arguments
	} `json:"function"`
}

// Message 结构体
type Message struct {
	Content   string      `json:"content,omitempty"`
	Role      string      `json:"role,omitempty"`
	ToolCalls []ToolCalls `json:"tool_calls,omitempty"`
}

// Function 结构体
type Function struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

// 请求体结构体
type RequestBody struct {
	Model             string                   `json:"model"`
	Messages          []map[string]interface{} `json:"messages"`
	Tools             []map[string]interface{} `json:"tools"`
	ParallelToolCalls bool                     `json:"parallel_tool_calls"`
}

// functionCalling 发送请求并解析 DeepSeek API 响应
func functionCalling(messages []map[string]interface{}) (gjson.Result, error) {
	// 定义工具列表
	var tools = []map[string]interface{}{
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "GetNowTime",
				"description": "获取当前时间",
				"parameters": map[string]interface{}{
					"description": "该函数不接受任何参数",
					"type":        "object",
					"properties":  map[string]interface{}{},
				},
			},
		},
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "CreateTask",
				"description": "创建新任务",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"task": map[string]interface{}{
							"type":        "models.Task",
							"description": "任务对象",
							"properties": map[string]interface{}{
								"title": map[string]interface{}{
									"type":        "string",
									"description": "任务标题，不要带时间描述的字段，长度不超过255字符，必填",
								},
								"description": map[string]interface{}{
									"type":        "string",
									"description": "任务描述，可选",
								},
								"status": map[string]interface{}{
									"type":        "string",
									"description": "任务状态，默认为'pending'，可选",
								},
								"start_date": map[string]interface{}{
									"type":        "string",
									"description": "任务开始日期，必须按照如下示例格式填写：2006-01-02 15:04:05,必填",
								},
								"due_date": map[string]interface{}{
									"type":        "string",
									"description": "任务截止日期，必须按照如下示例格式填写：2006-01-02 15:04:05，必填",
								},
							},
							"required": []string{"title", "due_date"},
						},
					},
					"required": []string{"task"},
				},
			},
		},

		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "UpdateTask",
				"description": "更新任务",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"task": map[string]interface{}{
							"type":        "*models.Task",
							"description": "任务对象指针，包含以下可更新字段，在提供符合用户描述的字段值时，需结合历史对话中的任务对象。如果发现当前任务与历史任务相关且字段值已存在，请尽量保持字段值与之前一致。",
							"properties": map[string]interface{}{
								"id": map[string]interface{}{
									"type":        "uint",
									"description": "任务ID，恒为0，必填",
								},
								"title": map[string]interface{}{
									"type":        "string",
									"description": "任务标题，不要带时间描述的字段，长度不超过255字符，必填",
								},
								"description": map[string]interface{}{
									"type":        "string",
									"description": "任务描述，可选",
								},
								"status": map[string]interface{}{
									"type":        "string",
									"description": "任务状态，可选",
								},
								"start_date": map[string]interface{}{
									"type":        "string",
									"description": "任务开始日期，必须按照如下示例格式填写：2006-01-02 15:04:05，必填，但可为空",
								},
								"due_date": map[string]interface{}{
									"type":        "string",
									"description": "任务截止日期，必须按照如下示例格式填写：2006-01-02 15:04:05，必填，但可为空",
								},
							},
						},
					},
					"required": []string{"task"},
				},
			},
		},
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "DeleteTask",
				"description": "删除任务",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"task": map[string]interface{}{
							"type":        "models.Task",
							"description": "任务对象，在提供符合用户描述的字段值时，需结合历史对话中的任务对象。如果发现当前任务与历史任务相关且字段值已存在，请尽量保持字段值与之前一致。",
							"properties": map[string]interface{}{
								"id": map[string]interface{}{
									"type":        "uint",
									"description": "任务ID，恒为0，必填",
								},
								"title": map[string]interface{}{
									"type":        "string",
									"description": "任务标题，不要带时间描述的字段，应尽量简短并且可以表达清楚用户要求，长度不超过255字符，必填",
								},
								"description": map[string]interface{}{
									"type":        "string",
									"description": "任务描述，应尽量简短，提取用户描述中的关键字，可选",
								},
								"status": map[string]interface{}{
									"type":        "string",
									"description": "任务状态，默认为'pending'，可选",
								},
								"start_date": map[string]interface{}{
									"type":        "string",
									"description": "任务开始日期，必须按照如下示例格式填写：2006-01-02 15:04:05，必填，但可为空",
								},
								"due_date": map[string]interface{}{
									"type":        "string",
									"description": "任务截止日期，必须按照如下示例格式填写：2006-01-02 15:04:05，必填，但可为空",
								},
							},
							"required": []string{"title", "due_date"},
						},
					},
					"required": []string{"task"},
				},
			},
		},
	}

	// 构造请求体
	requestBody := RequestBody{
		Model:             "qwen-plus",
		Messages:          messages,
		Tools:             tools,
		ParallelToolCalls: true,
	}

	// 将请求体转为 JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return gjson.Result{}, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// 创建 POST 请求
	req, err := http.NewRequest("POST", "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return gjson.Result{}, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	apiKey := os.Getenv("DASHSCOPE_API_KEY")
	if apiKey == "" {
		return gjson.Result{}, fmt.Errorf("API key is missing")
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return gjson.Result{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		bodyText, _ := io.ReadAll(resp.Body)
		return gjson.Result{}, fmt.Errorf("API Error: %s\nResponse: %s", resp.Status, string(bodyText))
	}

	// 读取响应体
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return gjson.Result{}, fmt.Errorf("failed to read response body: %w", err)
	}

	// 使用 gjson 解析响应体
	completion := gjson.ParseBytes(bodyBytes)
	if !completion.Exists() {
		return gjson.Result{}, fmt.Errorf("failed to parse response body")
	}

	// 提取 choices 数组
	choices := completion.Get("choices").Array()
	if len(choices) == 0 {
		return gjson.Result{}, fmt.Errorf("no response from AI")
	}

	return completion, nil
}
