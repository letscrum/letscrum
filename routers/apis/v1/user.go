package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	generalV1 "github.com/letscrum/letscrum/apis/general/v1"
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
	id, err := userService.Create(request.Name, request.Email, request.Password, request.IsSuperAdmin)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(err))
		return
	}

	ctx.JSON(http.StatusOK, userV1.CreateUserResponse{
		Id: id,
	})
}

func ListUser(ctx *gin.Context) {
	request := userV1.ListUserRequest{}

	errRequest := ctx.ShouldBindWith(&request, binding.Form)
	if errRequest != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(errRequest))
		return
	}

	if request.Page <= 0 {
		request.Page = 1
	}
	if request.PageSize <= 0 {
		request.PageSize = 10
	}

	list, count, err := userService.List(request.Keyword, request.Page, request.PageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(err))
		return
	}

	ctx.JSON(http.StatusOK, userV1.ListUserResponse{
		Items: list,
		Pagination: &generalV1.Pagination{
			Page:     request.Page,
			PageSize: request.PageSize,
			Total:    int32(count),
		},
	})
}
