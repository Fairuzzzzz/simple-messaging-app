package router

import (
	"github.com/Fairuzzzzz/fiber-boostrap/app/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

type ApiRouter struct{}

func (h ApiRouter) InstallRouter(app *fiber.App) {
	api := app.Group("/api", limiter.New())
	api.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Hello from api",
		})
	})

	userGroup := app.Group("/user")
	userGroupV1 := userGroup.Group("/v1")
	userGroupV1.Post("/register", controllers.Register)
	userGroupV1.Post("/login", controllers.Login)
	userGroupV1.Delete("/logout", MiddlewareValidateAuth, controllers.Logout)
	userGroupV1.Put("/refresh-token", MiddlewareRefreshToken, controllers.RefreshToken)

	messageGroup := app.Group("/message")
	messageGroupV1 := messageGroup.Group("/v1")
	messageGroupV1.Get("/history", MiddlewareValidateAuth, controllers.GetHistory)
}

func NewApiRouter() *ApiRouter {
	return &ApiRouter{}
}
