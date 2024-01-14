package model

const TableNamePost = "posts"

type Post struct {
	BaseContent
	SiteID       string `json:"siteId" gorm:"primaryKey;uniqueIndex:,composite:_site_id"`
	Site         Site   `json:"-"`
	ID           string `json:"id" gorm:"primaryKey;size:100;uniqueIndex:,composite:_site_id"`
	IsDraft      bool   `json:"isDraft"`
	Draft        string `json:"-"`
	Body         string `json:"body"`
	PreviewURL   string `json:"previewUrl,omitempty" gorm:"size:200"`
	CategoryID   string `json:"categoryId,omitempty" gorm:"size:64;index:,composite:_category_id_path" label:"Category"`
	CategoryPath string `json:"categoryPath,omitempty" gorm:"size:64;index:,composite:_category_id_path"`
}

func (*Post) TableName() string {
	return TableNamePost
}
