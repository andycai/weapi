package model

import (
	"time"
)

type Page struct {
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

type PublishLog struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	AuthorID  uint      `json:"-"`
	Author    User      `json:"author"`
	Content   string    `json:"content" gorm:"size:12;index:idx_content_with_id"`     // post_id or page_id
	ContentID string    `json:"content_id" gorm:"size:100;index:idx_content_with_id"` // post_id or page_id
	Body      string    `json:"body"`
}

type RelationContent struct {
	BaseContent
	SiteID string `json:"site_id"`
	ID     string `json:"id"`
}

type RenderContent struct {
	BaseContent
	ID          string            `json:"id"`
	SiteID      string            `json:"site_id"`
	Category    *RenderCategory   `json:"category,omitempty"`
	PageData    any               `json:"data,omitempty"`
	PostBody    string            `json:"body,omitempty"`
	IsDraft     bool              `json:"is_draft"`
	Relations   []RelationContent `json:"relations,omitempty"`
	Suggestions []RelationContent `json:"suggestions,omitempty"`
}
type ContentQueryResult struct {
	*QueryResult
	Relations   []RelationContent `json:"relations,omitempty"`
	Suggestions []RelationContent `json:"suggestions,omitempty"`
}
