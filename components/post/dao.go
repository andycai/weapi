package post

import (
	"github.com/andycai/weapi/lib/database"
	"github.com/andycai/weapi/model"
	"gorm.io/gorm"
)

type PostDao struct{}

var Dao = new(PostDao)

//#region Post

func (pd PostDao) GetBySlug(slug string) (*model.Post, error) {
	var post model.Post
	err := db.Model(&post).
		Where("slug = ?", slug).
		Find(&post).Error

	return &post, err
}

func (pd PostDao) GetByID(id uint) (*model.Post, error) {
	var post model.Post
	err := db.Model(&post).
		Where("id = ?", id).
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Order("tags.name asc")
		}).
		Find(&post).Error

	return &post, err
}

func (pd PostDao) CountByPublished() int64 {
	return pd.countByDraft(0)
}

func (pd PostDao) CountByDraft() int64 {
	return pd.countByDraft(1)
}

func (pd PostDao) countByDraft(draft int) int64 {
	var post model.Post
	var count int64
	db.Model(&post).Where("is_draft = ?", draft).Count(&count)

	return count
}

func (pd PostDao) CountByTrash() int64 {
	var post model.Post
	var count int64
	db.Model(&post).Unscoped().Where("deleted_at IS NOT NULL").Count(&count)

	return count
}

func (pd PostDao) Count() int64 {
	var post model.Post
	var count int64
	db.Model(&post).Count(&count)

	return count
}

func (pd PostDao) GetAllByPage(page, numPerPage int) []model.Post {
	return pd.GetListByPage(page, numPerPage, 0, "")
}

func (pd PostDao) GetListByPage(page, numPerPage int, categoryID int, q string) []model.Post {
	var posts []model.Post
	tx := db.Model(&posts).
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Order("tags.name asc")
		}).
		Preload("User").
		Preload("Category")
	tx = database.DecorateEqualInt(tx, "category_id", categoryID)
	tx = database.DecorateLike(tx, "title", q)
	tx.Limit(numPerPage).
		Offset(page * numPerPage).
		Order("created_at desc").
		Find(&posts)

	return posts
}

func (pd PostDao) GetPublishedListByPage(page, numPerPage int, categoryID int, q string) []model.Post {
	var posts []model.Post
	tx := db.Model(&posts).
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Order("tags.name asc")
		}).
		Preload("User").
		Preload("Category").
		Where("is_draft = ?", 0)
	tx = database.DecorateEqualInt(tx, "category_id", categoryID)
	tx = database.DecorateLike(tx, "title", q)
	tx.Limit(numPerPage).
		Offset(page * numPerPage).
		Order("created_at desc").
		Find(&posts)

	return posts
}

func (pd PostDao) GetDraftListByPage(page, numPerPage int, categoryID int, q string) []model.Post {
	var posts []model.Post
	tx := db.Model(&posts).
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Order("tags.name asc")
		}).
		Preload("User").
		Preload("Category").
		Where("is_draft = ?", 1)
	tx = database.DecorateEqualInt(tx, "category_id", categoryID)
	tx = database.DecorateLike(tx, "title", q)
	tx.Limit(numPerPage).
		Offset(page * numPerPage).
		Order("created_at desc").
		Find(&posts)

	return posts
}

func (pd PostDao) GetTrashListByPage(page, numPerPage int, categoryID int, q string) []model.Post {
	var posts []model.Post
	tx := db.Model(&posts).
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Order("tags.name asc")
		}).
		Preload("User").
		Preload("Category").
		Unscoped().
		Where("deleted_at IS NOT NULL")
	tx = database.DecorateEqualInt(tx, "category_id", categoryID)
	tx = database.DecorateLike(tx, "title", q)
	tx.Limit(numPerPage).
		Offset(page * numPerPage).
		Order("created_at desc").
		Find(&posts)

	return posts
}

func (pd PostDao) UpdateDraft(value uint) {
	db.Model(&model.Post{}).Update("is_draft", value)
}

func (pd PostDao) DeleteByIds(ids []uint) {
	var post model.Post
	db.Where(ids).Delete(&post)
}

func (pd PostDao) DeletePermanetlyByIds(ids []uint) {
	var post model.Post
	db.Unscoped().Where(ids).Delete(&post)
}

func (pd PostDao) RestoreByIds(ids []uint) {
	var post model.Post
	db.Unscoped().Model(&post).Where("id IN ?", ids).Update("deleted_at", nil)
}

//#endregion

//#region Category

func (pd PostDao) CountCatgegory() int64 {
	var category model.Category
	var count int64
	db.Model(&category).Count(&count)

	return count
}

func (pd PostDao) GetCategories() []model.Category {
	var categories []model.Category
	db.Model(&categories).
		Find(&categories)

	return categories
}

func (pd PostDao) GetCategoriesByPage(page, numPerPage int) []model.Category {
	var categories []model.Category
	db.Model(&categories).
		Limit(numPerPage).
		Offset(page * numPerPage).
		Find(&categories)

	return categories
}

func (pd PostDao) GetCategoryByID(id uint) (*model.Category, error) {
	var categoryVo model.Category
	err := db.Model(&categoryVo).Where("id = ?", id).Find(&categoryVo).Error

	return &categoryVo, err
}

func (pd PostDao) DeleteCategoriesByIds(ids []uint) {
	var category model.Category
	// categories := make([]model.Category, len(ids))
	// for i, v := range ids {
	// 	categories[i].ID = v
	// }
	db.Where(ids).Delete(&category)
}
