// controllers/analytics_controller.go
package controllers

import (
	"AITodo/dto"
	"AITodo/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetTrendHandler 获取任务趋势数据
// @Summary 获取任务完成趋势
// @Description 获取堆积柱状图所需数据
// @Tags 数据分析
// @Produce json
// @Param start query string true "开始时间 (格式: 2006-01-02)"
// @Param end query string true "结束时间 (格式: 2006-01-02)"
// @Param interval query string true "时间间隔 (week/month/year)"
// @Security ApiKeyAuth
// @Success 200 {object} TrendResponse
// @Router /analytics/trend [get]
func GetTrendHandler(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	startStr := c.Query("start")
	endStr := c.Query("end")
	interval := c.Query("interval")

	// 解析时间
	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date format"})
		return
	}

	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date format"})
		return
	}

	// 获取趋势数据
	data, err := services.GetTrendData(userID, interval, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := formatTrendResponse(interval, data)
	c.JSON(http.StatusOK, response)
}

// GetCategoryDistributionHandler 获取任务分类分布
// @Summary 获取任务分类分布
// @Description 获取环形图所需数据
// @Tags 数据分析
// @Produce json
// @Param start query string true "开始时间 (格式: 2006-01-02)"
// @Param end query string true "结束时间 (格式: 2006-01-02)"
// @Param interval query string true "时间间隔 (week/month/year)"
// @Security ApiKeyAuth
// @Success 200 {object} CategoryDistributionResponse
// @Router /analytics/category_distribution [get]
func GetCategoryDistributionHandler(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	startStr := c.Query("start")
	endStr := c.Query("end")
	interval := c.Query("interval")

	// 解析时间
	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date format"})
		return
	}

	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date format"})
		return
	}

	// 获取分类分布数据
	data, err := services.GetCategoryDistribution(userID, interval, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.CategoryDistributionResponse{
		Categories: []dto.CategoryData{
			{
				Name:  "工作",
				Value: data.Work,
				Color: "#4CAF50",
			},
			{
				Name:  "学习",
				Value: data.Study,
				Color: "#2196F3",
			},
			{
				Name:  "生活",
				Value: data.Life,
				Color: "#FFC107",
			},
			{
				Name:  "健身",
				Value: data.Fitness,
				Color: "#F44336",
			},
			{
				Name:  "其他",
				Value: data.Other,
				Color: "#9E9E9E",
			},
		},
	}

	c.JSON(http.StatusOK, response)
}

// GetHeatmapHandler 获取用户活跃时段分析
// @Summary 获取用户活跃时段分析
// @Description 获取热力图所需数据
// @Tags 数据分析
// @Produce json
// @Param start query string true "开始时间 (格式: 2006-01-02)"
// @Param end query string true "结束时间 (格式: 2006-01-02)"
// @Param interval query string true "时间间隔 (week/month/year)"
// @Security ApiKeyAuth
// @Success 200 {object} HeatmapResponse
// @Router /analytics/heatmap [get]
// 时段分析 热力图
func GetHeatmapHandler(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	startStr := c.Query("start")
	endStr := c.Query("end")
	interval := c.Query("interval")

	// 解析时间
	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date format"})
		return
	}

	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date format"})
		return
	}

	// 获取活跃时段数据
	data, err := services.GetHeatmapData(userID, interval, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.HeatmapResponse{
		TimeSlots: data,
	}

	c.JSON(http.StatusOK, response)
}

// GetCombinedDataHandler 获取三个图表的数据
// @Summary 获取三个图表的数据
// @Description 获取趋势图、分类分布图和热力图所需数据
// @Tags 数据分析
// @Produce json
// @Param start query string true "开始时间 (格式: 2006-01-02)"
// @Param end query string true "结束时间 (格式: 2006-01-02)"
// @Param interval query string true "时间间隔 (week/month/year)"
// @Security ApiKeyAuth
// @Success 200 {object} CombinedResponse
// @Router /analytics/combined [get]
func GetCombinedDataHandler(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	startStr := c.Query("start")
	endStr := c.Query("end")
	interval := c.Query("interval")

	// 解析时间
	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date format"})
		return
	}

	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date format"})
		return
	}

	// 获取趋势图数据
	trendData, err := services.GetTrendData(userID, interval, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	trendResponse := formatTrendResponse(interval, trendData)

	// 获取分类分布图数据
	categoryData, err := services.GetCategoryDistribution(userID, interval, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	categoryResponse := dto.CategoryDistributionResponse{
		Categories: []dto.CategoryData{
			{
				Name:  "工作",
				Value: categoryData.Work,
				Color: "#4CAF50",
			},
			{
				Name:  "学习",
				Value: categoryData.Study,
				Color: "#2196F3",
			},
			{
				Name:  "生活",
				Value: categoryData.Life,
				Color: "#FFC107",
			},
			{
				Name:  "健身",
				Value: categoryData.Fitness,
				Color: "#F44336",
			},
			{
				Name:  "其他",
				Value: categoryData.Other,
				Color: "#9E9E9E",
			},
		},
	}

	// 获取热力图数据
	heatmapData, err := services.GetHeatmapData(userID, interval, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	heatmapResponse := dto.HeatmapResponse{
		TimeSlots: heatmapData,
	}

	// 合并数据并返回
	response := dto.CombinedResponse{
		Trend:                trendResponse,
		CategoryDistribution: categoryResponse,
		Heatmap:              heatmapResponse,
	}

	c.JSON(http.StatusOK, response)
}

func formatTrendResponse(interval string, data *dto.TrendData) dto.TrendResponse {
	return dto.TrendResponse{
		TimeRange: interval,
		Labels:    generateLabels(interval, data.StartDate, len(data.Completed)),
		Series: []dto.TrendSeries{
			{
				Name:  "Completed",
				Data:  data.Completed,
				Color: "#4CAF50",
			},
			{
				Name:  "Pending",
				Data:  data.Pending,
				Color: "#FFC107",
			},
			{
				Name:  "Failed",
				Data:  data.Failed,
				Color: "#F44336",
			},
		},
	}
}

func generateLabels(interval string, start time.Time, count int) []string {
	labels := make([]string, count)
	for i := 0; i < count; i++ {
		switch interval {
		case "week":
			labels[i] = start.AddDate(0, 0, i*7).Format("2006-01-02")
		case "month":
			labels[i] = start.AddDate(0, i, 0).Format("2006-01")
		case "year":
			labels[i] = start.AddDate(i, 0, 0).Format("2006")
		}
	}
	return labels
}
