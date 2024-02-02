package model

type Post struct {
	BaseContent
	SiteID       string `json:"site_id" gorm:"primaryKey;uniqueIndex:,composite:_site_id"`
	Site         Site   `json:"-"`
	ID           string `json:"id" gorm:"primaryKey;size:100;uniqueIndex:,composite:_site_id"`
	IsDraft      bool   `json:"is_draft"`
	Draft        string `json:"-"`
	Body         string `json:"body"`
	PreviewURL   string `json:"preview_url,omitempty" gorm:"size:200"`
	CategoryID   string `json:"category_id,omitempty" gorm:"size:64;index:,composite:_category_id_path" label:"Category"`
	CategoryPath string `json:"category_path,omitempty" gorm:"size:64;index:,composite:_category_id_path"`
}
