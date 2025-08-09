package models

import (
	"time"
)

type Material struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	Filename         string    `json:"filename" gorm:"not null;size:255"`
	OriginalFilename string    `json:"original_filename" gorm:"not null;size:255"`
	FilePath         string    `json:"file_path" gorm:"not null;size:500"`
	FileSize         int64     `json:"file_size" gorm:"not null"`
	FileType         string    `json:"file_type" gorm:"not null;size:50"` // image, video
	MimeType         string    `json:"mime_type" gorm:"not null;size:100"`
	Width            *int      `json:"width,omitempty"`
	Height           *int      `json:"height,omitempty"`
	Duration         *int      `json:"duration,omitempty"` // 视频时长(秒)
	UploadedBy       uint      `json:"uploaded_by" gorm:"not null"`
	WorkflowID       *uint     `json:"workflow_id,omitempty"`
	UploadTime       time.Time `json:"upload_time" gorm:"autoCreateTime"`
	IsStarred        bool      `json:"is_starred" gorm:"default:false"`
	IsPublic         bool      `json:"is_public" gorm:"default:false"` // 是否公开
	ThumbnailPath    string    `json:"thumbnail_path,omitempty" gorm:"size:500"`

	// 关联关系
	Uploader     *User          `json:"uploader,omitempty" gorm:"foreignKey:UploadedBy"`
	Workflow     *WorkflowGroup `json:"workflow,omitempty" gorm:"foreignKey:WorkflowID"`
	MaterialTags []MaterialTag  `json:"material_tags,omitempty" gorm:"foreignKey:MaterialID"`
}

// MaterialResponse 用于返回给前端的素材信息，包含安全的用户信息
type MaterialResponse struct {
	ID               uint      `json:"id"`
	Filename         string    `json:"filename"`
	OriginalFilename string    `json:"original_filename"`
	FilePath         string    `json:"file_path"`
	FileSize         int64     `json:"file_size"`
	FileType         string    `json:"file_type"`
	MimeType         string    `json:"mime_type"`
	Width            *int      `json:"width,omitempty"`
	Height           *int      `json:"height,omitempty"`
	Duration         *int      `json:"duration,omitempty"`
	UploadedBy       uint      `json:"uploaded_by"`
	WorkflowID       *uint     `json:"workflow_id,omitempty"`
	UploadTime       time.Time `json:"upload_time"`
	IsStarred        bool      `json:"is_starred"`
	IsPublic         bool      `json:"is_public"`
	ThumbnailPath    string    `json:"thumbnail_path,omitempty"`

	// 安全的关联关系
	Uploader     *SafeUser      `json:"uploader,omitempty"`
	Workflow     *WorkflowGroup `json:"workflow,omitempty"`
	MaterialTags []MaterialTag  `json:"material_tags,omitempty"`
}

// ToMaterialResponse 将 Material 转换为 MaterialResponse
func (m *Material) ToMaterialResponse() *MaterialResponse {
	response := &MaterialResponse{
		ID:               m.ID,
		Filename:         m.Filename,
		OriginalFilename: m.OriginalFilename,
		FilePath:         m.FilePath,
		FileSize:         m.FileSize,
		FileType:         m.FileType,
		MimeType:         m.MimeType,
		Width:            m.Width,
		Height:           m.Height,
		Duration:         m.Duration,
		UploadedBy:       m.UploadedBy,
		WorkflowID:       m.WorkflowID,
		UploadTime:       m.UploadTime,
		IsStarred:        m.IsStarred,
		IsPublic:         m.IsPublic,
		ThumbnailPath:    m.ThumbnailPath,
		Workflow:         m.Workflow,
		MaterialTags:     m.MaterialTags,
	}

	// 安全地转换用户信息
	if m.Uploader != nil {
		response.Uploader = m.Uploader.ToSafeUser()
	}

	return response
}
