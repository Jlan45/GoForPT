package model

import "gorm.io/gorm"

type UserGroup struct {
	gorm.Model
	Name  string
	Level int
}
type User struct {
	gorm.Model
	Username    string
	Email       string
	PassHash    string
	Token       string //作为种子鉴别身份的token
	Uploaded    uint64
	Downloaded  uint64
	MagicPower  float32
	UserGroupID uint
	UserGroup   UserGroup `gorm:"foreignKey:UserGroupID"`
}
