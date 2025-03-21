package ai_service

import (
	"AITodo/dto"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

var transport = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout: 10 * time.Second,
	}).DialContext,
	ResponseHeaderTimeout: 15 * time.Second,
}

// HTTP Client 复用，避免每次创建新的 client
var client = &http.Client{
	//Timeout: 60 * time.Second, // 设置超时
	Transport: transport,
}

// functionCalling 发送请求并解析 DeepSeek API 响应
func FunctionCalling(messages []map[string]interface{}) (gjson.Result, error) {
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
								"category": map[string]interface{}{
									"type":        "string",
									"description": "任务类别，从{工作、学习、生活、健身、其他}这四个标签中选择一个，如果无法判断则默认选择其他，必填",
								},
								"location": map[string]interface{}{
									"type":        "string",
									"description": "地点，从用户的描述中提取地点，可选",
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
								"category": map[string]interface{}{
									"type":        "string",
									"description": "任务类别，从{工作、学习、生活、健身、其他}这四个标签中选择一个，如果无法判断则默认选择其他，必填",
								},
								"location": map[string]interface{}{
									"type":        "string",
									"description": "地点，从用户的描述中提取地点，可选",
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
									"description": "任务开始日期，必须按照如下示例格式填写：2006-01-02 15:04:05，必填",
								},
								"due_date": map[string]interface{}{
									"type":        "string",
									"description": "任务截止日期，必须按照如下示例格式填写：2006-01-02 15:04:05，必填",
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
								"category": map[string]interface{}{
									"type":        "string",
									"description": "任务类别，从{工作、学习、生活、健身、其他}这四个标签中选择一个，如果无法判断则默认选择其他，必填",
								},
								"location": map[string]interface{}{
									"type":        "string",
									"description": "地点，从用户的描述中提取地点，可选",
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
									"description": "任务开始日期，必须按照如下示例格式填写：2006-01-02 15:04:05，必填",
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
	}

	// 构造请求体
	requestBody := dto.RequestBody{
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

// StreamFunctionCalling 流式处理AI响应并将其直接写入HTTP响应
func StreamFunctionCalling(messages []map[string]interface{}, writer io.Writer) error {
	// 构造请求体
	requestBody := dto.RequestBody{
		Model:             "qwen-plus",
		Messages:          messages,
		ParallelToolCalls: true,
		Stream:            true,
	}

	// 将请求体转为 JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	// 创建 POST 请求
	req, err := http.NewRequest("POST", "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	apiKey := os.Getenv("DASHSCOPE_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("API key is missing")
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		bodyText, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API Error: %s\nResponse: %s", resp.Status, string(bodyText))
	}

	/*// 使用 bufio.Reader 读取响应流
	reader := io.Reader(resp.Body)
	buf := make([]byte, 1024)

	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read response: %w", err)
		}

		if n > 0 {
			// 将读取到的流式数据写入 HTTP 响应
			_, err = writer.Write([]byte(fmt.Sprintf("data: %s\n\n", string(buf[:n]))))
			if flusher, ok := writer.(http.Flusher); ok {
				flusher.Flush()
			}
			if err != nil {
				return fmt.Errorf("failed to write response: %w", err)
			}
		}
	}

	return nil*/
	// 按行读取响应流
	scanner := bufio.NewScanner(resp.Body)

	// 处理流式响应数据，转换为SSE格式并实时刷新到客户端
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue // 跳过空行
		}
		// 写入完整的 data: 事件
		if !strings.HasPrefix(line, "data: ") {
			line = "data: " + line
		}
		_, err = writer.Write([]byte(line + "\n\n"))
		if flusher, ok := writer.(http.Flusher); ok {
			flusher.Flush()
		}
		if err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}
	}

	// 检查流读取过程中是否发生错误
	if err := scanner.Err(); err != nil {
		//return fmt.Errorf("failed to read response: %w", err)
		errorMsg := fmt.Sprintf("data: {\"error\": \"%s\"}\n\n", err.Error())
		writer.Write([]byte(errorMsg))
		if flusher, ok := writer.(http.Flusher); ok {
			flusher.Flush()
		}
		return err
	}
	return nil
}
