package workflow

import (
	"ahsfnu-media-cloud/internal/database"
	"ahsfnu-media-cloud/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取工作流列表
func GetWorkflows(c *gin.Context) {
	db := database.GetDB()
	var workflows []models.WorkflowGroup
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	query := db.Model(&models.WorkflowGroup{})
	if keyword != "" {
		query = query.Where("name ILIKE ?", "%"+keyword+"%")
	}
	var total int64
	query.Model(&models.WorkflowGroup{}).Count(&total)
	err := query.Preload("Creator").Preload("Members.User").Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&workflows).Error
	if err != nil {
		c.JSON(500, gin.H{"error": "获取工作流列表失败"})
		return
	}
	c.JSON(200, gin.H{
		"data": workflows,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

// 创建工作流
func CreateWorkflow(c *gin.Context) {
	db := database.GetDB()
	userID, _ := c.Get("user_id")

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Members     []uint `json:"members"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	workflow := models.WorkflowGroup{
		Name:        req.Name,
		Description: req.Description,
		Status:      "active",
		CreatedBy:   userID.(uint),
	}
	if err := db.Create(&workflow).Error; err != nil {
		c.JSON(500, gin.H{"error": "创建工作流失败"})
		return
	}

	// 添加成员
	for _, uid := range req.Members {
		member := models.WorkflowMember{
			WorkflowID: workflow.ID,
			UserID:     uid,
			Role:       "member",
		}
		db.Create(&member)
	}

	db.Preload("Creator").Preload("Members.User").First(&workflow, workflow.ID)
	c.JSON(201, workflow)
}

// 获取单个工作流详情
func GetWorkflow(c *gin.Context) {
	db := database.GetDB()
	id := c.Param("id")
	var workflow models.WorkflowGroup
	err := db.Preload("Creator").Preload("Members.User").First(&workflow, id).Error
	if err != nil {
		c.JSON(404, gin.H{"error": "工作流不存在"})
		return
	}
	c.JSON(200, workflow)
}

// 更新工作流
func UpdateWorkflow(c *gin.Context) {
	db := database.GetDB()
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("role")

	var workflow models.WorkflowGroup
	if err := db.First(&workflow, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "工作流不存在"})
		return
	}
	if workflow.CreatedBy != userID.(uint) && userRole.(string) != "admin" {
		c.JSON(403, gin.H{"error": "没有权限"})
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Members     []uint `json:"members"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if err := db.Model(&workflow).Updates(updates).Error; err != nil {
		c.JSON(500, gin.H{"error": "更新失败"})
		return
	}

	// 更新成员（先删后加）
	if req.Members != nil {
		db.Where("workflow_id = ?", workflow.ID).Delete(&models.WorkflowMember{})
		for _, uid := range req.Members {
			member := models.WorkflowMember{
				WorkflowID: workflow.ID,
				UserID:     uid,
				Role:       "member",
			}
			db.Create(&member)
		}
	}
	db.Preload("Creator").Preload("Members.User").First(&workflow, workflow.ID)
	c.JSON(200, workflow)
}

// 删除工作流
func DeleteWorkflow(c *gin.Context) {
	db := database.GetDB()
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("role")

	var workflow models.WorkflowGroup
	if err := db.First(&workflow, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "工作流不存在"})
		return
	}
	if workflow.CreatedBy != userID.(uint) && userRole.(string) != "admin" {
		c.JSON(403, gin.H{"error": "没有权限"})
		return
	}

	// 先将该工作流下所有素材的 workflow_id 置为 NULL
	if err := db.Model(&models.Material{}).Where("workflow_id = ?", workflow.ID).Update("workflow_id", nil).Error; err != nil {
		c.JSON(500, gin.H{"error": "解除素材与工作流关系失败"})
		return
	}

	db.Where("workflow_id = ?", workflow.ID).Delete(&models.WorkflowMember{})
	db.Delete(&workflow)
	c.JSON(200, gin.H{"message": "删除成功，素材已解除关联"})
}

// 添加成员
func AddWorkflowMember(c *gin.Context) {
	db := database.GetDB()
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("role")

	var workflow models.WorkflowGroup
	if err := db.First(&workflow, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "工作流不存在"})
		return
	}
	if workflow.CreatedBy != userID.(uint) && userRole.(string) != "admin" {
		c.JSON(403, gin.H{"error": "没有权限"})
		return
	}

	var req struct {
		UserID uint   `json:"user_id" binding:"required"`
		Role   string `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if req.Role == "" {
		req.Role = "member"
	}
	member := models.WorkflowMember{
		WorkflowID: workflow.ID,
		UserID:     req.UserID,
		Role:       req.Role,
	}
	if err := db.Create(&member).Error; err != nil {
		c.JSON(500, gin.H{"error": "添加成员失败"})
		return
	}
	c.JSON(200, member)
}

// 移除成员
func RemoveWorkflowMember(c *gin.Context) {
	db := database.GetDB()
	id := c.Param("id")
	memberID := c.Param("userId")
	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("role")

	var workflow models.WorkflowGroup
	if err := db.First(&workflow, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "工作流不存在"})
		return
	}
	if workflow.CreatedBy != userID.(uint) && userRole.(string) != "admin" {
		c.JSON(403, gin.H{"error": "没有权限"})
		return
	}

	if err := db.Where("workflow_id = ? AND user_id = ?", id, memberID).Delete(&models.WorkflowMember{}).Error; err != nil {
		c.JSON(500, gin.H{"error": "移除成员失败"})
		return
	}
	c.JSON(200, gin.H{"message": "成员已移除"})
}
