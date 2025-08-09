package services

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/mojocn/base64Captcha"
)

// MemoryCaptchaStore 内存验证码存储
type MemoryCaptchaStore struct {
	store map[string]captchaItem
}

type captchaItem struct {
	Value     string
	ExpiresAt time.Time
}

// NewMemoryCaptchaStore 创建内存验证码存储
func NewMemoryCaptchaStore() *MemoryCaptchaStore {
	return &MemoryCaptchaStore{
		store: make(map[string]captchaItem),
	}
}

// Set 设置验证码
func (m *MemoryCaptchaStore) Set(id string, value string) error {
	m.store[id] = captchaItem{
		Value:     value,
		ExpiresAt: time.Now().Add(5 * time.Minute), // 5分钟过期
	}
	return nil
}

// Get 获取验证码 (实现 base64Captcha.Store 接口)
func (m *MemoryCaptchaStore) Get(id string, clear bool) string {
	item, exists := m.store[id]
	if !exists {
		return ""
	}
	if time.Now().After(item.ExpiresAt) {
		delete(m.store, id)
		return ""
	}

	if clear {
		delete(m.store, id)
	}
	return item.Value
}

// Verify 验证验证码 (实现 base64Captcha.Store 接口)
func (m *MemoryCaptchaStore) Verify(id, answer string, clear bool) bool {
	item, exists := m.store[id]
	if !exists {
		return false
	}
	if time.Now().After(item.ExpiresAt) {
		delete(m.store, id)
		return false
	}

	// 验证成功后删除验证码
	if item.Value == answer {
		if clear {
			delete(m.store, id)
		}
		return true
	}
	return false
}

// 全局验证码存储实例
var captchaStore = NewMemoryCaptchaStore()

// GenerateCaptcha 生成验证码
func GenerateCaptcha() (string, string, error) {
	// 生成随机ID
	id := generateRandomID()

	// 使用默认的数字验证码驱动
	driver := base64Captcha.NewDriverDigit(80, 240, 4, 0.7, 80)
	captcha := base64Captcha.NewCaptcha(driver, captchaStore)

	// 直接使用库的方法生成验证码和图片
	captchaID, captchaBase64, _, err := captcha.Generate()
	if err != nil {
		return "", "", err
	}

	// 从captchaID中获取答案
	answer := captchaStore.Get(captchaID, false)

	// 存储验证码答案到我们自己的ID
	captchaStore.Set(id, answer)

	return id, captchaBase64, nil
}

// VerifyCaptcha 验证验证码
func VerifyCaptcha(id, answer string) bool {
	return captchaStore.Verify(id, answer, true)
}

// GenerateAuthToken 生成认证token
func GenerateAuthToken() (string, error) {
	// 生成32字节随机数据
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// 创建token数据
	tokenData := map[string]interface{}{
		"token": base64.URLEncoding.EncodeToString(randomBytes),
		"exp":   time.Now().Add(10 * time.Minute).Unix(), // 10分钟过期
	}

	// 序列化为JSON
	tokenJSON, err := json.Marshal(tokenData)
	if err != nil {
		return "", err
	}

	// Base64编码
	return base64.URLEncoding.EncodeToString(tokenJSON), nil
}

// VerifyAuthToken 验证认证token
func VerifyAuthToken(token string) bool {
	// Base64解码
	tokenBytes, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return false
	}

	// 解析JSON
	var tokenData map[string]interface{}
	if err := json.Unmarshal(tokenBytes, &tokenData); err != nil {
		return false
	}

	// 检查过期时间
	exp, ok := tokenData["exp"].(float64)
	if !ok {
		return false
	}

	if time.Now().Unix() > int64(exp) {
		return false
	}

	return true
}

// generateRandomID 生成随机ID
func generateRandomID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)
}
