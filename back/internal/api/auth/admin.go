package auth

import (
	"ahsfnu-media-cloud/internal/database"
	"ahsfnu-media-cloud/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateUserRole(c *gin.Context) {
	exists := false
	if _, ok := c.Get("user_id"); ok {
		exists = true
	}
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限"})
		return
	}
	db := database.GetDB()
	targetID := c.Param("id")
	var req struct {
		Role string `json:"role" binding:"required,oneof=admin user"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	if err := db.First(&user, targetID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	user.Role = req.Role
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "权限修改失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "权限修改成功", "user": user})
}

func DeleteUser(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(403, gin.H{"error": "无权限"})
		return
	}
	db := database.GetDB()
	id := c.Param("id")
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "用户不存在"})
		return
	}
	if err := db.Delete(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(200, gin.H{"message": "删除成功"})
}

func GetUsers(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(403, gin.H{"error": "无权限"})
		return
	}
	db := database.GetDB()
	var users []models.User
	if err := db.Preload("Inviter").Find(&users).Error; err != nil {
		c.JSON(500, gin.H{"error": "获取用户列表失败"})
		return
	}
	c.JSON(200, users)
}
