package user

import (
	"time"

	adminuser "github.com/andycai/weapi/administrator/components/user"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/enum"
	"github.com/andycai/weapi/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func handleLogin(c *fiber.Ctx) error {
	loginVo := &adminuser.ReqLogin{}
	if err := c.BodyParser(&loginVo); err != nil {
		return core.Err(c, enum.ErrUserEmailOrPasswordError)
	}
	user := c.FormValue("user")
	pass := c.FormValue("pass")

	// Throws Unauthorized error
	if user != "john" || pass != "doe" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"name":  "John Doe",
		"admin": true,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(utils.GetEnv(enum.ENV_SESSION_SECRET)))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}
