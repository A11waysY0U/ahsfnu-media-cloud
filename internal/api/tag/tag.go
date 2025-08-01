package tag

import (
	"errors"
	"net/http"
	"strconv"

	"ahsfnu-media-cloud/internal/database"
	"ahsfnu-media-cloud/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TagService struct {
	db *gorm.DB
}

func NewTagService() *TagService {
	return &TagService{
		db: database.GetDB(),
	}
}

// GetTags 获取标签列表
func GetTags(c *gin.Context) {
	service := NewTagService()

	var tags []models.Tag
	err := service.db.Preload("Creator").Preload("MaterialTags.Material").Find(&tags).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取标签列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tags})
}

// CreateTag 创建标签
func CreateTag(c *gin.Context) {
	service := NewTagService()

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	var req struct {
		Name  string `json:"name" binding:"required"`
		Color string `json:"color"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查标签名是否已存在
	var existingTag models.Tag
	if err := service.db.Where("name = ?", req.Name).First(&existingTag).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标签名已存在"})
		return
	}

	// 设置默认颜色
	if req.Color == "" {
		req.Color = "#409EFF"
	}

	tag := &models.Tag{
		Name:      req.Name,
		Color:     req.Color,
		CreatedBy: userID.(uint),
	}

	if err := service.db.Create(tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建标签失败"})
		return
	}

	// 预加载创建者信息
	service.db.Preload("Creator").First(tag, tag.ID)

	c.JSON(http.StatusCreated, tag)
}

// UpdateTag 更新标签
func UpdateTag(c *gin.Context) {
	service := NewTagService()

	id := c.Param("id")
	tagID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的标签ID"})
		return
	}

	var tag models.Tag
	err = service.db.First(&tag, tagID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "标签不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取标签失败"})
		return
	}

	// 检查权限（只有创建者或管理员可以修改）
	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("role")
	if tag.CreatedBy != userID.(uint) && userRole.(string) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限修改此标签"})
		return
	}

	var updateData struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新字段
	updates := make(map[string]interface{})
	if updateData.Name != "" {
		updates["name"] = updateData.Name
	}
	if updateData.Color != "" {
		updates["color"] = updateData.Color
	}

	if err := service.db.Model(&tag).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新标签失败"})
		return
	}

	// 重新获取更新后的数据
	service.db.Preload("Creator").First(&tag, tagID)

	c.JSON(http.StatusOK, tag)
}

// DeleteTag 删除标签
func DeleteTag(c *gin.Context) {
	service := NewTagService()

	id := c.Param("id")
	tagID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的标签ID"})
		return
	}

	var tag models.Tag
	err = service.db.First(&tag, tagID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "标签不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取标签失败"})
		return
	}

	// 检查权限
	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("role")
	if tag.CreatedBy != userID.(uint) && userRole.(string) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限删除此标签"})
		return
	}

	// 删除标签（会级联删除关联的素材标签）
	if err := service.db.Delete(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除标签失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "标签删除成功"})
}

// AddTagToMaterial 为素材添加标签
func AddTagToMaterial(c *gin.Context) {
	service := NewTagService()

	tagID := c.Param("id")
	materialID := c.Param("materialId")

	tagIDUint, err := strconv.ParseUint(tagID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的标签ID"})
		return
	}

	materialIDUint, err := strconv.ParseUint(materialID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的素材ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 检查标签和素材是否存在
	var tag models.Tag
	if err := service.db.First(&tag, tagIDUint).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "标签不存在"})
		return
	}

	var material models.Material
	if err := service.db.First(&material, materialIDUint).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "素材不存在"})
		return
	}

	// 创建素材标签关联
	materialTag := &models.MaterialTag{
		MaterialID: uint(materialIDUint),
		TagID:      uint(tagIDUint),
		CreatedBy:  userID.(uint),
	}

	if err := service.db.Create(materialTag).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "该标签已存在于该素材"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加标签失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "标签添加成功"})
}

// RemoveTagFromMaterial 从素材移除标签
func RemoveTagFromMaterial(c *gin.Context) {
	service := NewTagService()

	tagID := c.Param("id")
	materialID := c.Param("materialId")

	tagIDUint, err := strconv.ParseUint(tagID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的标签ID"})
		return
	}

	materialIDUint, err := strconv.ParseUint(materialID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的素材ID"})
		return
	}

	// 删除素材标签关联
	if err := service.db.Where("material_id = ? AND tag_id = ?", materialIDUint, tagIDUint).Delete(&models.MaterialTag{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "移除标签失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "标签移除成功"})
}
