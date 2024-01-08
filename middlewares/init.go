package middlewares

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func Use(app *fiber.App) {
	// 日志
	// app.Use(func(c *fiber.Ctx) error {
	// 	// Log each request
	// 	log.Info(
	// 		"fetch URL",
	// 		zap.String("method", c.Method()),
	// 		zap.String("path", c.Path()),
	// 	)

	// 	// Go to next middleware
	// 	return c.Next()
	// })

	// 限流
	app.Use(
		limiter.New(limiter.Config{
			Next: func(c *fiber.Ctx) bool {
				return c.IP() == "127.0.0.1"
			},
			Max:        300,
			Expiration: 1 * time.Minute,
			// KeyGenerator: func(c *fiber.Ctx) string {
			// 	return c.Get("x-forwarded-for")
			// },
			LimitReached: func(c *fiber.Ctx) error {
				fmt.Println("==============")
				return c.SendString("too fast")
				// return c.SendFile("./toofast.html")
			},
			// Storage: customStarage{}
		}),
	)

	app.Use(
		recover.New(),
		cors.New(),
		requestid.New(),
		logger.New(logger.Config{
			Format:     "${time} ${pid} ${locals:requestid} ${status} - ${method} ${path}​\n​",
			TimeFormat: "2006-01-02 15:04:05",
			// TimeZone:   "America/New_York",
		}),
	)
}
