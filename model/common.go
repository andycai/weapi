package model

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"
)

type StringArray []string

func (s StringArray) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *StringArray) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &s)
}

type BaseContent struct {
	UpdatedAt   time.Time    `json:"updated_at" gorm:"index"`
	CreatedAt   time.Time    `json:"created_at" gorm:"index"`
	Thumbnail   string       `json:"thumbnail,omitempty" gorm:"size:500"`
	Tags        string       `json:"tags,omitempty" gorm:"size:200;index"`
	Title       string       `json:"title,omitempty" gorm:"size:200"`
	Alt         string       `json:"alt,omitempty"`
	Description string       `json:"description,omitempty"`
	Keywords    string       `json:"keywords,omitempty"`
	CreatorID   uint         `json:"-"`
	Creator     User         `json:"-"`
	Author      string       `json:"author" gorm:"size:64"`
	Published   bool         `json:"published"`
	PublishedAt sql.NullTime `json:"published_at"`
	ContentType string       `json:"content_type" gorm:"size:32"`
	Remark      string       `json:"remark"`
}

type SummaryResult struct {
	SiteCount     int64            `json:"sites"`
	PageCount     int64            `json:"pages"`
	PostCount     int64            `json:"posts"`
	CategoryCount int64            `json:"categories"`
	MediaCount    int64            `json:"media"`
	LatestPosts   []*RenderContent `json:"latestPosts"`
	BuildTime     string           `json:"buildTime"`
	CanExport     bool             `json:"canExport"`
}

type TagsForm struct {
	SiteId       string `json:"site_id"`
	CategoryId   string `json:"category_id"`
	CategoryPath string `json:"category_path"`
}

type QueryByTagsForm struct {
	Tags  []string `json:"tags" binding:"required"`
	Limit int      `json:"limit"`
	Pos   int      `json:"pos"`
	TagsForm
}
type QueryByTagsResult struct {
	Items []any `json:"items"`
	Total int   `json:"total"`
	Limit int   `json:"limit"`
	Pos   int   `json:"pos"`
}

type BaseField struct {
	ID          uint      `json:"id" gorm:"primarykey;autoIncrement:true"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoCreateTime"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoUpdateTime"`
	SiteID      string    `json:"site_id" gorm:"uniqueIndex:,composite:_site_name"`
	Site        Site      `json:"-"`
	Name        string    `json:"title" gorm:"size:200;uniqueIndex:,composite:_site_name"`
	Description string    `json:"description,omitempty"`
	CreatorID   uint      `json:"-" gorm:"index"`
	Creator     User      `json:"-"`
}
