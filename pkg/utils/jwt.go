package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type LetscrumClaims struct {
	jwt.StandardClaims

	IsSuperAdmin bool `json:"is_super_admin"`
}

func GenerateTokens(userId string, isSuperAdmin bool) (string, string, error) {
	nowTime := time.Now()
	accessTokenExpireTime := nowTime.Add(time.Hour * 720)
	accessTokenClaims := LetscrumClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessTokenExpireTime.Unix(),
			Id:        userId,
		},
		IsSuperAdmin: isSuperAdmin,
	}
	refreshTokenExpireTime := nowTime.Add(time.Hour * 720 * 2)
	refreshTokenClaims := LetscrumClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshTokenExpireTime.Unix(),
			Id:        userId,
		},
		IsSuperAdmin: isSuperAdmin,
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
func ParseToken(token string) (*LetscrumClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &LetscrumClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*LetscrumClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
