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
	Inviter   *User     `gorm:"foreignKey:InviterID" json:"inviter"`
	CreatedAt time.Time `json:"created_at"`
}
