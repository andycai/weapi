package user

import (
	"errors"

	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/enum"
	"github.com/andycai/weapi/lib/authentication"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func handleSigin(c *fiber.Ctx) error {
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

func handleSiginAction(c *fiber.Ctx) error {
	loginVo := &ReqLogin{}

	if err := c.BodyParser(&loginVo); err != nil {
		return core.Err(c, enum.ErrUserEmailOrPasswordError)
	}

	if loginVo.Email == "" || loginVo.Password == "" {
		return core.Err(c, enum.ErrUserEmailOrPasswordError)
	}

	err, userVo := GetByEmail(loginVo.Email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return core.Err(c, enum.ErrUserEmailOrPasswordError)
		}
	}

	if !core.CheckPassword(userVo.Password, loginVo.Password) {
		return core.Err(c, enum.ErrUserEmailOrPasswordError)
	}

	UpdateLoginTime(uint(userVo.ID))
	authentication.AuthStore(c, uint(userVo.ID))

	return core.Push(c, enum.Success)
}

func handleLogoutAction(c *fiber.Ctx) error {
	isAuthenticated, _ := authentication.AuthGet(c)
	if !isAuthenticated {
		return c.Redirect("/auth/login/")
	}

	authentication.AuthDestroy(c)

	return c.Redirect("/auth/login/")
}

func handleSigup(c *fiber.Ctx) error {
	return nil
}

func handleSigupAction(c *fiber.Ctx) error {
	return nil
}
