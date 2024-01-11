package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameUser = "users"

type User struct {
	gorm.Model
	Username    string    `gorm:"column:username;not null" json:"username"`
	Password    string    `gorm:"column:password;not null;default:123456" json:"password"`
	Name        string    `gorm:"column:name;not null" json:"name"`
	Token       string    `gorm:"column:token;not null" json:"token"`
	Avatar      string    `gorm:"column:avatar;not null" json:"avatar"`
	Gender      uint      `gorm:"column:gender;not null;default:1" json:"gender"`
	Phone       string    `gorm:"column:phone;not null" json:"phone"`
	Email       string    `gorm:"column:email;not null" json:"email"`
	Addr        string    `gorm:"column:addr;not null" json:"addr"`
	IP          string    `gorm:"column:ip;not null" json:"ip"`
	IsSuperUser bool      `gorm:"column:is_super_user;not null" json:"-"`
	LoginAt     time.Time `gorm:"column:login_at" json:"login_at"`
	LogoutAt    time.Time `gorm:"column:logout_at" json:"logout_at"`
}

func (*User) TableName() string {
	return TableNameUser
}
