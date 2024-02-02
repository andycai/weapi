package model

import (
	"fmt"
	"time"
)

type Site struct {
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	Domain    string    `json:"domain" gorm:"primarykey;size:200"`
	Name      string    `json:"name" gorm:"size:200"`
	Preview   string    `json:"preview" gorm:"size:200"`
	Disallow  bool      `json:"disallow"`
}

func (s Site) String() string {
	return fmt.Sprintf("%s(%s)", s.Name, s.Domain)
}
