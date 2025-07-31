package materials

import (
	"ahsfnu-media-cloud/internal/database"
	"ahsfnu-media-cloud/internal/models"
	"ahsfnu-media-cloud/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MaterialService struct {
	db            *gorm.DB
	uploadService *services.UploadService
}

// 单例模式
var materialService *MaterialService

func GetMaterialService() *MaterialService {
	if materialService == nil {
		materialService = &MaterialService{
			db:            database.GetDB(),
			uploadService: services.NewUploadService(),
		}
	}
	return materialService
}

// 参数验证工具函数
func validateMaterialID(c *gin.Context) (uint, bool) {
	id := c.Param("id")
	materialID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "无效的素材ID")
		return 0, false
	}
	return uint(materialID), true
}

// 权限检查工具函数
func checkMaterialPermission(c *gin.Context, material *models.Material) bool {
	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("role")

	if material.UploadedBy != userID.(uint) && userRole.(string) != "admin" {
		errorResponse(c, http.StatusForbidden, "没有权限操作此素材")
		return false
	}
	return true
}

// 获取素材工具函数
func getMaterialByID(service *MaterialService, materialID uint) (*models.Material, bool) {
	var material models.Material
	err := service.db.First(&material, materialID).Error
	if err != nil {
		return nil, false
	}
	return &material, true
}

// 统一错误响应
func errorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
}

// 统一成功响应
func successResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"data": data})
}

// UploadMaterial 上传素材
func UploadMaterial(c *gin.Context) {
	service := GetMaterialService()

	// 获取当前用户ID（中间件已确保用户已认证）
	userID, _ := c.Get("user_id")

	// 获取工作流ID（可选）
	var workflowID *uint
	if workflowIDStr := c.PostForm("workflow_id"); workflowIDStr != "" {
		if id, err := strconv.ParseUint(workflowIDStr, 10, 32); err == nil {
			workflowIDUint := uint(id)
			workflowID = &workflowIDUint
		}
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "请选择要上传的文件")
		return
	}

	// 使用数据库事务确保一致性
	tx := service.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 上传文件
	material, err := service.uploadService.UploadFile(file, userID.(uint), workflowID)
	if err != nil {
		tx.Rollback()
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// 保存到数据库
	if err := tx.Create(material).Error; err != nil {
		tx.Rollback()
		errorResponse(c, http.StatusInternalServerError, "保存素材记录失败")
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		errorResponse(c, http.StatusInternalServerError, "提交事务失败")
		return
	}

	// 预加载关联数据
	service.db.Preload("Uploader").First(material, material.ID)

	// 添加文件URL
	material.FilePath = service.uploadService.GetFileURL(material)

	successResponse(c, material)
}

// UpdateMaterial 更新素材
func UpdateMaterial(c *gin.Context) {
	service := GetMaterialService()

	materialID, valid := validateMaterialID(c)
	if !valid {
		return
	}

	material, found := getMaterialByID(service, materialID)
	if !found {
		errorResponse(c, http.StatusNotFound, "素材不存在")
		return
	}

	// 检查权限
	if !checkMaterialPermission(c, material) {
		return
	}

	// 更新字段
	var updateData struct {
		OriginalFilename string `json:"original_filename"`
		IsStarred        *bool  `json:"is_starred"`
		IsPublic         *bool  `json:"is_public"`
		WorkflowID       *uint  `json:"workflow_id"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// 更新数据库
	updates := make(map[string]interface{})
	if updateData.OriginalFilename != "" {
		updates["original_filename"] = updateData.OriginalFilename
	}
	if updateData.IsStarred != nil {
		updates["is_starred"] = *updateData.IsStarred
	}
	if updateData.IsPublic != nil {
		updates["is_public"] = *updateData.IsPublic
	}
	if updateData.WorkflowID != nil || (updateData.WorkflowID == nil && updateData.WorkflowID != material.WorkflowID) {
		updates["workflow_id"] = updateData.WorkflowID
	}

	if err := service.db.Model(&material).Updates(updates).Error; err != nil {
		errorResponse(c, http.StatusInternalServerError, "更新素材失败")
		return
	}

	// 重新获取更新后的数据
	service.db.Preload("Uploader").First(material, materialID)
	material.FilePath = service.uploadService.GetFileURL(material)

	successResponse(c, material)
}

func GetMaterialDetails(c *gin.Context) {
	service := GetMaterialService()

	materialID, valid := validateMaterialID(c)
	if !valid {
		return
	}

	material, found := getMaterialByID(service, materialID)
	if !found {
		errorResponse(c, http.StatusNotFound, "素材不存在")
		return
	}

	// 检查权限
	if !checkMaterialPermission(c, material) {
		return
	}

	// 添加文件URL
	material.FilePath = service.uploadService.GetFileURL(material)

	successResponse(c, material)
}
