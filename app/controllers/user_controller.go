package controllers

import (
	"fmt"
	"time"

	"github.com/Fairuzzzzz/fiber-boostrap/app/models"
	"github.com/Fairuzzzzz/fiber-boostrap/app/repository"
	"github.com/Fairuzzzzz/fiber-boostrap/pkg/jwt"
	"github.com/Fairuzzzzz/fiber-boostrap/pkg/response"
	"github.com/gofiber/fiber/v2"
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

func Login(ctx *fiber.Ctx) error {
	loginReq := new(models.LoginRequest)
	loginRes := models.LoginResponse{}

	err := ctx.BodyParser(loginReq)
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

	err = loginReq.Validate()
	if err != nil {
		errResponse := fmt.Errorf("failed to validate login request: %v", err)
		fmt.Println(errResponse)
		return response.SendFailureResponse(
			ctx,
			fiber.StatusInternalServerError,
			errResponse.Error(),
			nil,
		)
	}

	user, err := repository.GetUserByUsername(ctx.Context(), loginReq.Username)
	if err != nil {
		errResponse := fmt.Errorf("failed to get username: %v", err)
		fmt.Println(errResponse)
		return response.SendFailureResponse(
			ctx,
			fiber.StatusNotFound,
			"username/password is invalid",
			nil,
		)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if err != nil {
		errResponse := fmt.Errorf("failed to get username: %v", err)
		fmt.Println(errResponse)
		return response.SendFailureResponse(
			ctx,
			fiber.StatusNotFound,
			"username/password is invalid",
			nil,
		)
	}

	token, err := jwt.GenerateToken(ctx.Context(), user.Username, user.FullName, "token")
	if err != nil {
		errResponse := fmt.Errorf("failed to generate token: %v", err)
		fmt.Println(errResponse)
		return response.SendFailureResponse(
			ctx,
			fiber.StatusNotFound,
			errResponse.Error(),
			nil,
		)
	}

	refreshToken, err := jwt.GenerateToken(
		ctx.Context(),
		user.Username,
		user.FullName,
		"refresh_token",
	)
	if err != nil {
		errResponse := fmt.Errorf("failed to generate token: %v", err)
		fmt.Println(errResponse)
		return response.SendFailureResponse(
			ctx,
			fiber.StatusNotFound,
			errResponse.Error(),
			nil,
		)
	}

	userSession := &models.UserSession{
		UserID:              int(user.ID),
		Token:               token,
		RefreshToken:        refreshToken,
		TokenExpired:        time.Now().Add(jwt.MapTokenType["token"]),
		RefreshTokenExpired: time.Now().Add(jwt.MapTokenType["refresh_token"]),
	}

	err = repository.InsertNewUserSession(ctx.Context(), userSession)
	if err != nil {
		errResponse := fmt.Errorf("failed to insert user session: %v", err)
		fmt.Println(errResponse)
		return response.SendFailureResponse(
			ctx,
			fiber.StatusNotFound,
			errResponse.Error(),
			nil,
		)
	}

	loginRes.Username = user.Username
	loginRes.FullName = user.FullName
	loginRes.Token = token
	loginRes.RefreshToken = refreshToken

	return response.SendSuccessResponse(ctx, loginRes)
}
