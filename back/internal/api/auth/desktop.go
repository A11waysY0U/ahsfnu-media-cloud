package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"ahsfnu-media-cloud/internal/config"
	"ahsfnu-media-cloud/internal/services"

	"github.com/gin-gonic/gin"
)

type DesktopAuthRequest struct {
	Nonce string `json:"nonce" binding:"required"` // 随机字符串
	HMAC  string `json:"hmac" binding:"required"`  // HMAC摘要
}

type DesktopAuthResponse struct {
	AuthToken string `json:"auth_token"` // 认证token
	Message   string `json:"message"`    // 响应消息
}

// DesktopAuth 桌面端HMAC认证接口
func DesktopAuth(c *gin.Context) {
	var req DesktopAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 1) 验证HMAC（避免攻击者用错误HMAC消耗nonce）
	if !verifyHMAC(req.Nonce, req.HMAC) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "HMAC验证失败，请求可能被伪造"})
		return
	}

	// 2) 一次性 nonce：首次使用通过，再次使用拒绝
	if ok := services.UseNonceOnce(req.Nonce); !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "nonce已使用或已失效"})
		return
	}

	// 3) 生成认证token
	authToken, err := services.GenerateAuthToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成认证token失败"})
		return
	}

	response := DesktopAuthResponse{
		AuthToken: authToken,
		Message:   "桌面端认证成功",
	}

	c.JSON(http.StatusOK, response)
}

// verifyHMAC 验证HMAC
func verifyHMAC(nonce, clientHMAC string) bool {
	// 获取HMAC密钥
	secretKey := config.AppConfig.HMAC.SecretKey

	// 使用HMAC-SHA256计算摘要
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(nonce))
	expectedHMAC := hex.EncodeToString(h.Sum(nil))

	// 比较计算得到的HMAC和客户端发送的HMAC
	return hmac.Equal([]byte(expectedHMAC), []byte(clientHMAC))
}
