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
