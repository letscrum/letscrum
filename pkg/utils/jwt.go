package utils

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateTokens(userId string) (string, string, error) {
	nowTime := time.Now()
	accessTokenExpireTime := nowTime.Add(time.Hour * 720)
	accessTokenClaims := jwt.StandardClaims{
		ExpiresAt: accessTokenExpireTime.Unix(),
		Id:        userId,
	}
	refreshTokenExpireTime := nowTime.Add(time.Hour * 720 * 2)
	refreshTokenClaims := jwt.StandardClaims{
		ExpiresAt: refreshTokenExpireTime.Unix(),
		Id:        userId,
	}
	accessToken, errAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims).SignedString([]byte(os.Getenv("JWT_SECRET")))
	if errAccessToken != nil {
		return "", "", errAccessToken
	}
	refreshToken, errRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims).SignedString([]byte(os.Getenv("JWT_SECRET")))
	if errRefreshToken != nil {
		return "", "", errRefreshToken
	}
	return accessToken, refreshToken, nil
}

// ParseToken parsing token
func ParseToken(token string) (*jwt.StandardClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*jwt.StandardClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
