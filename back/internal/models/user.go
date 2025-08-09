package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Password  string    `gorm:"not null" json:"password"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Role      string    `gorm:"type:varchar(20);default:'user'" json:"role"`
	InviterID *uint     `json:"inviter_id"`
	CreatedAt time.Time `json:"created_at"`

	// 关联关系
	Materials       []Material       `json:"materials,omitempty" gorm:"foreignKey:UploadedBy"`
	Tags            []Tag            `json:"tags,omitempty" gorm:"foreignKey:CreatedBy"`
	Workflows       []WorkflowGroup  `json:"workflows,omitempty" gorm:"foreignKey:CreatedBy"`
	WorkflowMembers []WorkflowMember `json:"workflow_members,omitempty" gorm:"foreignKey:UserID"`
	MaterialTags    []MaterialTag    `json:"material_tags,omitempty" gorm:"foreignKey:CreatedBy"`
	Inviter         *User            `json:"inviter,omitempty" gorm:"foreignKey:InviterID"`
	Invitees        []User           `json:"invitees,omitempty" gorm:"foreignKey:InviterID"`
}

// SafeUser 安全的用户信息结构体，用于返回给前端，不包含敏感信息
type SafeUser struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	InviterID *uint     `json:"inviter_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// ToSafeUser 将 User 转换为 SafeUser
func (u *User) ToSafeUser() *SafeUser {
	return &SafeUser{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role,
		InviterID: u.InviterID,
		CreatedAt: u.CreatedAt,
	}
}
