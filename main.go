package main

import (
	"AITodo/config"
	"AITodo/db"
	"AITodo/globals"
	"AITodo/models"
	"AITodo/routes"
	"AITodo/util"
	"github.com/gin-gonic/gin"
	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
)

func main() {
	// 初始化协程池（最大100个协程）
	globals.TaskPool, _ = ants.NewPool(100)
	defer globals.TaskPool.Release() // 程序退出时释放协程池

	// 加载配置
	err := config.LoadConfig("./config/config.yaml")
	if err != nil {
		logrus.Fatal("failed to load config: %w", err)
	}

	// 初始化Redis
	if err = util.Initialize(config.Cfg.Redis); err != nil {
		logrus.Fatalf("Redis initialization failed: %v", err)
	}

	err = db.InitMysql()
	if err != nil {
		logrus.Fatal(err)
	}

	err = db.DB.AutoMigrate(&models.Task{}, &models.User{})
	if err != nil {
		logrus.Fatalf("数据库迁移失败: %v", err) // 终止程序
	}

	util.JwtSecret, util.PublicKey, err = util.ReadPEM("private_key.pem")
	if err != nil {
		logrus.Fatalf("私钥读取失败:%v", err)
	}

	//初始化gin
	router := gin.Default()

	//设置路由
	routes.SetupRoutes(router)

	//启动服务
	router.Run(":8080")

}
