package auth

import (
	"net/http"

	"ahsfnu-media-cloud/internal/services"

	"github.com/gin-gonic/gin"
)

type CaptchaRequest struct {
	CaptchaID   string `json:"captcha_id" binding:"required"`
	CaptchaCode string `json:"captcha_code" binding:"required"`
}

type CaptchaResponse struct {
	CaptchaID  string `json:"captcha_id"`
	CaptchaB64 string `json:"captcha_b64"`
	AuthToken  string `json:"auth_token"`
}

// GetCaptcha 获取验证码
func GetCaptcha(c *gin.Context) {
	// 生成验证码
	captchaID, captchaB64, err := services.GenerateCaptcha()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成验证码失败"})
		return
	}

	// 生成认证token
	authToken, err := services.GenerateAuthToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成认证token失败"})
		return
	}

	response := CaptchaResponse{
		CaptchaID:  captchaID,
		CaptchaB64: captchaB64,
		AuthToken:  authToken,
	}

	c.JSON(http.StatusOK, response)
}

// VerifyCaptcha 验证验证码并返回认证token
func VerifyCaptcha(c *gin.Context) {
	var req CaptchaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证验证码
	if !services.VerifyCaptcha(req.CaptchaID, req.CaptchaCode) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码错误"})
		return
	}

	// 生成新的认证token
	authToken, err := services.GenerateAuthToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成认证token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"auth_token": authToken,
		"message":    "验证码验证成功",
	})
}
