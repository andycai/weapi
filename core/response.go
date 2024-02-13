package core

import (
	"github.com/andycai/weapi/constant"
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
		"code": constant.Success,
		"data": data,
	})
}

func Push(c *fiber.Ctx, code int) error {
	return c.JSON(fiber.Map{
		"code": code,
		"msg":  constant.CodeText(code),
	})
}

func Err(c *Ctx, statusCode, code int) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"code":  code,
		"error": constant.CodeText(code),
	})
}

func Error(c *Ctx, statusCode int, err error) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"error": err.Error(),
	})
}
