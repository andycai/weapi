package user

import (
	"net/http"

	"github.com/andycai/weapi/administrator/components/config"
	"github.com/andycai/weapi/administrator/components/user"
	"github.com/andycai/weapi/enum"
	"github.com/andycai/weapi/utils"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func WithAPIAuth(c *fiber.Ctx) error {
	if user.Current(c) != nil {
		return c.Next()
	}
	guestAccess := config.GetBoolValue(enum.KEY_CMS_GUEST_ACCESS_API)
	if guestAccess {
		switch c.Method() {
		case http.MethodGet, http.MethodHead, http.MethodPost, http.MethodOptions:
			return c.Next()
		}
	}

	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(utils.GetEnv(enum.ENV_SESSION_SECRET))},
	})(c)
}
