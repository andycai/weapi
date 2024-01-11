package user

import (
	"errors"
	"fmt"
	"time"

	"github.com/andycai/weapi/components/page"
	"github.com/andycai/weapi/components/post"
	"github.com/andycai/weapi/components/user"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/library/authentication"
	"github.com/andycai/weapi/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func LoginPage(c *fiber.Ctx) error {
	isAuthenticated, _ := authentication.AuthGet(c)

	if isAuthenticated {
		return c.Redirect("/admin/dashboard")
	}

	return core.Render(c, "admin/login", fiber.Map{})
}

func LoginAction(c *fiber.Ctx) error {
	var userVo *model.User
	email := c.FormValue("email")
	password := c.FormValue("password")

	if email == "" || password == "" {
		return core.Render(c, "admin/login", fiber.Map{
			"Errors": []string{
				"Email or password cannot be null.",
			},
		})
	}

	err, userVo := user.Dao.GetByEmail(email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return core.Render(c, "admin/login", fiber.Map{
				"Errors": []string{
					"Email and password did not match.",
				},
			})
		}
	}

	if !core.CheckPassword(userVo.Password, password) {
		return core.Render(c, "admin/login", fiber.Map{
			"Errors": []string{
				"Email and password did not match.",
			},
		})
	}

	user.Dao.UpdateLoginTime(uint(userVo.ID))
	authentication.AuthStore(c, uint(userVo.ID))

	return c.Redirect("/admin/dashboard")
}

func LogoutAction(c *fiber.Ctx) error {
	isAuthenticated, userID := authentication.AuthGet(c)
	if !isAuthenticated {
		return core.Render(c, "admin/login", fiber.Map{})
	}

	user.Dao.UpdateLogoutTime(userID)
	authentication.AuthDestroy(c)

	return core.Render(c, "admin/login", fiber.Map{})
}

func DashBoardPage(c *fiber.Ctx) error {
	var userVo *model.User
	isAuthenticated, userID := authentication.AuthGet(c)

	userTotal := user.Dao.Count()
	postTotal := post.Dao.Count()
	pageTotal := page.Dao.Count()

	name := ""
	loginAt := time.Now()
	if isAuthenticated {
		userVo = user.Dao.GetByID(userID)
		name = userVo.Name
		loginAt = userVo.LoginAt
	}

	return core.Render(c, "admin/dashboard", fiber.Map{
		"PageTitle":    "DashBoard",
		"NavBarActive": "dashboard",
		"Path":         "/admin/dashboard",
		"UserName":     name,
		"UserTotal":    userTotal,
		"PostTotal":    postTotal,
		"PageTotal":    pageTotal,
		"Info": fiber.Map{
			"BlogName":     "Werite",
			"BlogSubTitle": "Content Management System",
			"LoginAt":      loginAt,
		},
	}, "admin/layouts/app")
}

func ProfilePage(c *fiber.Ctx) error {
	var userVo *model.User
	isAuthenticated, userID := authentication.AuthGet(c)

	if isAuthenticated {
		userVo = user.Dao.GetByID(userID)
	}

	return core.Render(c, "admin/users/profile", fiber.Map{
		"PageTitle":    "Profile",
		"NavBarActive": "users",
		"Path":         "/admin/users/profile",
		"UserName":     userVo.Name,
		"User":         userVo,
		"Info": fiber.Map{
			"BlogName":     "Werite",
			"BlogSubTitle": "Content Management System",
			"LoginAt":      userVo.LoginAt,
		},
	}, "admin/layouts/app")
}

func SecurityPage(c *fiber.Ctx) error {
	var userVo *model.User
	isAuthenticated, userID := authentication.AuthGet(c)

	if isAuthenticated {
		userVo = user.Dao.GetByID(userID)
	}

	return core.Render(c, "admin/users/security", fiber.Map{
		"PageTitle":    "Security",
		"NavBarActive": "users",
		"Path":         "/admin/users/security",
		"UserName":     userVo.Name,
		"Info": fiber.Map{
			"BlogName":     "Werite",
			"BlogSubTitle": "Content Management System",
			"LoginAt":      userVo.LoginAt,
		},
	}, "admin/layouts/app")
}

func BlogPage(c *fiber.Ctx) error {
	var userVo *model.User
	isAuthenticated, userID := authentication.AuthGet(c)

	if isAuthenticated {
		userVo = user.Dao.GetByID(userID)
	}

	return core.Render(c, "admin/users/blog", fiber.Map{
		"PageTitle":    "Blog",
		"NavBarActive": "users",
		"Path":         "/admin/users/blog",
		"UserName":     userVo.Name,
		"Info": fiber.Map{
			"BlogName":     "Werite",
			"BlogSubTitle": "Content Management System",
			"LoginAt":      userVo.LoginAt,
		},
	}, "admin/layouts/app")
}

func ProfileSave(c *fiber.Ctx) error {
	var userVo *model.User
	isAuthenticated, userID := authentication.AuthGet(c)

	if isAuthenticated {
		userVo = user.Dao.GetByID(userID)
	}

	err := user.BindProfile(c, userVo)
	if err != nil {
		return err
	}

	db.Model(userVo).Updates(map[string]interface{}{"gender": userVo.Gender, "phone": userVo.Phone, "email": userVo.Email, "addr": userVo.Addr})

	core.PushMessages(fmt.Sprintf("Updated profile"))

	return c.Redirect("/admin/users/profile")
}

func PasswordSave(c *fiber.Ctx) error {
	var userVo *model.User
	isAuthenticated, userID := authentication.AuthGet(c)

	if isAuthenticated {
		userVo = user.Dao.GetByID(userID)
	}

	err := user.BindPassword(c, userVo)
	if err != nil {
		return err
	}

	db.Model(userVo).Update("password", userVo.Password)

	core.PushMessages(fmt.Sprintf("Updated password"))

	return c.Redirect("/admin/users/security")
}
