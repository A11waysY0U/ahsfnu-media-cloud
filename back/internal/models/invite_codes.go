package models

import (
	"time"
)

type InviteCode struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Code      string     `gorm:"unique;not null" json:"code"`
	CreatedBy uint       `json:"created_by"`
	UsedBy    *uint      `json:"used_by"`
	CreatedAt time.Time  `json:"created_at"`
	UsedAt    *time.Time `json:"used_at"`
	Status    int        `json:"status"`

	// 关联关系
	Creator *User `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
	User    *User `json:"user,omitempty" gorm:"foreignKey:UsedBy"`
}

// InviteCodeResponse 用于返回给前端的邀请码信息，包含安全的用户信息
type InviteCodeResponse struct {
	ID        uint       `json:"id"`
	Code      string     `json:"code"`
	CreatedBy uint       `json:"created_by"`
	UsedBy    *uint      `json:"used_by"`
	CreatedAt time.Time  `json:"created_at"`
	UsedAt    *time.Time `json:"used_at"`
	Status    int        `json:"status"`

	// 安全的关联关系
	Creator *SafeUser `json:"creator,omitempty"`
	User    *SafeUser `json:"user,omitempty"`
}

// ToInviteCodeResponse 将 InviteCode 转换为 InviteCodeResponse
func (ic *InviteCode) ToInviteCodeResponse() *InviteCodeResponse {
	response := &InviteCodeResponse{
		ID:        ic.ID,
		Code:      ic.Code,
		CreatedBy: ic.CreatedBy,
		UsedBy:    ic.UsedBy,
		CreatedAt: ic.CreatedAt,
		UsedAt:    ic.UsedAt,
		Status:    ic.Status,
	}

	// 安全地转换创建者信息
	if ic.Creator != nil {
		response.Creator = ic.Creator.ToSafeUser()
	}

	// 安全地转换使用者信息
	if ic.User != nil {
		response.User = ic.User.ToSafeUser()
	}

	return response
}
