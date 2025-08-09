package services

import (
	"ahsfnu-media-cloud/internal/models"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

// 生成指定数量的邀请码
func GenerateInviteCodes(db *gorm.DB, count int, adminID uint) ([]models.InviteCode, error) {
	codes := make([]models.InviteCode, 0, count)
	for i := 0; i < count; i++ {
		code := randomCode(10)
		invite := models.InviteCode{
			Code:      code,
			CreatedBy: adminID,
			Status:    0,
			CreatedAt: time.Now(),
		}
		if err := db.Create(&invite).Error; err != nil {
			return nil, err
		}
		codes = append(codes, invite)
	}
	return codes, nil
}

// 查询邀请码
func GetInviteCode(db *gorm.DB, code string) (*models.InviteCode, error) {
	var invite models.InviteCode
	if err := db.Where("code = ?", code).First(&invite).Error; err != nil {
		return nil, err
	}
	return &invite, nil
}

// 标记邀请码已用
func UseInviteCode(db *gorm.DB, code string, userID uint) error {
	var invite models.InviteCode
	if err := db.Where("code = ? AND status = 0", code).First(&invite).Error; err != nil {
		return err
	}
	now := time.Now()
	invite.UsedBy = &userID
	invite.UsedAt = &now
	invite.Status = 1
	return db.Save(&invite).Error
}

// 随机生成邀请码字符串
func randomCode(length int) string {
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
