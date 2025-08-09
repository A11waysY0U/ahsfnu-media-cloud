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

// TagResponse 用于返回给前端的标签信息，包含安全的用户信息
type TagResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedBy uint      `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`

	// 安全的关联关系
	Creator      *SafeUser     `json:"creator,omitempty"`
	MaterialTags []MaterialTag `json:"material_tags,omitempty"`
}

// MaterialTagResponse 用于返回给前端的素材标签信息，包含安全的用户信息
type MaterialTagResponse struct {
	ID         uint      `json:"id"`
	MaterialID uint      `json:"material_id"`
	TagID      uint      `json:"tag_id"`
	CreatedBy  uint      `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`

	// 安全的关联关系
	Material *MaterialResponse `json:"material,omitempty"`
	Tag      *TagResponse      `json:"tag,omitempty"`
	Creator  *SafeUser         `json:"creator,omitempty"`
}

// ToTagResponse 将 Tag 转换为 TagResponse
func (t *Tag) ToTagResponse() *TagResponse {
	response := &TagResponse{
		ID:           t.ID,
		Name:         t.Name,
		Color:        t.Color,
		CreatedBy:    t.CreatedBy,
		CreatedAt:    t.CreatedAt,
		MaterialTags: t.MaterialTags,
	}

	// 安全地转换创建者信息
	if t.Creator.ID != 0 {
		response.Creator = t.Creator.ToSafeUser()
	}

	return response
}

// ToMaterialTagResponse 将 MaterialTag 转换为 MaterialTagResponse
func (mt *MaterialTag) ToMaterialTagResponse() *MaterialTagResponse {
	response := &MaterialTagResponse{
		ID:         mt.ID,
		MaterialID: mt.MaterialID,
		TagID:      mt.TagID,
		CreatedBy:  mt.CreatedBy,
		CreatedAt:  mt.CreatedAt,
	}

	// 安全地转换创建者信息
	if mt.Creator.ID != 0 {
		response.Creator = mt.Creator.ToSafeUser()
	}

	// 安全地转换素材信息
	if mt.Material.ID != 0 {
		response.Material = mt.Material.ToMaterialResponse()
	}

	// 安全地转换标签信息
	if mt.Tag.ID != 0 {
		response.Tag = mt.Tag.ToTagResponse()
	}

	return response
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
