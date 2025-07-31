package main

import (
	"ahsfnu-media-cloud/internal/api"
	"ahsfnu-media-cloud/internal/config"
	"ahsfnu-media-cloud/internal/database"
	"ahsfnu-media-cloud/internal/models"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	config.Init()

	// 初始化数据库
	database.Init()

	db := database.GetDB()
	err := db.AutoMigrate(
		&models.User{},
		&models.Material{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	r := gin.Default()

	// 设置路由
	api.SetupRoutes(r)

	// 启动服务器
	log.Printf("Server starting on port %s", config.AppConfig.Server.Port)
	if err := r.Run(":" + config.AppConfig.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
