package app

import (
	"github.com/gin-gonic/gin"

	"github.com/letscrum/letscrum/pkg/errors"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Response settings gin.JSON
func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	g.C.JSON(httpCode, Response{
		Code: errCode,
		Msg:  errors.GetMsg(errCode),
		Data: data,
	})
	return
}
