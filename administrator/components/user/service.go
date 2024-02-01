package user

import (
	"errors"
	"strings"
	"time"

	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/lib/authentication"
	"github.com/andycai/weapi/model"
	"github.com/gofiber/fiber/v2"
)

func SuperAccessCheck(c *fiber.Ctx, obj *model.AdminObject) error {
	isAuthenticated, _ := authentication.AuthGet(c)
	if isAuthenticated {
		return nil
	}
	return errors.New("not authorized")
}

func GetByID(id uint) *model.User {
	var user model.User
	db.Model(&user).
		Where("id = ?", id).
		First(&user)

	return &user
}

func CurrentUser(c *fiber.Ctx) *model.User {
	var userVo *model.User
	isAuthenticated, userID := authentication.AuthGet(c)

	if isAuthenticated {
		userVo = GetByID(userID)
	}

	return userVo
}

func GetByEmail(email string) (error, *model.User) {
	var user model.User
	result := db.Where("email", strings.ToLower(email)).Take(&user)

	return result.Error, &user
}

func CreateUser(user *model.User) error {
	result := db.Create(user)

	return result.Error
}

func UpdatePassword(user *model.User, password string) error {
	p := core.HashPassword(password)
	err := UpdateFields(user, map[string]any{
		"Password": p,
	})
	if err != nil {
		return err
	}
	user.Password = p

	return err
}

func UpdateLoginTime(userID uint) error {
	db.Model(&model.User{}).Where("id = ?", userID).Update("last_login", time.Now())

	return nil
}

func UpdateFields(user *model.User, vals map[string]any) error {
	return db.Model(user).Updates(vals).Error
}