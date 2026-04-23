package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	accessSecret  string
	refreshSecret string
}

func NewJWTService(accessSecret, refreshSecret string) *JWTService {
	return &JWTService{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
	}
}

// generate access token
func (j *JWTService) GenerateAccessToken(userID uint, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.accessSecret))
}

// generate refresh token
func (j *JWTService) GenerateRefreshToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.refreshSecret))
}
