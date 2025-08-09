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
