package user

import (
	"errors"
	"net/http"
	"time"

	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/enum"
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
	}, "layout/app")
}

func handleSiginAction(c *fiber.Ctx) error {
	loginVo := &ReqLogin{}

	if err := c.BodyParser(&loginVo); err != nil {
		return core.Err(c, http.StatusBadRequest, enum.ErrUserEmailOrPasswordError)
	}

	if loginVo.Email == "" || loginVo.Password == "" {
		return core.Err(c, http.StatusBadRequest, enum.ErrUserEmailOrPasswordError)
	}

	err, userVo := GetByEmail(loginVo.Email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return core.Err(c, http.StatusForbidden, enum.ErrUserEmailOrPasswordError)
		}
	}

	if !core.CheckPassword(userVo.Password, loginVo.Password) {
		return core.Err(c, http.StatusForbidden, enum.ErrUserEmailOrPasswordError)
	}

	if !userVo.Enabled {
		return core.Err(c, http.StatusForbidden, enum.ErrUserDisabled)
	}

	if !userVo.Activated {
		return core.Err(c, http.StatusForbidden, enum.ErrUserNotActivated)
	}

	UpdateLogin(c, userVo.ID)
	authentication.AuthStore(c, userVo.ID)

	if !loginVo.Remember {
		core.Push(c, enum.Success)
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
	t, err := token.SignedString([]byte(utils.GetEnv(enum.ENV_SESSION_SECRET)))
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
