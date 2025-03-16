package dto

import (
	"time"
)

// 完成趋势 堆积柱状图
type TrendResponse struct {
	TimeRange string        `json:"time_range"` // weekly/monthly/yearly
	Labels    []string      `json:"labels"`     // 时间轴标签
	Series    []TrendSeries `json:"series"`     // 数据系列
}
type TrendSeries struct {
	Name  string `json:"name"`  // 状态名称
	Data  []int  `json:"data"`  // 数据值
	Color string `json:"color"` // 颜色代码
}
type TrendData struct {
	StartDate time.Time
	Pending   []int
	Progress  []int
	Completed []int
	Failed    []int
}

// 任务分类 环形图
type CategoryDistributionResponse struct {
	Categories []CategoryData `json:"categories"`
}
type CategoryData struct {
	Name  string `json:"name"`  // 类别名称
	Value int    `json:"value"` // 任务数量
	Color string `json:"color"` // 颜色代码
}
type CategoryDistribution struct {
	Work    int
	Study   int
	Life    int
	Fitness int
	Other   int
}

// 时段分析 热力图
type HeatmapResponse struct {
	TimeSlots []TimeSlot `json:"time_slots"`
}
type TimeSlot struct {
	Time  string `json:"time"`  // 时间段 (格式: HH:mm)
	Count int    `json:"count"` // 任务数量
}

// 总数据
type CombinedResponse struct {
	Trend                TrendResponse                `json:"trend"`
	CategoryDistribution CategoryDistributionResponse `json:"category_distribution"`
	Heatmap              HeatmapResponse              `json:"heatmap"`
}
