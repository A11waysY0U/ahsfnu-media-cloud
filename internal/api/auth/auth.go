package auth

import (
	"net/http"

	"ahsfnu-media-cloud/internal/database"
	"ahsfnu-media-cloud/internal/models"
	"ahsfnu-media-cloud/internal/services"
	"ahsfnu-media-cloud/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	AuthToken string `json:"auth_token" binding:"required"`
}

type RegisterRequest struct {
	Username   string `json:"username" binding:"required,min=3,max=50"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	InviteCode string `json:"invite_code" binding:"required"`
	AuthToken  string `json:"auth_token" binding:"required"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// 登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证认证token
	if !services.VerifyAuthToken(req.AuthToken) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "认证token无效或已过期"})
		return
	}

	db := database.GetDB()
	var user models.User

	// 通过用户名或邮箱查找用户
	if err := db.Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		return
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	response := AuthResponse{
		Token: token,
		User:  user,
	}

	c.JSON(http.StatusOK, response)
}

// Register 用户注册
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证认证token
	if !services.VerifyAuthToken(req.AuthToken) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "认证token无效或已过期"})
		return
	}

	db := database.GetDB()

	// 检查用户名是否已存在
	var existingUser models.User
	if err := db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	// 检查邮箱是否已存在
	if err := db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱已存在"})
		return
	}

	// 校验邀请码
	invite, err := services.GetInviteCode(db, req.InviteCode)
	if err != nil || invite.Status != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邀请码无效或已被使用"})
		return
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 创建新用户
	user := models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		InviterID: &invite.CreatedBy,
		Role:      "user", // 注册用户强制为普通用户
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	// 标记邀请码已用
	if err := services.UseInviteCode(db, req.InviteCode, user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "邀请码状态更新失败"})
		return
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	response := AuthResponse{
		Token: token,
		User:  user,
	}

	c.JSON(http.StatusCreated, response)
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
