package api

import (
	"ahsfnu-media-cloud/internal/api/auth"
	"ahsfnu-media-cloud/internal/api/materials"
	"ahsfnu-media-cloud/internal/config"
	"ahsfnu-media-cloud/internal/database"
	"ahsfnu-media-cloud/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine) {
	// 添加中间件
	r.Use(middleware.CORSMiddleware())

	// 静态文件服务 - 提供上传文件的访问
	r.Static("/uploads", config.AppConfig.Upload.UploadPath)

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 认证相关路由
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/login", auth.Login)
			authGroup.POST("/register", auth.Register)
		}
	}
	protected := v1.Group("/")
	protected.Use(middleware.AuthMiddleware(database.GetDB()))
	{

		materialGroup := protected.Group("/materials")
		{
			materialGroup.POST("", materials.UploadMaterial)
			materialGroup.PUT("/:id", materials.UpdateMaterial)
		}
	}
}
