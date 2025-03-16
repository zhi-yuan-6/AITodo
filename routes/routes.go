package routes

import (
	"AITodo/controllers"
	"AITodo/middleware"
	"AITodo/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func SetupRoutes(router *gin.Engine) {
	// 全局应用 CORS 中间件
	router.Use(middleware.SetupCORS())

	// AI 相关路由（假设需要认证）
	ai := router.Group("/ai").Use(middleware.JWTAuth())
	{
		ai.POST("/assist", controllers.AIAssistant)
	}

	// 用户管理路由（无需认证）
	user := router.Group("/user")
	{
		user.POST("/login", controllers.LoginHandler)
		user.POST("/register", controllers.RegisterHandler)
		user.GET("/send_code", controllers.SendVerifyCodeHandler)
	}

	// 任务管理路由（需要认证）
	task := router.Group("/task").Use(middleware.JWTAuth())
	{
		task.GET("/", controllers.GetAllTasks)
		task.POST("/", controllers.CreateTask)
		task.PUT("/:id", controllers.UpdateTask)
		task.DELETE("/:id", controllers.DeleteTask)
	}

	analyticsGroup := router.Group("/analytics").Use(middleware.JWTAuth())
	{
		analyticsGroup.GET("/", controllers.GetCombinedDataHandler)
		analyticsGroup.GET("/trend", controllers.GetTrendHandler)
		analyticsGroup.GET("/category_distribution", controllers.GetCategoryDistributionHandler)
		analyticsGroup.GET("/heatmap", controllers.GetHeatmapHandler)
		analyticsGroup.POST("/ai_report", controllers.AIAnalytics)
	}
	// 认证相关路由
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/refresh", middleware.RefreshToken())
		authGroup.POST("/logout", middleware.JWTAuth(), logoutHandler)
	}
}

func logoutHandler(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
		return
	}

	parts := strings.SplitN(tokenString, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "令牌格式错误"})
		return
	}
	//token黑名单那
	if err := util.InvalidateToken(parts[1]); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注销失败:" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "成功注销"})
}
