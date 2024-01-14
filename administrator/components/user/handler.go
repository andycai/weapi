package user

import (
	"errors"
	"fmt"
	"time"

	"github.com/andycai/weapi/components/page"
	"github.com/andycai/weapi/components/post"
	"github.com/andycai/weapi/components/user"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/enum"
	"github.com/andycai/weapi/library/authentication"
	"github.com/andycai/weapi/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SigninPage(c *fiber.Ctx) error {
	isAuthenticated, _ := authentication.AuthGet(c)

	if isAuthenticated {
		return c.Redirect("/admin/dashboard")
	}

	return core.Render(c, "signin", fiber.Map{
		"signup_url":  "/auth/register",
		"signuptext":  "Sign up",
		"login_next":  "/admin/dashboard",
		"sitename":    "Weapi",
		"logo_url":    "/static/img/logo.svg",
		"favicon_url": "/static/img/favicon.png",
		"title":       "Sign in",
	}, "layout/app")
}

func SigninAction(c *fiber.Ctx) error {
	userVo := &model.User{}

	err := user.BindLogin(c, userVo)
	if err != nil {
		return core.Err(c, enum.ErrUserEmailOrPasswordError)
	}
	email := userVo.Email
	password := userVo.Password

	if email == "" || password == "" {
		return core.Err(c, enum.ErrUserEmailOrPasswordIsEmpty)
	}

	err, userVo = user.Dao.GetByEmail(email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return core.Err(c, enum.ErrUserEmailOrPasswordError)
		}
	}

	if !core.CheckPassword(userVo.Password, password) {
		return core.Err(c, enum.ErrUserEmailOrPasswordError)
	}

	user.Dao.UpdateLoginTime(uint(userVo.ID))
	authentication.AuthStore(c, uint(userVo.ID))

	return core.Push(c, enum.Success)
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

func SignupPage(c *fiber.Ctx) error {
	return nil
}

func SignupAction(c *fiber.Ctx) error {
	return nil
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
		name = userVo.FirstName
		loginAt = *userVo.LastLogin
	}

	return core.Render(c, "admin/index", fiber.Map{
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
	})
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
		"UserName":     userVo.FirstName,
		"User":         userVo,
		"Info": fiber.Map{
			"BlogName":     "Werite",
			"BlogSubTitle": "Content Management System",
			"LoginAt":      userVo.LastLogin,
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
		"UserName":     userVo.FirstName,
		"Info": fiber.Map{
			"BlogName":     "Werite",
			"BlogSubTitle": "Content Management System",
			"LoginAt":      userVo.LastLogin,
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
		"UserName":     userVo.FirstName,
		"Info": fiber.Map{
			"BlogName":     "Werite",
			"BlogSubTitle": "Content Management System",
			"LoginAt":      userVo.LastLogin,
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

	db.Model(userVo).Updates(map[string]interface{}{
		"gender": userVo.Profile.Gender,
		"phone":  userVo.Phone,
		"email":  userVo.Email,
		"addr":   userVo.Profile.City})

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
