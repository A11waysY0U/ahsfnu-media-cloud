package main

import (
	"ahsfnu-media-cloud/internal/api"
	"ahsfnu-media-cloud/internal/config"
	"ahsfnu-media-cloud/internal/database"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	config.Init()

	// 初始化数据库
	database.Init()

	r := gin.Default()

	// 设置路由
	api.SetupRoutes(r)

	// 启动服务器
	log.Printf("Server starting on port %s", config.AppConfig.Server.Port)
	if err := r.Run(":" + config.AppConfig.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
