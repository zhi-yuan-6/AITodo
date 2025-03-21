// Package services/analytics_service.go
package services

import (
	"AITodo/db"
	"AITodo/dto"
	"fmt"
	"time"
)

// GetTrendData 获取趋势数据
func GetTrendData(userID uint, interval string, start, end time.Time) (*dto.TrendData, error) {
	// 执行数据库查询
	rawData, err := fetchTrendDataFromDB(userID, interval, start, end)
	if err != nil {
		return nil, err
	}

	// 初始化结果集
	periods := calculatePeriods(start, end, interval)
	result := &dto.TrendData{
		StartDate: start,
		Completed: make([]int, periods),
		Pending:   make([]int, periods),
		Progress:  make([]int, periods),
		Failed:    make([]int, periods),
	}

	// 填充数据
	for _, item := range rawData {
		index := calculatePeriodIndex(start, item.PeriodStart, interval)
		if index < 0 || index >= periods {
			continue
		}

		switch item.Status {
		case "completed":
			result.Completed[index] += item.Count
		case "pending":
			result.Pending[index] += item.Count
		case "failed":
			result.Failed[index] += item.Count
		}
	}

	return result, nil
}

// fetchTrendDataFromDB 从数据库中获取趋势数据
func fetchTrendDataFromDB(userID uint, interval string, start, end time.Time) ([]struct {
	PeriodStart time.Time
	Status      string
	Count       int
}, error) {
	var result []struct {
		PeriodStart string
		Status      string
		Count       int
	}

	var query string

	// 使用 DATE_FORMAT 确保返回 YYYY-MM-DD 格式的字符串
	switch interval {
	case "week":
		query = `
            SELECT 
                DATE_FORMAT(DATE_SUB(created_at, INTERVAL WEEKDAY(created_at) DAY), '%Y-%m-%d') AS period_start,
                status,
                COUNT(*) AS count
            FROM tasks
            WHERE 
                user_id = ? AND
                created_at BETWEEN ? AND ?
            GROUP BY period_start, status
            ORDER BY period_start`
	case "month":
		query = `
            SELECT 
                DATE_FORMAT(created_at, '%Y-%m-01') AS period_start,
                status,
                COUNT(*) AS count
            FROM tasks
            WHERE 
                user_id = ? AND
                created_at BETWEEN ? AND ?
            GROUP BY period_start, status
            ORDER BY period_start`
	case "year":
		query = `
            SELECT 
                DATE_FORMAT(created_at, '%Y-01-01') AS period_start,
                status,
                COUNT(*) AS count
            FROM tasks
            WHERE 
                user_id = ? AND
                created_at BETWEEN ? AND ?
            GROUP BY period_start, status
            ORDER BY period_start`
	default: // daily
		query = `
            SELECT 
                DATE_FORMAT(created_at, '%Y-%m-%d') AS period_start,
                status,
                COUNT(*) AS count
            FROM tasks
            WHERE 
                user_id = ? AND
                created_at BETWEEN ? AND ?
            GROUP BY period_start, status
            ORDER BY period_start`
	}

	// 执行查询
	err := db.DB.Raw(query, userID, start, end).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	// 将字符串解析为 time.Time
	parsedResult := make([]struct {
		PeriodStart time.Time
		Status      string
		Count       int
	}, len(result))

	for i, item := range result {
		parsedTime, err := time.Parse("2006-01-02", item.PeriodStart)
		if err != nil {
			return nil, fmt.Errorf("failed to parse period_start '%s': %v", item.PeriodStart, err)
		}
		parsedResult[i] = struct {
			PeriodStart time.Time
			Status      string
			Count       int
		}{
			PeriodStart: parsedTime,
			Status:      item.Status,
			Count:       item.Count,
		}
	}

	return parsedResult, nil
}

// 计算时间段数量
func calculatePeriods(start, end time.Time, interval string) int {
	switch interval {
	case "week":
		return int(end.Sub(start).Hours()/(24*7)) + 1
	case "month":
		return (end.Year()-start.Year())*12 + int(end.Month()-start.Month()) + 1
	case "year":
		return end.Year() - start.Year() + 1
	default:
		return int(end.Sub(start).Hours()/24) + 1
	}
}

// 计算时间段索引
func calculatePeriodIndex(start, current time.Time, interval string) int {
	switch interval {
	case "week":
		weeks := int(current.Sub(start).Hours() / (24 * 7))
		if weeks < 0 {
			return -1
		}
		return weeks
	case "month":
		months := (current.Year()-start.Year())*12 + int(current.Month()-start.Month())
		return months
	case "year":
		return current.Year() - start.Year()
	default:
		return -1
	}
}

func GetCategoryDistribution(userID uint, interval string, start, end time.Time) (*dto.CategoryDistribution, error) {
	// 执行数据库查询
	rawData, err := fetchCategoryDataFromDB(userID, interval, start, end)
	if err != nil {
		return nil, err
	}

	// 初始化结果集
	result := &dto.CategoryDistribution{
		Work:    0,
		Study:   0,
		Life:    0,
		Fitness: 0,
		Other:   0,
	}

	// 填充数据
	for _, item := range rawData {
		switch item.Category {
		case "工作":
			result.Work += item.Count
		case "学习":
			result.Study += item.Count
		case "生活":
			result.Life += item.Count
		case "健身":
			result.Fitness += item.Count
		case "其他":
			result.Other += item.Count
		}
	}

	return result, nil
}

// 数据库查询实现
func fetchCategoryDataFromDB(userID uint, interval string, start, end time.Time) ([]struct {
	Category string
	Count    int
}, error) {
	var result []struct {
		Category string
		Count    int
	}

	var query string

	switch interval {
	case "week":
		query = `
			SELECT 
				category,
				COUNT(*) AS count
			FROM tasks
			WHERE 
				user_id = ? AND
				created_at BETWEEN ? AND ?
			GROUP BY category`
	case "month":
		query = `
			SELECT 
				category,
				COUNT(*) AS count
			FROM tasks
			WHERE 
				user_id = ? AND
				created_at BETWEEN ? AND ?
			GROUP BY category`
	case "year":
		query = `
			SELECT 
				category,
				COUNT(*) AS count
			FROM tasks
			WHERE 
				user_id = ? AND
				created_at BETWEEN ? AND ?
			GROUP BY category`
	default:
		query = `
			SELECT 
				category,
				COUNT(*) AS count
			FROM tasks
			WHERE 
				user_id = ? AND
				created_at BETWEEN ? AND ?
			GROUP BY category`
	}

	err := db.DB.Raw(query, userID, start, end).Scan(&result).Error
	return result, err
}

func GetHeatmapData(userID uint, interval string, start, end time.Time) ([]dto.TimeSlot, error) {
	// 执行数据库查询
	rawData, err := fetchHeatmapDataFromDB(userID, interval, start, end)
	if err != nil {
		return nil, err
	}

	return rawData, nil
}

// 数据库查询实现
func fetchHeatmapDataFromDB(userID uint, interval string, start, end time.Time) ([]dto.TimeSlot, error) {
	var result []dto.TimeSlot

	var query string

	switch interval {
	case "week":
		query = `
			SELECT 
				DATE_FORMAT(created_at, '%H:%i') AS time,
				COUNT(*) AS count
			FROM tasks
			WHERE 
				user_id = ? AND
				created_at BETWEEN ? AND ?
			GROUP BY time
			ORDER BY time`
	case "month":
		query = `
			SELECT 
				DATE_FORMAT(created_at, '%H:%i') AS time,
				COUNT(*) AS count
			FROM tasks
			WHERE 
				user_id = ? AND
				created_at BETWEEN ? AND ?
			GROUP BY time
			ORDER BY time`
	case "year":
		query = `
			SELECT 
				DATE_FORMAT(created_at, '%H:%i') AS time,
				COUNT(*) AS count
			FROM tasks
			WHERE 
				user_id = ? AND
				created_at BETWEEN ? AND ?
			GROUP BY time
			ORDER BY time`
	default:
		query = `
			SELECT 
				DATE_FORMAT(created_at, '%H:%i') AS time,
				COUNT(*) AS count
			FROM tasks
			WHERE 
				user_id = ? AND
				created_at BETWEEN ? AND ?
			GROUP BY time
			ORDER BY time`
	}

	err := db.DB.Raw(query, userID, start, end).Scan(&result).Error
	return result, err
}
