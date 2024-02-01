package model

import (
	"time"
)

type Page struct {
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

type PublishLog struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt"`
	AuthorID  uint      `json:"-"`
	Author    User      `json:"author"`
	Content   string    `json:"content" gorm:"size:12;index:idx_content_with_id"`    // post_id or page_id
	ContentID string    `json:"contentId" gorm:"size:100;index:idx_content_with_id"` // post_id or page_id
	Body      string    `json:"body"`
}

type RelationContent struct {
	BaseContent
	SiteID string `json:"siteId"`
	ID     string `json:"id"`
}

type RenderContent struct {
	BaseContent
	ID          string            `json:"id"`
	SiteID      string            `json:"siteId"`
	Category    *RenderCategory   `json:"category,omitempty"`
	PageData    any               `json:"data,omitempty"`
	PostBody    string            `json:"body,omitempty"`
	IsDraft     bool              `json:"isDraft"`
	Relations   []RelationContent `json:"relations,omitempty"`
	Suggestions []RelationContent `json:"suggestions,omitempty"`
}
type ContentQueryResult struct {
	*QueryResult
	Relations   []RelationContent `json:"relations,omitempty"`
	Suggestions []RelationContent `json:"suggestions,omitempty"`
}
