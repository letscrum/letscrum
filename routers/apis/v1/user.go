package v1

import (
	"github.com/gin-gonic/gin"
	userV1 "github.com/letscrum/letscrum/apis/user/v1"
	"github.com/letscrum/letscrum/pkg/errors"
	userService "github.com/letscrum/letscrum/services/user"
	"net/http"
)

func SignIn(ctx *gin.Context) {
	request := userV1.SignInRequest{}
	errRequest := ctx.ShouldBindJSON(&request)
	if errRequest != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(errRequest))
		return
	}

	user, err := userService.SignIn(request.Name, request.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(err))
		return
	}

	ctx.JSON(http.StatusOK, userV1.SignInResponse{
		Item: user,
	})
}
