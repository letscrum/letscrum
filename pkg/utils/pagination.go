package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"

	"github.com/letscrum/letscrum/pkg/settings"
)

// GetPage get page parameters
func GetPage(c *gin.Context) int {
	result := 0
	page := com.StrTo(c.Query("page")).MustInt()
	if page > 0 {
		result = (page - 1) * settings.AppSetting.PageSize
	}

	return result
}
