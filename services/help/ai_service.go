package help

/*//单个函数调用
import (
	"encoding/json"
	"fmt"
)

type ToolFunction func(map[string]interface{}) (interface{}, error)

func ProcessTaskWithAI(input string) (string, error) {
	messages := []map[string]interface{}{
		{
			"role": "system",
			"content": "你是一名任务管理助手，当我像你提问有关事务、事项的事情时，你需要帮助我对我提的要求进行分析，如果我想要添加或创建任务，请调用‘controller.CreateTask’;" +
				"如果我想要删除任务，请调用‘controller.DeleteTask’; 如果我想要对任务进行更新，请调用‘controller.UpdateTask’",
		},
		{
			"role":    "user",
			"content": input,
		},
	}

	completion, err := functionCalling(messages)
	if err != nil {
		return "", err
	}

	// 提取函数名称和参数
	toolCalls := completion.Get("choices.0.message.tool_calls").Array()
	if len(toolCalls) == 0 {
		return "", fmt.Errorf("no tool calls found in response")
	}

	functionName := toolCalls[0].Get("function.name").String()
	argumentsString := toolCalls[0].Get("function.arguments").Str

	// 工具函数映射表
	functionMapper := map[string]ToolFunction{
		"GetAllTasks": adaptGetAllTasks,
		"CreateTask":  adaptCreateTask,
		"UpdateTask":  adaptUpdateTask,
		"DeleteTask":  adaptDeleteTask,
		"GetTaskById": adaptGetTaskById,
	}

	// 查找对应的工具函数
	function, exists := functionMapper[functionName]
	if !exists {
		return "", fmt.Errorf("unknown function: %s", functionName)
	}

	// 执行工具函数
	var arguments map[string]interface{}
	err = json.Unmarshal([]byte(argumentsString), &arguments)
	if err != nil {
		return "", fmt.Errorf("解析参数失败:%v", err)
	}
	fmt.Println(arguments)
	functionOutput, err := function(arguments)
	if err != nil {
		return "", fmt.Errorf("function execution failed: %v", err)
	}

	// 输出函数执行结果
	fmt.Println("Function output:", functionOutput)

	message := completion.Get("choices.0.message")
	var messageMap map[string]interface{}
	err = json.Unmarshal([]byte(message.Raw), &messageMap)
	if err != nil {
		return "", fmt.Errorf("解析message失败: %v", err)
	}
	messages = append(messages, messageMap)
	messages = append(messages, map[string]interface{}{"role": "tool", "content": fmt.Sprintf("%s", functionOutput)})
	// 再次调用 functionCalling，获取响应内容
	completion, err = functionCalling(messages)
	if err != nil {
		return "", err
	}

	return completion.Get("choices.0.message.content").String(), nil
}
*/

/*每次循环只处理第一个工具调用
单次处理原则：

每次循环只处理第一个工具调用（toolCalls[0]）

处理完成后会将工具执行结果加入对话上下文

重新调用模型生成新的响应

实际场景示例：

// 模型第一次响应：
"tool_calls": [
    { "id": "call_1", "function": "CreateTask" },
    { "id": "call_2", "function": "UpdateTask" }
]

// 代码处理：
1. 处理 call_1（创建任务）
2. 将创建结果加入上下文
3. 触发第二次模型调用

// 模型第二次响应：
"tool_calls": [  // 此时是全新的工具调用
    { "id": "call_3", "function": "UpdateTask" }
]

type ToolFunction func(map[string]interface{}) (interface{}, error)

func ProcessTaskWithAI(input string) (string, error) {
	messages := []map[string]interface{}{
		{
			"role": "system",
			"content": "你是一名任务管理助手，当我像你提问有关事务、事项的事情时，你需要帮助我对我提的要求进行分析，如果我想要添加或创建任务，请调用‘CreateTask’;" +
				"如果我想要删除任务，请调用‘DeleteTask’; 如果我想要对任务进行更新，请调用‘UpdateTask’",
		},
		{
			"role":    "user",
			"content": input,
		},
	}

	fmt.Println(strings.Repeat("-", 60))

	// 工具函数映射表
	functionMapper := map[string]ToolFunction{
		"GetAllTasks": adaptGetAllTasks,
		"CreateTask":  adaptCreateTask,
		"UpdateTask":  adaptUpdateTask,
		"DeleteTask":  adaptDeleteTask,
		"GetTaskById": adaptGetTaskById,
	}

	for {
		completion, err := functionCalling(messages)
		if err != nil {
			return "", err
		}
		// 解析模型响应
		message := completion.Get("choices.0.message")
		toolCalls := message.Get("tool_calls").Array()

		// 如果没有工具调用，直接返回内容
		if len(toolCalls) == 0 {
			return message.Get("content").String(), nil
		}

		// 处理第一个工具调用（假设每次只处理一个）
		firstToolCall := toolCalls[0]
		functionName := firstToolCall.Get("function.name").String()
		argumentsString := firstToolCall.Get("function.arguments").Str

		// 将assistant响应加入消息历史
		messages = append(messages, map[string]interface{}{
			"role":    "assistant",
			"content": message.Get("content").String(),
			"tool_calls": []interface{}{
				map[string]interface{}{
					"id":   firstToolCall.Get("id").String(),
					"type": "function",
					"function": map[string]interface{}{
						"name":      functionName,
						"arguments": argumentsString,
					},
				},
			},
		})

		// 查找对应的工具函数
		function, exists := functionMapper[functionName]
		if !exists {
			return "", fmt.Errorf("unknown function: %s", functionName)
		}

		// 执行工具函数
		var arguments map[string]interface{}
		err = json.Unmarshal([]byte(argumentsString), &arguments)
		if err != nil {
			return "", fmt.Errorf("解析参数失败:%v", err)
		}
		fmt.Println(arguments)
		functionOutput, err := function(arguments)
		if err != nil {
			return "", fmt.Errorf("function execution failed: %v", err)
		}

		// 将工具响应加入消息历史
		messages = append(messages, map[string]interface{}{
			"role":         "tool",
			"content":      fmt.Sprintf("%v", functionOutput),
			"tool_call_id": firstToolCall.Get("id").String(),
		})

		// 输出函数执行结果
		fmt.Println("Function output:", functionOutput)
	}

}*/
