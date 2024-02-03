package model

import (
	"fmt"
	"time"
)

const (
	ClubPositionOwner   = 1
	ClubPositionMember  = 2
	ClubPositionManager = 3
)

type Club struct {
	BaseField
	Scores  uint   `json:"scores" gorm:"default:0;comment:俱乐部积分"` // 俱乐部积分
	Level   uint   `json:"level" gorm:"default:1;comment:俱乐部等级"`  // 俱乐部等级
	Logo    string `json:"logo" gorm:"comment:俱乐部LOGO"`           // 俱乐部LOGO
	Notice  string `json:"notice" gorm:"comment:俱乐部公告"`           // 俱乐部公告
	Address string `json:"address" gorm:"comment:俱乐部总部地址"`        // 俱乐部总部地址
}

type ClubMember struct {
	ID          uint      `json:"-" gorm:"primaryKey;autoIncrement"`
	CreatedAt   time.Time `json:"-" gorm:"autoCreateTime"`
	Scores      uint      `json:"scores" gorm:"default:0"`
	UserID      uint      `json:"-"`
	User        User      `json:"user"`
	ClubID      uint      `json:"-"`
	Club        Club      `json:"club"`
	Position    uint      `json:"position" gorm:"default:1;comment:群组职位"`    // 群组职位
	DisplayName string    `json:"display_name" gorm:"size:200;comment:群组昵称"` // 群组昵称
	EnterAt     time.Time `json:"enter_at" gorm:"comment:进入群组时间"`            // 进入群组时间
}

func (s Club) String() string {
	return fmt.Sprintf("%s(%d)", s.Name, s.ID)
}
