package api

import (
	"ahsfnu-media-cloud/internal/api/auth"
	"ahsfnu-media-cloud/internal/api/materials"
	"ahsfnu-media-cloud/internal/api/tag"
	"ahsfnu-media-cloud/internal/api/workflow"
	"ahsfnu-media-cloud/internal/config"
	"ahsfnu-media-cloud/internal/database"
	"ahsfnu-media-cloud/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine) {
	// 添加中间件
	r.Use(middleware.CORSMiddleware())

	// 静态文件服务 - 提供文件的访问
	r.Static("/uploads", config.AppConfig.Upload.UploadPath)

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 认证相关路由
		authGroup := v1.Group("/auth")
		{
			authGroup.GET("/captcha", auth.GetCaptcha)
			authGroup.POST("/verify-captcha", auth.VerifyCaptcha)
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
			materialGroup.GET("/:id", materials.GetMaterial)
			materialGroup.DELETE("/:id", materials.DeleteMaterial)
			materialGroup.GET("", materials.SearchMaterials)
		}
		protected.POST("/invite_codes", auth.GenerateInviteCodes)
		protected.GET("/invite_codes", auth.ListInviteCodes)
		// 用户相关
		protected.GET("/profile", auth.GetProfile)
		protected.PUT("/profile", auth.UpdateProfile)
		protected.PUT("/profile/password", auth.ChangePassword)
		protected.GET("/users", auth.GetUsers)
		protected.PUT("/users/:id/role", auth.UpdateUserRole) // 管理员修改用户权限
		protected.DELETE("/users/:id", auth.DeleteUser)       // 管理员删除用户

		// 标签相关路由
		tagGroup := protected.Group("/tags")
		{
			tagGroup.GET("", tag.GetTags)
			tagGroup.POST("", tag.CreateTag)
			tagGroup.PUT("/:id", tag.UpdateTag)
			tagGroup.DELETE("/:id", tag.DeleteTag)
			tagGroup.POST("/:id/materials/:materialId", tag.AddTagToMaterial)
			tagGroup.DELETE("/:id/materials/:materialId", tag.RemoveTagFromMaterial)
		}

		// 工作流相关路由
		workflowGroup := protected.Group("/workflows")
		{
			workflowGroup.GET("", workflow.GetWorkflows)
			workflowGroup.POST("", workflow.CreateWorkflow)
			workflowGroup.GET("/:id", workflow.GetWorkflow)
			workflowGroup.PUT("/:id", workflow.UpdateWorkflow)
			workflowGroup.DELETE("/:id", workflow.DeleteWorkflow)
			workflowGroup.POST("/:id/members", workflow.AddWorkflowMember)
			workflowGroup.DELETE("/:id/members/:userId", workflow.RemoveWorkflowMember)
		}

	}

}
