package ai_service

import (
	"AITodo/models"
	"AITodo/services"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"log"
	"strconv"
	"time"
)

// 统一错误定义
var (
	ErrInvalidID      = fmt.Errorf("无效的ID格式")
	ErrMissingTitle   = fmt.Errorf("title为必填字段")
	ErrMissingDueDate = fmt.Errorf("due_date为必填字段")
	ErrMissStartDate  = fmt.Errorf("start_date为必填字段")
)

// CreateTask 适配器
func adaptCreateTask(id uint, args map[string]interface{}) (interface{}, error) {
	// 参数验证
	if err := validateCreateArgs(args); err != nil {
		return nil, err
	}

	// 提取 task 字段
	task, ok := args["task"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("task字段缺失或格式错误")
	}
	startDate, err := parseTime(task["start_date"])
	if err != nil {
		log.Printf("任务 %s 的 start_date 解析失败: %v", task["title"], err)
	}
	dueDate, err := parseTime(task["due_date"])
	if err != nil {
		log.Printf("任务 %s 的 due_date 解析失败: %v", task["title"], err)
	}

	// 构建 Task 对象
	taskModel := models.Task{
		UserID:      id,
		Title:       task["title"].(string),
		Category:    task["category"].(string),
		Location:    parseString(task["location"]),
		Description: parseString(task["description"]),
		Status:      parseStatus(task["status"]),
		StartDate:   startDate,
		DueDate:     dueDate,
	}

	err = services.CreateTask(taskModel)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// UpdateTask 适配器
func adaptUpdateTask(args map[string]interface{}) (interface{}, error) {
	// 参数验证
	if err := validateUpdateArgs(args); err != nil {
		return nil, err
	}

	id := args["id"].(uint)

	// 提取 task 字段
	req, ok := args["task"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("task字段缺失或格式错误")
	}

	var task models.Task
	if err := mapstructure.Decode(req, &task); err != nil {
		return nil, fmt.Errorf("参数转换失败: %w", err)
	}

	updatedTask, err := services.UpdateTask(id, task)
	if err != nil {
		return nil, fmt.Errorf("更新任务失败: %v", err)
	}
	return updatedTask, nil

}

// DeleteTask 适配器
func adaptDeleteTask(args map[string]interface{}) (interface{}, error) {
	id, err := parseUint(args["id"])
	if err != nil {
		return nil, ErrInvalidID
	}
	// 添加存在性检查
	if _, err := models.GetTaskById(id); err != nil {
		return nil, fmt.Errorf("任务不存在")
	}

	err = services.DeleteTask(id)
	return nil, err
}

// 参数验证函数
func validateCreateArgs(args map[string]interface{}) error {
	task, ok := args["task"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("task字段缺失或格式错误")
	}

	// 验证 title 和 due_date 是否存在
	if _, ok := task["title"].(string); !ok {
		return ErrMissingTitle
	}
	if _, ok := task["start_date"].(string); !ok {
		return ErrMissStartDate
	}
	if _, ok := task["due_date"].(string); !ok {
		return ErrMissingDueDate
	}
	return nil
}

func validateUpdateArgs(args map[string]interface{}) error {
	task, ok := args["task"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("task字段缺失或格式错误")
	}

	// 验证 id 是否存在
	if _, ok := task["id"].(float64); !ok {
		return ErrInvalidID
	}
	return validateCreateArgs(args)
}

// 类型转换辅助函数
func parseUint(value interface{}) (uint, error) {
	switch v := value.(type) {
	case string:
		id, err := strconv.ParseUint(v, 10, 64)
		return uint(id), err
	case float64:
		return uint(v), nil
	case int:
		return uint(v), nil
	default:
		return 0, ErrInvalidID
	}
}

func parseString(value interface{}) string {
	if s, ok := value.(string); ok {
		return s
	}
	return ""
}

func parseStatus(value interface{}) string {
	status := "pending"
	if s, ok := value.(string); ok {
		switch s {
		case "pending", "in_progress", "completed":
			status = s
		}
	}
	return status
}

func parseTime(value interface{}) (time.Time, error) {
	// 使用正确的格式字符串
	const layout = "2006-01-02 15:04:05"

	// 检查 value 是否为字符串
	strValue, ok := value.(string)
	if !ok {
		return time.Time{}, fmt.Errorf("value is not a string")
	}

	if strValue == "" {
		return time.Time{}, nil
	}

	// 解析时间字符串
	parsedTime, err := time.ParseInLocation(layout, strValue, time.Local)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time: %w", err)
	}

	return parsedTime, nil
}
