package page

import (
	"github.com/andycai/weapi/lib/database"
	"github.com/andycai/weapi/model"
)

type PageDao struct {
}

var Dao = new(PageDao)

func (ad PageDao) Count() int64 {
	var page model.Page
	var count int64
	db.Model(&page).Count(&count)

	return count
}

func (ad PageDao) CountTrash() int64 {
	var page model.Page
	var count int64
	db.Model(&page).Unscoped().Where("deleted_at IS NOT NULL").Count(&count)

	return count
}

func (pd PageDao) GetBySlug(slug string) (*model.Page, error) {
	var page model.Page
	err := db.Model(&page).
		Where("slug = ?", slug).
		Find(&page).Error

	return &page, err
}

func (pd PageDao) GetByID(id uint) (*model.Page, error) {
	var page model.Page
	err := db.Model(&page).
		Where("id = ?", id).
		Find(&page).Error

	return &page, err
}

func (pd PageDao) GetAllByPage(page, numPerPage int) []model.Page {
	return pd.GetListByPage(page, numPerPage, "")
}

func (pd PageDao) GetListByPage(page, numPerPage int, q string) []model.Page {
	var pages []model.Page
	tx := db.Model(&pages).
		Preload("User").
		Limit(numPerPage)
	tx = database.DecorateLike(tx, "title", q)
	tx.Offset(page * numPerPage).
		Order("created_at desc").
		Find(&pages)

	return pages
}

func (pd PageDao) GetTrashListByPage(page, numPerPage int, q string) []model.Page {
	var pages []model.Page
	tx := db.Model(&pages).
		Preload("User").
		Unscoped().
		Where("deleted_at IS NOT NULL")
	tx = database.DecorateLike(tx, "title", q)
	tx.Limit(numPerPage).
		Offset(page * numPerPage).
		Order("created_at desc").
		Find(&pages)

	return pages
}

func (pd PageDao) DeleteByIds(ids []uint) {
	var page model.Page
	db.Where(ids).Delete(&page)
}

func (pd PageDao) DeletePermanetlyByIds(ids []uint) {
	var page model.Page
	db.Unscoped().Where(ids).Delete(&page)
}

func (pd PageDao) RestoreByIds(ids []uint) {
	var page model.Page
	db.Unscoped().Model(&page).Where("id IN ?", ids).Update("deleted_at", nil)
}
