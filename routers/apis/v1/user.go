package v1

import (
	"github.com/gin-gonic/gin"
	userV1 "github.com/letscrum/letscrum/apis/user/v1"
	"github.com/letscrum/letscrum/pkg/errors"
	"github.com/letscrum/letscrum/services/userService"
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

func CreateUser(ctx *gin.Context) {
	request := userV1.CreateUserRequest{}
	errRequest := ctx.ShouldBindJSON(&request)
	if errRequest != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(errRequest))
		return
	}
	user := userV1.User{
		Name:         request.Name,
		Email:        request.Email,
		Password:     request.Password,
		IsSuperAdmin: request.IsSuperAdmin,
	}
	err := userService.Create(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(err))
		return
	}

	ctx.JSON(http.StatusOK, userV1.CreateUserResponse{
		Item: &user,
	})
}
