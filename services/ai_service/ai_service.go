package ai_service

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"time"
)

type AIResult struct {
	Response string
	Err      error
}

type ToolFunction func(map[string]interface{}) (interface{}, error)

/*// 异步处理入口
func ProcessTaskWithAIAsync(input string, callback func(string, error)) {
	//提交任务到线程池
	_ = globals.TaskPool.Submit(func() {
		response, err := ProcessTaskWithAI(input)
		callback(response, err)
	})
}*/

func ProcessTaskWithAI(userID uint, input string) (string, error) {
	//两种处理方式：1.是先使用大模型对字符串进行分析analyze，判断用户是什么操作，是创建、更新还是删除。如果是创建，则再调用模型返回创建所需要的参数；若是更新或删除，则首先调用搜索任务函数searchTask，让大模型返回搜索所需要的参数关键字string，然后搜所函数返回id和新的任务参数（若为删除则为空字段），然后再根据是更新还是删除调用对应函数。总共调用两次模型（searchTask先查）
	//2.先让大模型分析用户输入的字符串，然后选择调用的函数，返回对应函数所需要的参数（函数名称和task）；根据函数名称调用对应函数，若为创建，则直接调用创建函数，若为更新或删除则根据返回的task参数调用searchTask函数（也可使用大模型实现检索，或者普通的方式），若返回为空则直接进行返回，查到了对应的id则进行删除或更新操作，总共调用一次模型（searchTask后查）

	messages := []map[string]interface{}{
		{
			"role": "system",
			"content": "你是一个智能任务管理助手，用户可能会给你发他要做的事情，也可能会给你发一个通知（如班级、工作里发的通知等），你需要智能的帮助用户去分析，根据用户的需求执行以下操作,：" +
				"创建任务：当用户请求添加或创建任务时，调用 CreateTask 函数。" +
				"删除任务：当用户请求删除或取消任务时，或者使用不想不要这种带否定的字段时，调用 DeleteTask 函数。" +
				"更新任务：当用户请求修改或更新任务时，调用 UpdateTask 函数。" +
				"在调用 CreateTask 或 DeleteTask 时，确保返回的参数符合用户描述，并结合历史对话中的任务对象。如果任务与历史任务相关且某些字段值已存在，请尽量保持这些字段值一致。" +
				"时间处理逻辑：" +
				"时间格式要求：所有时间字段的返回格式必须为：YYYY-MM-DD HH:MM:SS（例如：2006-01-02 15:04:05）。" +
				"明确时间：如果用户提供了明确的时间（例如：明天下午三点、下周五上午十点），直接返回相应的 start_date 和 due_date。" +
				"模糊时间和时间段：" +
				"“今天”：即当前日期，start_date 为当天的00:00，due_date 为当天的23:59。\n“明天”：即明天的日期，start_date 为明天的00:00，due_date 为明天的23:59。\n“后天”：即后天的日期，start_date 为后天的00:00，due_date 为后天的23:59。\n“上午/中午/下午/晚上”：这些时间段应根据常规认知来设置。" +
				"例如：\n上午：从 06:00 到 12:00；\n中午：从 12:00 到 14:00；\n下午：从 14:00 到 18:00；\n晚上：从 18:00 到 22:00；\n深夜/凌晨：从 22:00 到 06:00。\n“周一到周五的工作日”：如果没有指定具体日期，可以根据“工作日”的常识推测任务。例如：\n“周一到周五的上午”：从周一到周五的 06:00 到 12:00；\n“周一到周五的下午”：从周一到周五的 12:00 到 18:00。\n时间段与重复性任务：\n\n如果用户提到一个时间段（例如：“明天从下午三点到五点有个会议”），应创建两个任务：\n任务1：start_date 为明天下午3点，due_date 为明天下午5点。\n如果用户提到重复性任务（例如：“我每周三都开会”），应按周期设置任务。\n不确定时间：\n\n如果用户只说了大概的时间（例如：“下午开会”），则默认返回该时间段（如：14:00 到 18:00）。\n如果用户提到的是模糊时间点（例如：“下周某天开会”），则可以推测出合理的时间范围，默认设置任务时间为该日的00:00到23:59。\n注意事项：\n时间段：尽量遵循常见的社会和文化认知，推测合理的时间范围。\n时区处理：确保所有时间都以用户本地时区为准，自动转换时区差异。\n时间上的模糊性：当用户没有明确指定具体时间时，尽量根据常识和常规习惯给出合理的时间段。\n",
		},
		{
			"role":    "user",
			"content": fmt.Sprintf("当前时间是%v,你需要根据当前时间来帮助我来对任务进行管理", time.Now()),
		},
		{
			"role":    "user",
			"content": input,
		},
	}

	// 工具函数映射表
	functionMapper := map[string]ToolFunction{
		"CreateTask": func(args map[string]interface{}) (interface{}, error) {
			return adaptCreateTask(userID, args) // ✅ 闭包传递 userID
		},
		"UpdateTask": adaptUpdateTask,
		"DeleteTask": adaptDeleteTask,
	}

	//添加for循环但是，会多次调用大模型损失效率
	//for {
	completion, err := FunctionCalling(messages)
	if err != nil {
		return "", err
	}

	message := completion.Get("choices.0.message")
	toolCalls := message.Get("tool_calls")

	// 如果没有工具调用则直接返回
	if len(toolCalls.Array()) == 0 {
		return message.Get("content").String(), nil
	}

	// 添加助手的工具调用消息到上下文
	assistantMsg := createAssistantMessage(message, toolCalls)
	messages = append(messages, assistantMsg)

	fmt.Println(toolCalls)
	// 收集所有工具调用
	for _, toolCall := range toolCalls.Array() {
		response, err := processSingleToolCall(toolCall, functionMapper)
		// 记录错误但继续处理
		if err != nil {
			response = map[string]interface{}{
				"role":                    "tool",
				"content":                 err.Error(),
				"tool_call_function_name": toolCall.Get("function.name").String(),
				"is_error":                false,
			}
		}

		messages = append(messages, response)
		logToolResponse(response)

	}

	messages = append(messages, map[string]interface{}{
		"role": "system",
		"content": "请根据工具执行结果生成最终响应：" +
			"1. 汇总所有成功操作\n" +
			"2. 列出所有失败操作及原因\n" +
			"3. 使用自然语言组织成用户友好的回复",
	})
	// 获取最终总结
	finalCompletion, err := FunctionCalling(messages)
	if err != nil {
		return "", err
	}

	return finalCompletion.Get("choices.0.message.content").String(), nil
}

//}

// 创建助手消息
func createAssistantMessage(message gjson.Result, toolCalls gjson.Result) map[string]interface{} {
	msg := map[string]interface{}{
		"role":       "assistant",
		"content":    message.Get("content").String(),
		"tool_calls": make([]interface{}, 0),
	}

	for _, tc := range toolCalls.Array() {
		msg["tool_calls"] = append(msg["tool_calls"].([]interface{}), map[string]interface{}{
			"id":   tc.Get("id").String(),
			"type": "function",
			"function": map[string]interface{}{
				"name":      tc.Get("function.name").String(),
				"arguments": tc.Get("function.arguments").Str,
			},
		})
	}
	return msg
}

func processSingleToolCall(toolCall gjson.Result, mapper map[string]ToolFunction) (map[string]interface{}, error) {
	functionName := toolCall.Get("function.name").String()
	argumentsString := toolCall.Get("function.arguments").Str

	// 参数解析
	var arguments map[string]interface{}
	if err := json.Unmarshal([]byte(argumentsString), &arguments); err != nil {
		return buildErrorResponse(toolCall, "参数解析失败: %v", err)
	}

	// 任务检索处理
	if needsTaskLookup(functionName) {
		taskID, err := searchTask(arguments)
		if err != nil {
			return buildErrorResponse(toolCall, "任务查找失败: %v", err)
		}
		arguments["id"] = taskID
	}

	// 执行目标函数
	function, exists := mapper[functionName]
	if !exists {
		return buildErrorResponse(toolCall, "未知函数: %s", functionName)
	}

	output, err := function(arguments)
	if err != nil {
		return buildErrorResponse(toolCall, "执行失败: %v", err)
	}

	return map[string]interface{}{
		"role":                    "tool",
		"content":                 fmt.Sprintf("%v", output),
		"tool_call_function_name": toolCall.Get("function.name").String(),
	}, nil
}

// 构建错误响应
func buildErrorResponse(toolCall gjson.Result, format string, args ...interface{}) (map[string]interface{}, error) {
	errMsg := fmt.Sprintf(format, args...)
	return map[string]interface{}{
		"role":                    "tool",
		"content":                 errMsg,
		"tool_call_function_name": toolCall.Get("function.name").String(),
		"is_error":                false,
	}, fmt.Errorf(errMsg)
}

// 判断是否需要任务检索
func needsTaskLookup(funcName string) bool {
	return funcName == "DeleteTask" || funcName == "UpdateTask"
}

func logToolResponse(resp map[string]interface{}) {
	status := "true"
	if _, ok := resp["is_error"]; ok {
		status = "false"
	}
	fmt.Printf("[%s] 工具调用 %v → %s\n",
		status,
		resp["tool_call_function_name"],
		truncateString(resp["content"].(string), 100))
}

func truncateString(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
