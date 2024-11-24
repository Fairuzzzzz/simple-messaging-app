package router

import (
	"log"
	"time"

	"github.com/Fairuzzzzz/fiber-boostrap/app/repository"
	"github.com/Fairuzzzzz/fiber-boostrap/pkg/jwt"
	"github.com/Fairuzzzzz/fiber-boostrap/pkg/response"
	"github.com/gofiber/fiber/v2"
	"go.elastic.co/apm"
)

func MiddlewareValidateAuth(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "MiddlewareValidateAuth", "router")
	defer span.End()

	auth := ctx.Get("Authorization")
	if auth == "" {
		log.Println("Authorization is empty")
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	_, err := repository.GetUserSession(spanCtx, auth)
	if err != nil {
		log.Println("failed to get user session on database: ", err)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	claim, err := jwt.ValidateToken(spanCtx, auth)
	if err != nil {
		log.Println(err)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	if time.Now().Unix() > claim.ExpiresAt.Unix() {
		log.Println("jwt token expired: ", claim.ExpiresAt)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	ctx.Locals("username", claim.Username)
	ctx.Locals("full_name", claim.Fullname)

	return ctx.Next()
}

func MiddlewareRefreshToken(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "MiddlewareRefreshToken", "router")
	defer span.End()

	auth := ctx.Get("Authorization")
	if auth == "" {
		log.Println("Authorization is empty")
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	claim, err := jwt.ValidateToken(spanCtx, auth)
	if err != nil {
		log.Println(err)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	if time.Now().Unix() > claim.ExpiresAt.Unix() {
		log.Println("jwt token expired: ", claim.ExpiresAt)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	ctx.Locals("username", claim.Username)
	ctx.Locals("full_name", claim.Fullname)

	return ctx.Next()
}
