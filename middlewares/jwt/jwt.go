package jwtMiddleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/letscrum/letscrum/pkg/utils"
	"net/http"
	"strconv"

	"github.com/letscrum/letscrum/pkg/errors"
)

// JWT is jwt middleware
func JWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, errors.HandleErr(fmt.Errorf("miss token")))
			ctx.Abort()
			return
		}
		claims, err := utils.ParseToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, errors.HandleErr(err))
			ctx.Abort()
			return
		}
		userId, errSetUserId := strconv.ParseInt(claims.Id, 10, 64)
		if errSetUserId != nil {
			ctx.JSON(http.StatusUnauthorized, errors.HandleErr(errSetUserId))
			ctx.Abort()
			return
		}
		ctx.Set("userId", userId)
		ctx.Next()
	}
}
