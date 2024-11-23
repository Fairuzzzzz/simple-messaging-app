package repository

import (
	"context"
	"time"

	"github.com/Fairuzzzzz/fiber-boostrap/app/models"
	"github.com/Fairuzzzzz/fiber-boostrap/pkg/database"
)

func InsertNewUser(ctx context.Context, user *models.User) error {
	return database.DB.Create(user).Error
}

func InsertNewUserSession(ctx context.Context, user *models.UserSession) error {
	return database.DB.Create(user).Error
}

func DeleteUserSession(ctx context.Context, token string) error {
	return database.DB.Exec("DELETE FROM user_sessions WHERE token = ?", token).Error
}

func UpdateUserSessionToken(
	ctx context.Context,
	token string,
	tokenExpired time.Time,
	refreshToken string,
) error {
	return database.DB.Exec(
		"UPDATE user_sessions SET token = ?, token_expired = ? WHERE refresh_token = ?",
		token,
		tokenExpired,
		refreshToken,
	).Error
}

func GetUserSession(ctx context.Context, token string) (models.UserSession, error) {
	var (
		resp models.UserSession
		err  error
	)

	err = database.DB.Where("token = ?", token).Last(&resp).Error
	return resp, err
}

func GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	var (
		resp models.User
		err  error
	)

	err = database.DB.Where("username = ?", username).Last(&resp).Error
	return resp, err
}
