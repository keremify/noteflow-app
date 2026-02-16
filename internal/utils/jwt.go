package utils

import (
	"time"

	"crypto/rand"
	"encoding/hex"
	"github.com/golang-jwt/jwt/v5"
	"saasproject/internal/config"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, role string) (string, error) {
	cfg := config.Get()

	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(time.Duration(cfg.JWT.ExpiresIn) * time.Hour),
			),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}

func GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32) // 256 bit
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
