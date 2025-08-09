package models

import (
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null;size:100"`
	Color     string    `json:"color" gorm:"default:'#409EFF';size:7"` // 十六进制颜色
	CreatedBy uint      `json:"created_by" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	// 关联关系
	Creator      User          `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
	MaterialTags []MaterialTag `json:"material_tags,omitempty" gorm:"foreignKey:TagID"`
}

type MaterialTag struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	MaterialID uint      `json:"material_id" gorm:"not null"`
	TagID      uint      `json:"tag_id" gorm:"not null"`
	CreatedBy  uint      `json:"created_by" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`

	// 关联关系
	Material Material `json:"material,omitempty" gorm:"foreignKey:MaterialID"`
	Tag      Tag      `json:"tag,omitempty" gorm:"foreignKey:TagID"`
	Creator  User     `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
}

func (mt *MaterialTag) BeforeCreate(tx *gorm.DB) error {
	// 确保一个素材不会重复添加同一个标签
	var count int64
	tx.Model(&MaterialTag{}).Where("material_id = ? AND tag_id = ?", mt.MaterialID, mt.TagID).Count(&count)
	if count > 0 {
		return gorm.ErrDuplicatedKey
	}
	return nil
}
