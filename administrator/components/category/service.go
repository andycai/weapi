package category

import (
	"fmt"
	"strings"

	"github.com/andycai/weapi/model"
	"gorm.io/gorm"
)

func QueryCategoryWithCount(siteId, contentObject string) ([]model.Category, error) {
	var tx *gorm.DB
	switch contentObject {
	case "post":
		tx = db.Model(&model.Post{}).Where("site_id", siteId)
	case "page":
		tx = db.Model(&model.Page{}).Where("site_id", siteId)
	default:
		return nil, fmt.Errorf("invalid content object: %s", contentObject)
	}

	var vals []model.Category
	r := db.Model(&model.Category{}).Where("site_id", siteId).Find(&vals)
	if r.Error != nil {
		return nil, r.Error
	}
	for i := range vals {
		val := &vals[i]
		tx := tx.Where("category_id", val.UUID)
		var count int64
		tx.Count(&count)
		val.Count = int(count)
	}
	return vals, r.Error
}

func NewRenderCategory(categoryID, categoryPath string) *model.RenderCategory {
	var category model.Category
	r := db.Model(&model.Category{}).Where("uuid", categoryID).First(&category)
	if r.Error != nil {
		return nil
	}

	selected := category.FindItem(categoryPath)

	obj := &model.RenderCategory{
		UUID: category.UUID,
		Name: category.Name,
	}
	if selected != nil {
		obj.Path = selected.Path
		obj.PathName = selected.Name
	}
	return obj
}

// Query tags by category
func GetTagsByCategory(contentType string, form *model.TagsForm) ([]string, error) {
	var tx *gorm.DB
	switch contentType {
	case "post":
		tx = db.Model(&model.Post{})
	case "page":
		tx = db.Model(&model.Page{})
	}

	if form.SiteId != "" {
		tx = tx.Where("site_id", form.SiteId)
	}

	if form.CategoryId != "" {
		tx = tx.Where("category_id", form.CategoryId)
	}

	if form.CategoryPath != "" {
		tx = tx.Where("category_path", form.CategoryPath)
	}

	var rawTags []string
	r := tx.Pluck("tags", &rawTags)
	if r.Error != nil {
		return nil, r.Error
	}

	var uniqueTags map[string]string = make(map[string]string)
	for _, tag := range rawTags {
		if tag == "" {
			continue
		}
		vals := strings.Split(tag, ",")
		for _, val := range vals {
			val = strings.TrimSpace(val)
			if val == "" {
				continue
			}
			uniqueTags[strings.ToLower(val)] = val
		}
	}

	var tags []string = make([]string, 0, len(uniqueTags))
	for k, v := range uniqueTags {
		if k == "" {
			continue
		}
		tags = append(tags, v)
	}
	return tags, r.Error
}

func QueryContentByTags(contentType string, form *model.QueryByTagsForm) ([]string, error) {
	return nil, nil
}
