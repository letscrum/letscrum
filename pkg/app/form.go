package app

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"

	"github.com/letscrum/letscrum/pkg/errors"
)

// BindAndValid binds and validates data
func BindAndValid(c *gin.Context, form interface{}) (int, int) {
	err := c.BindWith(form, binding.Form)
	if err != nil {
		return http.StatusBadRequest, errors.INVALID_PARAMS
	}

	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		return http.StatusInternalServerError, errors.ERROR
	}
	if !check {
		MarkErrors(valid.Errors)
		return http.StatusBadRequest, errors.INVALID_PARAMS
	}

	return http.StatusOK, errors.SUCCESS
}
