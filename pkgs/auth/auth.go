package auth

import (
	"errors"

	"github.com/GazDuckington/go-gin/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateJWT(tokenString string, cfg *config.Config) (any, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token.Claims, nil
}
