package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	UID       uint           `gorm:"primaryKey;not null;autoIncrement" json:"uid"` //uid
	Email     string         `gorm:"size:128;uniqueIndex" json:"email"`            //邮箱
	Password  string         `gorm:"size:128" json:"password"`                     //密码
	Gender    string         `gorm:"size:4;default:'-'" json:"gender"`             //性别
	Nickname  string         `gorm:"size:512" json:"nickname"`                     //昵称
	CreatedAt time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoCreateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}
