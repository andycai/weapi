package post

import (
	"errors"
	"math/rand"

	"github.com/andycai/weapi/model"
	"gorm.io/gorm"
)

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
