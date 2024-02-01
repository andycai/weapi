package page

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/andycai/weapi/enum"
	"github.com/andycai/weapi/model"
	"github.com/andycai/weapi/utils/random"
	"gorm.io/gorm"
)

func NewRenderContentFromPage(page *model.Page) *model.RenderContent {
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

	return &model.RenderContent{
		BaseContent: page.BaseContent,
		ID:          page.ID,
		SiteID:      page.SiteID,
		PageData:    data,
		IsDraft:     page.IsDraft,
	}
}

func MakeDuplicate(obj any) error {
	if page, ok := obj.(*model.Page); ok {
		page.ID = page.ID + "-copy-" + random.RandText(3)
		page.Title = page.Title + "-copy"
		page.IsDraft = true
		page.PreviewURL = ""
		page.Published = false
		page.CreatedAt = time.Now()
		page.UpdatedAt = time.Now()
		return db.Create(page).Error
	} else if post, ok := obj.(*model.Post); ok {
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

func MakePublish(siteID, ID string, obj any, publish bool) error {
	tx := db.Model(obj).Where("site_id", siteID).Where("id", ID)
	vals := map[string]any{"published": publish}

	vals["published"] = publish
	if publish {
		vals["body"] = gorm.Expr("draft")
		vals["is_draft"] = false
	}
	return tx.Updates(vals).Error
}

func SafeDraft(siteID, ID string, obj any, draft string) error {
	tx := db.Model(obj).Where("site_id", siteID).Where("id", ID)
	vals := map[string]any{
		"is_draft": true,
		"draft":    draft,
	}
	return tx.Updates(vals).Error
}

func QueryTags() ([]string, error) {
	var vals []string
	r := db.Model(&model.Page{}).Select("DISTINCT(tags)").Find(&vals)
	return vals, r.Error
}
