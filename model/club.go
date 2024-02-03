package model

import (
	"fmt"
)

type Club struct {
	BaseField
	Scores  uint   `json:"scores" gorm:"not null;comment:俱乐部积分"` // 俱乐部积分
	Level   uint   `json:"level" gorm:"default:1;comment:俱乐部等级"` // 俱乐部等级
	Logo    string `json:"logo" gorm:"comment:俱乐部LOGO"`          // 俱乐部LOGO
	Notice  string `json:"notice" gorm:"comment:俱乐部公告"`          // 俱乐部公告
	Address string `json:"address" gorm:"comment:俱乐部总部地址"`       // 俱乐部总部地址
}

func (s Club) String() string {
	return fmt.Sprintf("%s(%d)", s.Name, s.ID)
}
