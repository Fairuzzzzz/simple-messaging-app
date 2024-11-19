package router

import (
	"github.com/Fairuzzzzz/fiber-boostrap/app/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
)

type HttpRouter struct{}

func (h HttpRouter) InstallRouter(app *fiber.App) {
	group := app.Group("", cors.New(), csrf.New())
	group.Get("/", controllers.RenderUi)
}

func NewHttpRouter() *HttpRouter {
	return &HttpRouter{}
}
