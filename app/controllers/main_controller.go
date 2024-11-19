package controllers

import "github.com/gofiber/fiber/v2"

func RenderUi(c *fiber.Ctx) error {
	return c.Render("index", nil)
}
