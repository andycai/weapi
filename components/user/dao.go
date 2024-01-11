package user

import (
	"strings"
	"time"

	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/model"
)

type UserDao struct{}

var Dao = new(UserDao)

func (ud UserDao) GetByID(id uint) *model.User {
	var user model.User
	db.Model(&user).
		Where("id = ?", id).
		First(&user)

	return &user
}

func (ud UserDao) Count() int64 {
	var user model.User
	var count int64
	db.Model(&user).Count(&count)

	return count
}

func (ud UserDao) GetByEmail(email string) (error, *model.User) {
	var user model.User
	result := db.Where("email", strings.ToLower(email)).Take(&user)

	return result.Error, &user
}

func (ud UserDao) Create(user *model.User) error {
	result := db.Create(user)

	return result.Error
}

func (ud UserDao) UpdatePassword(user *model.User, password string) error {
	p := core.HashPassword(password)
	err := ud.UpdateFields(user, map[string]any{
		"Password": p,
	})
	if err != nil {
		return err
	}
	user.Password = p

	return err
}

func (ud UserDao) UpdateLoginTime(userID uint) error {
	db.Model(&model.User{}).Where("id = ?", userID).Update("login_at", time.Now())

	return nil
}

func (ud UserDao) UpdateLogoutTime(userID uint) error {
	db.Model(&model.User{}).Where("id = ?", userID).Update("logout_at", time.Now())

	return nil
}

func (ud UserDao) UpdateFields(user *model.User, vals map[string]any) error {
	return db.Model(user).Updates(vals).Error
}
