package routes

import (
	"AITodo/controllers"
	"AITodo/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	//全局应用jwt中间件和cors中间件
	router.Use(middleware.SetupCORS(), middleware.JWTAuth())

	//AI相关路由
	ai := router.Group("/ai")
	{
		ai.POST("/assist", controllers.AIAssistant)
	}

	//用户管理路由
	user := router.Group("/user")
	{
		user.POST("/login", controllers.LoginHandler)
		user.POST("/register", controllers.RegisterHandler)
		user.GET("/send_code", controllers.SendVerifyCodeHandler)
	}

	//任务管理路由
	task := router.Group("/task")
	{
		task.GET("/", controllers.GetAllTasks)
		task.POST("/", controllers.CreateTask)
		task.PUT("/:id", controllers.UpdateTask)
		task.DELETE("/:id", controllers.DeleteTask)
	}
}

//添加注册限流中间件
/*func RegisterRateLimit() gin.HandlerFunc {
	return rateLimiter(5,time.Minute)
}*/
