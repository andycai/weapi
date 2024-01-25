package core

import (
	"github.com/andycai/weapi/enum"
	"github.com/gofiber/fiber/v2"
)

// Msg push common response
func Msg(c *Ctx, code int, msg string) error {
	return c.JSON(fiber.Map{
		"code": code,
		"msg":  msg,
	})
}

func Ok(c *fiber.Ctx, data interface{}) error {
	return c.JSON(fiber.Map{
		"code": enum.Success,
		"data": data,
	})
}

func Push(c *fiber.Ctx, code int) error {
	return c.JSON(fiber.Map{
		"code": code,
		"msg":  enum.CodeText(code),
	})
}

func Err(c *Ctx, code int) error {
	return c.Status(code).JSON(fiber.Map{
		"code":  code,
		"error": enum.CodeText(code),
	})
}
