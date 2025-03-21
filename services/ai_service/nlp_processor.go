package ai_service

import (
	"AITodo/db"
	"AITodo/models"
	"errors"
	"fmt"
	"github.com/agnivade/levenshtein"
	"gorm.io/gorm"
	"log"
	"math"
	"strings"
	"time"
)

// 添加评分权重常量
const (
	TitleSimilarityWeight  = 0.4
	KeywordCoverageWeight  = 0.2
	StatusMatchWeight      = 0.2
	DateProximityWeight    = 0.2
	MinAcceptableScore     = 0.6 // 最低可接受匹配分数
	DateProximityThreshold = 7   // 日期邻近阈值（天）
)

func searchTask(input map[string]interface{}) (string, error) {
	taskInterface, err := parseTask(input)
	if err != nil {
		return "", fmt.Errorf("input解析失败: %v", err)
	}
	task, ok := taskInterface.(models.Task)
	if !ok {
		return "", fmt.Errorf("任务数据解析失败，类型转换失败")
	}

	query := db.DB.Model(&models.Task{})
	// 动态构建查询条件
	buildQueryConditions(query, task)

	var tasks []models.Task
	if err := query.Order("created_at DESC").Find(&tasks).Error; err != nil {
		return "", fmt.Errorf("数据库查询失败: %v", err)
	}

	//添加对tasks的匹配
	if len(tasks) > 1 {
		return handleMultipleMatches(tasks, task)
	}

	return handleSearchResults(tasks)
}

// 构建动态查询条件
func buildQueryConditions(query *gorm.DB, task models.Task) *gorm.DB {
	if !task.StartDate.IsZero() {
		// 时间范围查询（±1天）
		start := task.StartDate.Add(-24 * time.Hour)
		end := task.StartDate.Add(24 * time.Hour)
		query = query.Where("start_date BETWEEN ? AND ?", start, end)
	}
	if !task.DueDate.IsZero() {
		// 精确匹配截止日期
		query = query.Where("due_date <= ?", task.DueDate)
	}

	if task.Status != "" {
		query = query.Where("status = ?", task.Status)
	}

	return query
}

// 新增多匹配处理函数
func handleMultipleMatches(tasks []models.Task, candidate models.Task) (string, error) {
	scores := make(map[uint]float64)

	for _, t := range tasks {
		score := calculateMatchScore(t, candidate)
		scores[t.ID] = score
	}

	// 寻找最佳匹配
	maxScore := 0.0
	var bestTask models.Task
	for _, t := range tasks {
		if scores[t.ID] > maxScore {
			maxScore = scores[t.ID]
			bestTask = t
		}
	}

	// 验证匹配质量
	if maxScore < MinAcceptableScore {
		return "", fmt.Errorf("找到%d个可能匹配，但均未达到匹配阈值（%.2f/%f）",
			len(tasks), maxScore, MinAcceptableScore)
	}

	return fmt.Sprintf("%d", bestTask.ID), nil
}

// 综合评分计算函数
func calculateMatchScore(task models.Task, candidate models.Task) float64 {
	score := 0.0

	// 1. 标题相似度
	score += CalculateStringSimilarity(task.Title, candidate.Title) * TitleSimilarityWeight

	// 2. 描述关键词覆盖率
	score += KeywordCoverage(task.Description, candidate.Description) * KeywordCoverageWeight

	// 3. 状态匹配（新增）
	if task.Status == candidate.Status {
		score += StatusMatchWeight
	}

	// 4. 日期邻近性（新增）
	if !candidate.StartDate.IsZero() && !task.StartDate.IsZero() {
		daysDiff := math.Abs(candidate.StartDate.Sub(task.StartDate).Hours() / 24)
		if daysDiff <= DateProximityThreshold {
			score += (1 - daysDiff/DateProximityThreshold) * DateProximityWeight
		}
	}

	return score
}

// 字符串相似度计算（保持原有实现）
func CalculateStringSimilarity(a, b string) float64 {
	distance := levenshtein.ComputeDistance(a, b)
	maxLen := math.Max(float64(len(a)), float64(len(b)))
	return 1 - float64(distance)/maxLen
}

// 关键词覆盖率计算（保持原有实现）
func KeywordCoverage(target, query string) float64 {
	queryWords := strings.Fields(query)
	if len(queryWords) == 0 {
		return 0
	}

	count := 0
	for _, word := range queryWords {
		if strings.Contains(strings.ToLower(target), strings.ToLower(word)) {
			count++
		}
	}
	return float64(count) / float64(len(queryWords))
}

// 修改后的结果处理函数
func handleSearchResults(tasks []models.Task) (string, error) {
	switch len(tasks) {
	case 0:
		return "", errors.New("未找到匹配任务")
	case 1:
		return fmt.Sprintf("%d", tasks[0].ID), nil
	default:
		// 现在由handleMultipleMatches处理多结果情况
		return "", errors.New("匹配结果处理异常")
	}
}

func parseTask(args map[string]interface{}) (interface{}, error) {
	// 提取 task 字段
	task, ok := args["task"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("task字段缺失或格式错误")
	}

	// 验证 title 和 due_date 是否存在
	if err, ok := task["title"].(string); !ok {
		return nil, fmt.Errorf("task字段确实或者格式错误:%w", err)
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
		Title:       task["title"].(string),
		Description: parseString(task["description"]),
		Status:      parseStatus(task["status"]),
		StartDate:   startDate,
		DueDate:     dueDate,
	}

	return taskModel, nil
}
