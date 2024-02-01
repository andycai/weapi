package post

import (
	"errors"
	"math/rand"

	"github.com/andycai/weapi/administrator/components/category"
	"github.com/andycai/weapi/administrator/components/config"
	"github.com/andycai/weapi/enum"
	"github.com/andycai/weapi/model"
	"gorm.io/gorm"
)

func NewRenderContentFromPost(post *model.Post, relations bool) *model.RenderContent {
	r := &model.RenderContent{
		BaseContent: post.BaseContent,
		ID:          post.ID,
		SiteID:      post.SiteID,
		PostBody:    post.Body,
		IsDraft:     post.IsDraft,
		Category:    category.NewRenderCategory(post.CategoryID, post.CategoryPath),
	}

	if relations {
		relationCount := config.GetIntValue(enum.KEY_CMS_RELATION_COUNT, 3)
		suggestionCount := config.GetIntValue(enum.KEY_CMS_SUGGESTION_COUNT, 3)

		r.Relations, _ = GetRelations(post.SiteID, post.CategoryID, post.CategoryPath, post.ID, relationCount)
		r.Suggestions, _ = GetSuggestions(post.SiteID, post.CategoryID, post.CategoryPath, post.ID, suggestionCount)
	}
	return r
}

func QueryTags() ([]string, error) {
	var vals []string
	r := db.Model(&model.Post{}).Select("DISTINCT(tags)").Find(&vals)
	return vals, r.Error
}

func GetSuggestions(siteId, categoryId, categoryPath, postId string, maxCount int) ([]model.RelationContent, error) {
	return GetRelations(siteId, categoryId, categoryPath, postId, maxCount)
}

func GetRelations(siteId, categoryId, categoryPath, postId string, maxCount int) ([]model.RelationContent, error) {
	if maxCount <= 0 {
		return nil, nil
	}
	var r []model.RelationContent
	tx := db.Model(&model.Post{}).Where("site_id", siteId).Where("published", true)
	if categoryId != "" {
		tx = tx.Where("category_id", categoryId)
	}

	var totalCount int64
	tx.Count(&totalCount)
	if totalCount == 0 {
		return nil, nil
	}
	excludeIds := []string{}
	if postId != "" {
		excludeIds = append(excludeIds, postId)
	}
	for i := 0; i < maxCount; i++ {
		// random select
		offset := rand.Intn(int(totalCount))
		var val model.Post
		subTx := tx
		if len(excludeIds) > 0 {
			subTx = subTx.Where("id NOT IN (?)", excludeIds)
		}
		result := subTx.Offset(offset).Limit(1).Take(&val)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				continue
			}
			return nil, result.Error
		}

		excludeIds = append(excludeIds, val.ID)

		r = append(r, model.RelationContent{
			BaseContent: val.BaseContent,
			SiteID:      val.SiteID,
			ID:          val.ID})
	}
	return r, nil
}
