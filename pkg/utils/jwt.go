package utils

import (
	"context"
	"errors"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func GenerateTokens(userId float64, userName string, isSuperAdmin bool) (string, string, error) {
	nowTime := time.Now()
	accessTokenExpireTime := nowTime.Add(time.Hour * 720)
	accessTokenClaims := jwt.MapClaims{
		"iss": userId,
		"sub": userName,
		"aud": isSuperAdmin,
		"exp": accessTokenExpireTime.Unix(),
	}
	refreshTokenExpireTime := nowTime.Add(time.Hour * 720 * 2)
	refreshTokenClaims := jwt.MapClaims{
		"iss": userId,
		"sub": userName,
		"aud": isSuperAdmin,
		"exp": refreshTokenExpireTime.Unix(),
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
func ParseToken(token string) (jwt.MapClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	if tokenClaims != nil {
		claims, ok := tokenClaims.Claims.(jwt.MapClaims)
		if ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

type UserClaims struct {
	ID             float64   `json:"id"`
	Name           string    `json:"name"`
	IsSuperAdmin   bool      `json:"is_super_admin"`
	ExpirationTime time.Time `json:"exp"`
}

func GetTokenDetails(ctx context.Context) (UserClaims, error) {
	claims := ctx.Value("token").(jwt.MapClaims)
	if claims == nil {
		return UserClaims{}, errors.New("token claims not found")
	}
	user := UserClaims{
		ID:             claims["iss"].(float64),
		Name:           claims["sub"].(string),
		IsSuperAdmin:   claims["aud"].(bool),
		ExpirationTime: time.Unix(int64(claims["exp"].(float64)), 0),
	}
	return user, nil
}
