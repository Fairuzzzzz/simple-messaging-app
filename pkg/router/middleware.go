package router

import (
	"fmt"
	"time"

	"github.com/Fairuzzzzz/fiber-boostrap/app/repository"
	"github.com/Fairuzzzzz/fiber-boostrap/pkg/jwt"
	"github.com/Fairuzzzzz/fiber-boostrap/pkg/response"
	"github.com/gofiber/fiber/v2"
)

func MiddlewareValidateAuth(ctx *fiber.Ctx) error {
	auth := ctx.Get("Authorization")
	if auth == "" {
		fmt.Println("Authorization is empty")
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	_, err := repository.GetUserSession(ctx.Context(), auth)
	if err != nil {
		fmt.Println("failed to get user session on database: ", err)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	claim, err := jwt.ValidateToken(ctx.Context(), auth)
	if err != nil {
		fmt.Println(err)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	if time.Now().Unix() > claim.ExpiresAt.Unix() {
		fmt.Println("jwt token expired: ", claim.ExpiresAt)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	ctx.Locals("username", claim.Username)
	ctx.Locals("full_name", claim.Fullname)

	return ctx.Next()
}

func MiddlewareRefreshToken(ctx *fiber.Ctx) error {
	auth := ctx.Get("Authorization")
	if auth == "" {
		fmt.Println("Authorization is empty")
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	claim, err := jwt.ValidateToken(ctx.Context(), auth)
	if err != nil {
		fmt.Println(err)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	if time.Now().Unix() > claim.ExpiresAt.Unix() {
		fmt.Println("jwt token expired: ", claim.ExpiresAt)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	ctx.Locals("username", claim.Username)
	ctx.Locals("full_name", claim.Fullname)

	return ctx.Next()
}
