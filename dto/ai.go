package dto

type AIRequest struct {
	Input string `json:"input"`
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
	Stream            bool                     `json:"stream"`
}

type Analytics struct {
	TimeRange struct {
		Start    string `json:"start"`
		End      string `json:"end"`
		Interval string `json:"interval"`
	} `json:"timeRange"`
	AnalyticsData CombinedResponse `json:"analyticsData"`
}
