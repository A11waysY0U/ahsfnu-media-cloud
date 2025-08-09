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
