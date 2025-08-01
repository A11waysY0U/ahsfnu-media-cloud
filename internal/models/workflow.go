package models

import (
	"time"
)

type WorkflowGroup struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name" gorm:"not null;size:100"`
	Description string     `json:"description,omitempty"`
	Status      string     `json:"status" gorm:"default:'active';size:20"` // active, archived
	CreatedBy   uint       `json:"created_by" gorm:"not null"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	EndedAt     *time.Time `json:"ended_at,omitempty"`

	// 关联关系
	Creator   *User      `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
	Materials []Material `json:"materials,omitempty" gorm:"foreignKey:WorkflowID"`
}

type WorkflowMember struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	WorkflowID uint      `json:"workflow_id" gorm:"not null"`
	UserID     uint      `json:"user_id" gorm:"not null"`
	Role       string    `json:"role" gorm:"not null;size:20"` // photographer, reviewer, admin
	JoinedAt   time.Time `json:"joined_at" gorm:"autoCreateTime"`

	// 关联关系
	Workflow *WorkflowGroup `json:"workflow,omitempty" gorm:"foreignKey:WorkflowID"`
	User     *User          `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}
