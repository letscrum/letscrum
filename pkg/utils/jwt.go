package utils

import (
	"context"
	"errors"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type key string

func GenerateTokens(userId uuid.UUID, isSuperAdmin bool) (string, string, error) {
	nowTime := time.Now()
	accessTokenExpireTime := nowTime.Add(time.Hour * 720)
	accessTokenClaims := jwt.MapClaims{
		"iss": userId,
		"aud": isSuperAdmin,
		"exp": accessTokenExpireTime.Unix(),
	}
	refreshTokenExpireTime := nowTime.Add(time.Hour * 720 * 2)
	refreshTokenClaims := jwt.MapClaims{
		"iss": userId,
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
	Id             uuid.UUID `json:"id"`
	IsSuperAdmin   bool      `json:"is_super_admin"`
	ExpirationTime time.Time `json:"exp"`
}

func GetTokenDetails(ctx context.Context) (UserClaims, error) {
	claims := ctx.Value("token").(jwt.MapClaims)
	if claims == nil {
		return UserClaims{}, errors.New("token claims not found")
	}
	user := UserClaims{
		Id:             uuid.MustParse(claims["iss"].(string)),
		IsSuperAdmin:   claims["aud"].(bool),
		ExpirationTime: time.Unix(int64(claims["exp"].(float64)), 0),
	}
	return user, nil
}
