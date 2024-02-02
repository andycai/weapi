package model

import "time"

type Club struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;index"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;index"`
	Scores    uint      `json:"scores" gorm:"column:scores;not null;comment:积分"`           // 积分
	Level     uint      `json:"level" gorm:"column:level;not null;default:1;comment:群组等级"` // 群组等级
	Name      string    `json:"name" gorm:"column:name;not null;comment:群组名称"`             // 群组名称
	Logo      string    `json:"logo" gorm:"column:logo;comment:群组LOGO"`                    // 群组LOGO
	Notice    string    `json:"notice" gorm:"column:notice;comment:群组公告"`                  // 群组公告
	Addr      string    `json:"addr" gorm:"column:addr;comment:群组总部地址"`                    // 群组总部地址
}
