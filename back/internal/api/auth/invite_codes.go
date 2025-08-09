package auth

import (
	"ahsfnu-media-cloud/internal/database"
	"ahsfnu-media-cloud/internal/models"
	"ahsfnu-media-cloud/internal/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GenerateInviteCodes(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限"})
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

	// 转换为安全的响应格式
	var codeResponses []models.InviteCodeResponse
	for i := range codes {
		codeResponse := codes[i].ToInviteCodeResponse()
		codeResponses = append(codeResponses, *codeResponse)
	}

	c.JSON(http.StatusOK, gin.H{"data": codeResponses})
}

func ListInviteCodes(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限"})
		return
	}
	db := database.GetDB()
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	offset := (page - 1) * pageSize
	var codes []models.InviteCode
	var total int64
	db.Model(&models.InviteCode{}).Count(&total)
	err := db.Preload("Creator").Preload("User").Order("created_at desc").Offset(offset).Limit(pageSize).Find(&codes).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取邀请码列表失败"})
		return
	}

	// 转换为安全的响应格式
	var codeResponses []models.InviteCodeResponse
	for i := range codes {
		codeResponse := codes[i].ToInviteCodeResponse()
		codeResponses = append(codeResponses, *codeResponse)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": codeResponses,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

// GetInviteCodeStats 获取邀请码统计信息
func GetInviteCodeStats(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限"})
		return
	}

	db := database.GetDB()

	// 统计总数
	var total int64
	db.Model(&models.InviteCode{}).Count(&total)

	// 统计未使用的邀请码（status = 0）
	var unused int64
	db.Model(&models.InviteCode{}).Where("status = ?", 0).Count(&unused)

	// 统计已使用的邀请码（status = 1）
	var used int64
	db.Model(&models.InviteCode{}).Where("status = ?", 1).Count(&used)

	// 统计已过期的邀请码（创建时间超过30天的未使用邀请码）
	var expired int64
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	db.Model(&models.InviteCode{}).Where("status = ? AND created_at < ?", 0, thirtyDaysAgo).Count(&expired)

	stats := gin.H{
		"total":   total,
		"unused":  unused,
		"used":    used,
		"expired": expired,
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}
