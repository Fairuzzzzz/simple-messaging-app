package controllers

import (
	"fmt"
	"log"
	"time"

	"github.com/Fairuzzzzz/fiber-boostrap/app/models"
	"github.com/Fairuzzzzz/fiber-boostrap/app/repository"
	"github.com/Fairuzzzzz/fiber-boostrap/pkg/jwt"
	"github.com/Fairuzzzzz/fiber-boostrap/pkg/response"
	"github.com/gofiber/fiber/v2"
	"go.elastic.co/apm"
	"golang.org/x/crypto/bcrypt"
)

func Register(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "Register", "controller")
	defer span.End()

	user := new(models.User)
	err := ctx.BodyParser(user)
	if err != nil {
		errResponse := fmt.Errorf("failed to parsing request: %v", err)
		log.Println(errResponse)
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
		log.Println(errResponse)
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
		log.Println(errResponse)
		return response.SendFailureResponse(
			ctx,
			fiber.StatusInternalServerError,
			errResponse.Error(),
			nil,
		)
	}

	user.Password = string(hashPassword)

	err = repository.InsertNewUser(spanCtx, user)
	if err != nil {
		errResponse := fmt.Errorf("failed to insert new user: %v", err)
		log.Println(errResponse)
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
	span, spanCtx := apm.StartSpan(ctx.Context(), "Login", "controller")
	defer span.End()

	now := time.Now()
	loginReq := new(models.LoginRequest)
	loginRes := models.LoginResponse{}

	err := ctx.BodyParser(loginReq)
	if err != nil {
		errResponse := fmt.Errorf("failed to parsing request: %v", err)
		log.Println(errResponse)
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
		log.Println(errResponse)
		return response.SendFailureResponse(
			ctx,
			fiber.StatusInternalServerError,
			errResponse.Error(),
			nil,
		)
	}

	user, err := repository.GetUserByUsername(spanCtx, loginReq.Username)
	if err != nil {
		errResponse := fmt.Errorf("failed to get username: %v", err)
		log.Println(errResponse)
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
		log.Println(errResponse)
		return response.SendFailureResponse(
			ctx,
			fiber.StatusNotFound,
			"username/password is invalid",
			nil,
		)
	}

	token, err := jwt.GenerateToken(spanCtx, user.Username, user.FullName, "token", now)
	if err != nil {
		errResponse := fmt.Errorf("failed to generate token: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(
			ctx,
			fiber.StatusInternalServerError,
			errResponse.Error(),
			nil,
		)
	}

	refreshToken, err := jwt.GenerateToken(
		spanCtx,
		user.Username,
		user.FullName,
		"refresh_token",
		now,
	)
	if err != nil {
		errResponse := fmt.Errorf("failed to generate token: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(
			ctx,
			fiber.StatusInternalServerError,
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

	err = repository.InsertNewUserSession(spanCtx, userSession)
	if err != nil {
		errResponse := fmt.Errorf("failed to insert user session: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(
			ctx,
			fiber.StatusInternalServerError,
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

func Logout(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "Logout", "controller")
	defer span.End()

	token := ctx.Get("Authorization")
	err := repository.DeleteUserSession(spanCtx, token)
	if err != nil {
		errResponse := fmt.Errorf("failed to delete user session: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(
			ctx,
			fiber.StatusInternalServerError,
			errResponse.Error(),
			nil,
		)
	}
	return response.SendSuccessResponse(ctx, nil)
}

func RefreshToken(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "RefreshToken", "controller")
	defer span.End()

	refreshToken := ctx.Get("Authorization")

	now := time.Now()
	username := ctx.Locals("username").(string)
	fullName := ctx.Locals("full_name").(string)

	token, err := jwt.GenerateToken(spanCtx, username, fullName, "token", now)
	if err != nil {
		errResponse := fmt.Errorf("failed to generate token: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(
			ctx,
			fiber.StatusInternalServerError,
			errResponse.Error(),
			nil,
		)
	}

	err = repository.UpdateUserSessionToken(
		spanCtx,
		token,
		now.Add(jwt.MapTokenType["token"]),
		refreshToken,
	)
	if err != nil {
		errResponse := fmt.Errorf("failed to generate token: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(
			ctx,
			fiber.StatusInternalServerError,
			errResponse.Error(),
			nil,
		)
	}

	return response.SendSuccessResponse(ctx, fiber.Map{
		"token": token,
	})
}
