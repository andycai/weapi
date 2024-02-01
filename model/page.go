package model

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/andycai/weapi/enum"
	"github.com/andycai/weapi/utils/random"
	"gorm.io/gorm"
)

const TableNamePage = "pages"

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

func (*Page) TableName() string {
	return TableNamePage
}

func MakeDuplicate(db *gorm.DB, obj any) error {
	if page, ok := obj.(*Page); ok {
		page.ID = page.ID + "-copy-" + random.RandText(3)
		page.Title = page.Title + "-copy"
		page.IsDraft = true
		page.PreviewURL = ""
		page.Published = false
		page.CreatedAt = time.Now()
		page.UpdatedAt = time.Now()
		return db.Create(page).Error
	} else if post, ok := obj.(*Post); ok {
		post.ID = post.ID + "-copy-" + random.RandText(3)
		post.Title = post.Title + "-copy"
		post.IsDraft = true
		post.PreviewURL = ""
		post.CreatedAt = time.Now()
		post.UpdatedAt = time.Now()
		post.Published = false
		return db.Create(post).Error
	}
	return errors.New("invalid object, must be page or post")
}

func MakePublish(db *gorm.DB, siteID, ID string, obj any, publish bool) error {
	tx := db.Model(obj).Where("site_id", siteID).Where("id", ID)
	vals := map[string]any{"published": publish}

	vals["published"] = publish
	if publish {
		vals["body"] = gorm.Expr("draft")
		vals["is_draft"] = false
	}
	return tx.Updates(vals).Error
}

func SafeDraft(db *gorm.DB, siteID, ID string, obj any, draft string) error {
	tx := db.Model(obj).Where("site_id", siteID).Where("id", ID)
	vals := map[string]any{
		"is_draft": true,
		"draft":    draft,
	}
	return tx.Updates(vals).Error
}

func QueryTags(db *gorm.DB) ([]string, error) {
	var vals []string
	r := db.Model(&Post{}).Select("DISTINCT(tags)").Find(&vals)
	return vals, r.Error
}

func NewRenderContentFromPage(db *gorm.DB, page *Page) *RenderContent {
	var data any
	if page.ContentType == enum.ContentTypeJson {
		data = make(map[string]any)
		err := json.Unmarshal([]byte(page.Body), &data)
		if err != nil {
			// carrot.Warning("unmarshal json error: ", page.SiteID, page.ID, page.Title, err)
		}
	} else {
		data = page.Body
	}

	return &RenderContent{
		BaseContent: page.BaseContent,
		ID:          page.ID,
		SiteID:      page.SiteID,
		PageData:    data,
		IsDraft:     page.IsDraft,
	}
}

func NewRenderContentFromPost(db *gorm.DB, post *Post, relations bool) *RenderContent {
	r := &RenderContent{
		BaseContent: post.BaseContent,
		ID:          post.ID,
		SiteID:      post.SiteID,
		PostBody:    post.Body,
		IsDraft:     post.IsDraft,
		Category:    NewRenderCategory(db, post.CategoryID, post.CategoryPath),
	}

	if relations {
		// relationCount := GetIntValue(db, enum.KEY_CMS_RELATION_COUNT, 3)
		// suggestionCount := GetIntValue(db, enum.KEY_CMS_SUGGESTION_COUNT, 3)

		// r.Relations, _ = GetRelations(db, post.SiteID, post.CategoryID, post.CategoryPath, post.ID, relationCount)
		// r.Suggestions, _ = GetSuggestions(db, post.SiteID, post.CategoryID, post.CategoryPath, post.ID, suggestionCount)
	}
	return r
}
