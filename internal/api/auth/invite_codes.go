package auth

import (
	"ahsfnu-media-cloud/internal/database"
	"ahsfnu-media-cloud/internal/models"
	"ahsfnu-media-cloud/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GenerateInviteCodes(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(403, gin.H{"error": "无权限"})
		return
	}
	db := database.GetDB()
	var req struct {
		Count int `json:"count" binding:"required,min=1,max=100"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}
	codes, err := services.GenerateInviteCodes(db, req.Count, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成邀请码失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": codes})
}

// 查询邀请码列表（仅管理员可用）
func ListInviteCodes(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(403, gin.H{"error": "无权限"})
		return
	}
	db := database.GetDB()
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	offset := (page - 1) * pageSize
	var codes []models.InviteCode
	var total int64
	db.Model(&models.InviteCode{}).Count(&total)
	db.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&codes)
	c.JSON(http.StatusOK, gin.H{
		"data": codes,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}
