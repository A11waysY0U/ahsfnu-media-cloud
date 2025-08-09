package models

import (
	"time"
)

type WorkflowGroup struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name" gorm:"not null;size:100"`
	Description string     `json:"description,omitempty"`
	Type        string     `json:"type" gorm:"default:'custom';size:50"`  // image_processing, video_processing, file_conversion, batch_operation, custom
	Color       string     `json:"color" gorm:"default:'#409EFF';size:7"` // 十六进制颜色
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	Config      string     `json:"config,omitempty" gorm:"type:text"`      // JSON配置
	Status      string     `json:"status" gorm:"default:'active';size:20"` // active, archived
	CreatedBy   uint       `json:"created_by" gorm:"not null"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	EndedAt     *time.Time `json:"ended_at,omitempty"`

	// 关联关系
	Creator   *User            `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
	Materials []Material       `json:"materials,omitempty" gorm:"foreignKey:WorkflowID"`
	Members   []WorkflowMember `json:"members,omitempty" gorm:"foreignKey:WorkflowID"`
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

// WorkflowGroupResponse 用于返回给前端的工作流信息，包含安全的用户信息
type WorkflowGroupResponse struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Type        string     `json:"type"`
	Color       string     `json:"color"`
	IsActive    bool       `json:"is_active"`
	Config      string     `json:"config,omitempty"`
	Status      string     `json:"status"`
	CreatedBy   uint       `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	EndedAt     *time.Time `json:"ended_at,omitempty"`

	// 安全的关联关系
	Creator   *SafeUser                `json:"creator,omitempty"`
	Materials []Material               `json:"materials,omitempty"`
	Members   []WorkflowMemberResponse `json:"members,omitempty"`
}

// WorkflowMemberResponse 用于返回给前端的工作流成员信息，包含安全的用户信息
type WorkflowMemberResponse struct {
	ID         uint      `json:"id"`
	WorkflowID uint      `json:"workflow_id"`
	UserID     uint      `json:"user_id"`
	Role       string    `json:"role"`
	JoinedAt   time.Time `json:"joined_at"`

	// 安全的关联关系
	Workflow *WorkflowGroupResponse `json:"workflow,omitempty"`
	User     *SafeUser              `json:"user,omitempty"`
}

// ToWorkflowGroupResponse 将 WorkflowGroup 转换为 WorkflowGroupResponse
func (w *WorkflowGroup) ToWorkflowGroupResponse() *WorkflowGroupResponse {
	response := &WorkflowGroupResponse{
		ID:          w.ID,
		Name:        w.Name,
		Description: w.Description,
		Type:        w.Type,
		Color:       w.Color,
		IsActive:    w.IsActive,
		Config:      w.Config,
		Status:      w.Status,
		CreatedBy:   w.CreatedBy,
		CreatedAt:   w.CreatedAt,
		EndedAt:     w.EndedAt,
		Materials:   w.Materials,
	}

	// 安全地转换创建者信息
	if w.Creator != nil {
		response.Creator = w.Creator.ToSafeUser()
	}

	// 安全地转换成员信息
	if w.Members != nil {
		response.Members = make([]WorkflowMemberResponse, len(w.Members))
		for i, member := range w.Members {
			response.Members[i] = *member.ToWorkflowMemberResponse()
		}
	}

	return response
}

// ToWorkflowMemberResponse 将 WorkflowMember 转换为 WorkflowMemberResponse
func (w *WorkflowMember) ToWorkflowMemberResponse() *WorkflowMemberResponse {
	response := &WorkflowMemberResponse{
		ID:         w.ID,
		WorkflowID: w.WorkflowID,
		UserID:     w.UserID,
		Role:       w.Role,
		JoinedAt:   w.JoinedAt,
	}

	// 安全地转换用户信息
	if w.User != nil {
		response.User = w.User.ToSafeUser()
	}

	// 安全地转换工作流信息（避免循环引用，只返回基本信息）
	if w.Workflow != nil {
		response.Workflow = &WorkflowGroupResponse{
			ID:          w.Workflow.ID,
			Name:        w.Workflow.Name,
			Description: w.Workflow.Description,
			Status:      w.Workflow.Status,
			CreatedBy:   w.Workflow.CreatedBy,
			CreatedAt:   w.Workflow.CreatedAt,
			EndedAt:     w.Workflow.EndedAt,
		}
	}

	return response
}
