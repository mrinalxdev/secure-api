package utils

import (
	"secure-api/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWTSecret []byte

func InitJWT (secret string) {
	JWTSecret = []byte(secret)
}

type Claims struct {
	UserID uint `json:"user_id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID uint, role string) (string, string, error) {
	config := config.LoadConfig()

	accessTokenClaims := Claims {
		UserID: userID,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.JWTExpiry)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Issuer: "secure-api",
		},
	}


	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)


	accessTokenString, err := accessToken.SignedString(JWTSecret)

	if err != nil {
		return "", "", err
	}


	refreshTokenClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.RefreshExpiry)),

		IssuedAt : jwt.NewNumericDate(time.Now()),
		Issuer : "secure-api",
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString(JWTSecret)

	if err != nil {
		return "", "", err
	}


	return accessTokenString, refreshTokenString, nil
}

func ValidateToken(tokenString string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}