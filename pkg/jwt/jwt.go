package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/Fairuzzzzz/fiber-boostrap/pkg/env"
	"github.com/golang-jwt/jwt/v5"
	"go.elastic.co/apm"
)

type ClaimToken struct {
	Username string `json:"username"`
	Fullname string `json:"full_name"`
	jwt.RegisteredClaims
}

var MapTokenType = map[string]time.Duration{
	"token":         time.Hour * 3,
	"refresh_token": time.Hour * 72,
}

var jwtSecret = []byte(env.GetEnv("APP_SECRET", ""))

func GenerateToken(
	ctx context.Context,
	username, fullname string,
	tokenType string,
	now time.Time,
) (string, error) {
	span, _ := apm.StartSpan(ctx, "GenerateToken", "jwt")
	defer span.End()

	claimToken := ClaimToken{
		Username: username,
		Fullname: fullname,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    env.GetEnv("APP_NAME", ""),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(MapTokenType[tokenType])),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimToken)

	resultToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return resultToken, fmt.Errorf("failed to generate token: %v", err)
	}
	return resultToken, nil
}

func ValidateToken(ctx context.Context, token string) (*ClaimToken, error) {
	span, _ := apm.StartSpan(ctx, "ValidateToken", "jwt")
	defer span.End()

	var (
		claimToken *ClaimToken
		ok         bool
	)

	jwtToken, err := jwt.ParseWithClaims(
		token,
		&ClaimToken{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("failed to validate method jwt: %v", t.Header["alg"])
			}
			return jwtSecret, nil
		})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	if claimToken, ok = jwtToken.Claims.(*ClaimToken); !ok || !jwtToken.Valid {
		return nil, fmt.Errorf("token invalid")
	}

	return claimToken, nil
}
