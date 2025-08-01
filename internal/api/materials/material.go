package materials

import (
	"ahsfnu-media-cloud/internal/database"
	"ahsfnu-media-cloud/internal/models"
	"ahsfnu-media-cloud/internal/services"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MaterialService struct {
	db            *gorm.DB
	uploadService *services.UploadService
}

type MaterialQueryBuilder struct {
	query *gorm.DB
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

func NewMaterialQueryBuilder(db *gorm.DB) *MaterialQueryBuilder {
	return &MaterialQueryBuilder{
		query: db.Model(&models.Material{}).Preload("Uploader").Preload("MaterialTags.Tag"),
	}
}

func (qb *MaterialQueryBuilder) WithWorkflow(workflowID string) *MaterialQueryBuilder {
	if workflowID != "" {
		qb.query = qb.query.Where("workflow_id = ?", workflowID)
	}
	return qb
}

func (qb *MaterialQueryBuilder) WithFileType(fileType string) *MaterialQueryBuilder {
	if fileType != "" {
		qb.query = qb.query.Where("file_type = ?", fileType)
	}
	return qb
}

func (qb *MaterialQueryBuilder) WithKeyword(keyword string) *MaterialQueryBuilder {
	if keyword != "" {
		qb.query = qb.query.Where("original_filename ILIKE ?", "%"+keyword+"%")
	}
	return qb
}

func (qb *MaterialQueryBuilder) WithTags(tagsParam string) *MaterialQueryBuilder {
	if tagsParam != "" {
		tagIDs := []uint{}
		for _, idStr := range SplitAndTrim(tagsParam, ",") {
			if id, err := strconv.ParseUint(idStr, 10, 32); err == nil {
				tagIDs = append(tagIDs, uint(id))
			}
		}
		if len(tagIDs) > 0 {
			qb.query = qb.query.Joins("JOIN material_tags ON materials.id = material_tags.material_id").
				Where("material_tags.tag_id IN ?", tagIDs).Group("materials.id")
		}
	}
	return qb
}

func (qb *MaterialQueryBuilder) WithPublic() *MaterialQueryBuilder {
	qb.query = qb.query.Where("is_public = ?", true)
	return qb
}

func (qb *MaterialQueryBuilder) Build() *gorm.DB {
	return qb.query
}

func paginatedResponse(c *gin.Context, data interface{}, page, pageSize int, total int64) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
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

// GetMaterial 获取素材详情
func GetMaterial(c *gin.Context) {
	service := GetMaterialService()

	materialID, valid := validateMaterialID(c)
	if !valid {
		return
	}

	var material models.Material
	err := service.db.Preload("Uploader").Preload("MaterialTags.Tag").First(&material, materialID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			errorResponse(c, http.StatusNotFound, "素材不存在")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "获取素材详情失败")
		return
	}

	// 添加文件URL
	material.FilePath = service.uploadService.GetFileURL(&material)

	successResponse(c, material)
}

// DeleteMaterial 删除素材
func DeleteMaterial(c *gin.Context) {
	service := GetMaterialService()

	id := c.Param("id")
	materialID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "无效的素材ID")
		return
	}

	var material models.Material
	err = service.db.First(&material, materialID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			errorResponse(c, http.StatusNotFound, "素材不存在")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "获取素材失败")
		return
	}

	// 检查权限
	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("role")
	if material.UploadedBy != userID.(uint) && userRole.(string) != "admin" {
		errorResponse(c, http.StatusForbidden, "没有权限删除此素材")
		return
	}

	// 删除文件
	if err := service.uploadService.DeleteFile(&material); err != nil {
		errorResponse(c, http.StatusInternalServerError, "删除文件失败")
		return
	}

	// 删除数据库记录
	if err := service.db.Delete(&material).Error; err != nil {
		errorResponse(c, http.StatusInternalServerError, "删除素材记录失败")
		return
	}

	successResponse(c, gin.H{"message": "素材删除成功"})
}

// GetMaterials 获取素材列表
func SearchMaterials(c *gin.Context) {
	service := GetMaterialService()

	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	workflowID := c.Query("workflow_id")
	fileType := c.Query("file_type")
	keyword := c.Query("keyword")
	tagsParam := c.Query("tags")

	// 使用查询构建器
	queryBuilder := NewMaterialQueryBuilder(service.db)
	query := queryBuilder.
		WithWorkflow(workflowID).
		WithFileType(fileType).
		WithKeyword(keyword).
		WithTags(tagsParam).
		Build()

	// 分页
	offset := (page - 1) * pageSize
	var materials []models.Material
	var total int64

	query.Count(&total)
	err := query.Offset(offset).Limit(pageSize).Order("upload_time DESC").Find(&materials).Error
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "获取素材列表失败")
		return
	}

	// 添加文件URL
	for i := range materials {
		materials[i].FilePath = service.uploadService.GetFileURL(&materials[i])
	}

	paginatedResponse(c, materials, page, pageSize, total)
}

// SplitAndTrim 工具函数：分割字符串并去除空格
func SplitAndTrim(s, sep string) []string {
	res := []string{}
	for _, part := range strings.Split(s, sep) {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			res = append(res, trimmed)
		}
	}
	return res
}

func ToggleStar(c *gin.Context) {
	service := GetMaterialService()

	id := c.Param("id")
	materialID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "无效的素材ID")
		return
	}

	var material models.Material
	err = service.db.First(&material, materialID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			errorResponse(c, http.StatusNotFound, "素材不存在")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "获取素材失败")
		return
	}

	// 切换星标状态
	material.IsStarred = !material.IsStarred

	if err := service.db.Save(&material).Error; err != nil {
		errorResponse(c, http.StatusInternalServerError, "更新星标状态失败")
		return
	}

	successResponse(c, gin.H{
		"is_starred": material.IsStarred,
		"message":    "星标状态更新成功",
	})
}
