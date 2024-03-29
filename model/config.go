package model

type Config struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Key   string `json:"key" gorm:"size:128;uniqueIndex"`
	Desc  string `json:"desc" gorm:"size:200"`
	Value string
}
