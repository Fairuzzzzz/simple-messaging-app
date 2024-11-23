package controllers

import (
	"log"

	"github.com/Fairuzzzzz/fiber-boostrap/app/repository"
	"github.com/Fairuzzzzz/fiber-boostrap/pkg/response"
	"github.com/gofiber/fiber/v2"
)

func GetHistory(ctx *fiber.Ctx) error {
	resp, err := repository.GetAllMessage(ctx.Context())
	if err != nil {
		log.Println(err)
		return response.SendFailureResponse(
			ctx,
			fiber.StatusInternalServerError,
			"error internal server",
			nil,
		)
	}
	return response.SendSuccessResponse(ctx, resp)
}
