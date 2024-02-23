package user

import (
	"errors"
	"net/http"
	"time"

	"github.com/andycai/weapi/constant"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/lib/authentication"
	"github.com/andycai/weapi/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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
		"copyright":   "2024 Weapi",
	}, "layout/app")
}

func handleSiginAction(c *fiber.Ctx) error {
	loginVo := &ReqLogin{}

	if err := c.BodyParser(&loginVo); err != nil {
		return core.Err(c, http.StatusBadRequest, constant.ErrUserEmailOrPasswordError)
	}

	if loginVo.Email == "" || loginVo.Password == "" {
		return core.Err(c, http.StatusBadRequest, constant.ErrUserEmailOrPasswordError)
	}

	err, userVo := GetByEmail(loginVo.Email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return core.Err(c, http.StatusForbidden, constant.ErrUserEmailOrPasswordError)
		}
	}

	if !core.CheckPassword(userVo.Password, loginVo.Password) {
		return core.Err(c, http.StatusForbidden, constant.ErrUserEmailOrPasswordError)
	}

	if !userVo.Enabled {
		return core.Err(c, http.StatusForbidden, constant.ErrUserDisabled)
	}

	if !userVo.Activated {
		return core.Err(c, http.StatusForbidden, constant.ErrUserNotActivated)
	}

	UpdateLogin(c, userVo.ID)
	authentication.AuthStore(c, userVo.ID)

	if !loginVo.Remember {
		core.Push(c, constant.Success)
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"name":  userVo.Email,
		"admin": userVo.IsSuperUser,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(utils.GetEnv(constant.ENV_SESSION_SECRET)))
	if err != nil {
		return core.Error(c, http.StatusInternalServerError, err)
	}

	return c.JSON(fiber.Map{"token": t})
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
