package utils

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetParamInt64(ctx *gin.Context, name string) int64 {
	int64Param, errGetParam := strconv.ParseInt(ctx.Param(name), 10, 64)
	if errGetParam != nil {
		return 0
	}
	return int64Param
}
