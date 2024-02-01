package user

import (
	"errors"

	"github.com/andycai/weapi/components/user"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/enum"
	"github.com/andycai/weapi/lib/authentication"
	"github.com/andycai/weapi/model"
	"github.com/andycai/weapi/object"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SuperAccessCheck(c *fiber.Ctx, obj *object.AdminObject) error {
	isAuthenticated, _ := authentication.AuthGet(c)
	if isAuthenticated {
		return nil
	}
	return errors.New("not authorized")
}

func SigninPage(c *fiber.Ctx) error {
	isAuthenticated, _ := authentication.AuthGet(c)

	if isAuthenticated {
		return c.Redirect("/admin/")
	}

	return core.Render(c, "signin", fiber.Map{
		"signup_url":  "/auth/register",
		"signuptext":  "Sign up",
		"login_next":  "/admin/",
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
		return c.Redirect("/auth/login/")
	}

	user.Dao.UpdateLogoutTime(userID)
	authentication.AuthDestroy(c)

	return c.Redirect("/auth/login/")
}

func SignupPage(c *fiber.Ctx) error {
	return nil
}

func SignupAction(c *fiber.Ctx) error {
	return nil
}

func CurrentUser(c *fiber.Ctx) *model.User {
	var userVo *model.User
	isAuthenticated, userID := authentication.AuthGet(c)

	if isAuthenticated {
		userVo = user.Dao.GetByID(userID)
	}

	return userVo
}
