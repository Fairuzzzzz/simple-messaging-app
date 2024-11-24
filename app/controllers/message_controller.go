package controllers

import (
	"log"

	"github.com/Fairuzzzzz/fiber-boostrap/app/repository"
	"github.com/Fairuzzzzz/fiber-boostrap/pkg/response"
	"github.com/gofiber/fiber/v2"
	"go.elastic.co/apm"
)

func GetHistory(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "GetHistory", "controller")
	defer span.End()

	resp, err := repository.GetAllMessage(spanCtx)
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
