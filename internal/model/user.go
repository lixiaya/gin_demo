package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	UID       uint           `gorm:"primaryKey;not null;autoIncrement"` //uid
	Email     string         `gorm:"size:128;uniqueIndex"`              //邮箱
	Password  string         `gorm:"size:128"`                          //密码
	Gender    string         `gorm:"size:4;default:'-'"`                //性别
	Nickname  string         `gorm:"size:512"`                          //昵称
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"autoCreateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
