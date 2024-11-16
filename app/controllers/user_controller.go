package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/kooroshh/fiber-boostrap/app/models"
	"github.com/kooroshh/fiber-boostrap/app/repository"
	"github.com/kooroshh/fiber-boostrap/pkg/response"
	"golang.org/x/crypto/bcrypt"
)

func Register(ctx *fiber.Ctx) error {
	user := new(models.User)
	err := ctx.BodyParser(user)
	if err != nil {
		errResponse := fmt.Errorf("failed to parsing request: %v", err)
		fmt.Println(errResponse)
		return response.SendFailureResponse(
			ctx,
			fiber.StatusInternalServerError,
			errResponse.Error(),
			nil,
		)
	}

	err = user.Validate()
	if err != nil {
		errResponse := fmt.Errorf("failed to validate user: %v", err)
		fmt.Println(errResponse)
		return response.SendFailureResponse(
			ctx,
			fiber.StatusInternalServerError,
			errResponse.Error(),
			nil,
		)
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		errResponse := fmt.Errorf("failed to encrypt password: %v", err)
		fmt.Println(errResponse)
		return response.SendFailureResponse(
			ctx,
			fiber.StatusInternalServerError,
			errResponse.Error(),
			nil,
		)
	}

	user.Password = string(hashPassword)

	err = repository.InsertNewUser(ctx.Context(), user)
	if err != nil {
		errResponse := fmt.Errorf("failed to insert new user: %v", err)
		fmt.Println(errResponse)
		return response.SendFailureResponse(
			ctx,
			fiber.StatusInternalServerError,
			errResponse.Error(),
			nil,
		)
	}

	userResponse := models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		FullName: user.FullName,
	}

	return response.SendSuccessResponse(ctx, userResponse)
}
