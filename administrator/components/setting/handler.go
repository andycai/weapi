package setting

import (
	"errors"
	"fmt"

	"github.com/andycai/weapi/components/user"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/library/authentication"
	"github.com/andycai/weapi/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func BlogPage(c *fiber.Ctx) error {
	var userVo *model.User
	isAuthenticated, userID := authentication.AuthGet(c)

	if isAuthenticated {
		userVo = user.Dao.GetByID(userID)
	}

	var blogVo model.Blog
	err := db.Model(blogVo).Where("user_id= ?", userID).First(&blogVo).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		blogVo = model.Blog{}
	}

	return core.Render(c, "admin/settings/blog", fiber.Map{
		"PageTitle":    "Blog",
		"NavBarActive": "settings",
		"Path":         "/admin/settings/blog",
		"UserName":     userVo.Name,
		"Blog":         blogVo,
		"Info": fiber.Map{
			"BlogName":     "Werite",
			"BlogSubTitle": "Content Management System",
			"LoginAt":      userVo.LoginAt,
		},
	}, "admin/layouts/app")
}

func BlogSave(c *fiber.Ctx) error {
	blogVo := model.Blog{}

	_, userID := authentication.AuthGet(c)

	err := user.BindBlog(c, &blogVo)
	if err != nil {
		return err
	}

	blogVo.UserID = userID

	err = db.Model(&blogVo).Where("user_id= ?", userID).First(&blogVo).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		db.Create(&blogVo)
	} else {
		db.Model(&blogVo).Updates(map[string]interface{}{"name": blogVo.Name, "description": blogVo.Description})
	}

	core.PushMessages(fmt.Sprintf("Updated blog infomation"))

	return c.Redirect("/admin/settings/blog")
}
