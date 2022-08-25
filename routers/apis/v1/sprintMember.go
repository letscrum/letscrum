package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	generalV1 "github.com/letscrum/letscrum/apis/general/v1"
	projectV1 "github.com/letscrum/letscrum/apis/project/v1"
	"github.com/letscrum/letscrum/internal/service/sprintmemberservice"
	"github.com/letscrum/letscrum/pkg/errors"
	"github.com/letscrum/letscrum/pkg/utils"
	"net/http"
)

func CreateSprintMember(ctx *gin.Context) {
	request := projectV1.CreateSprintMemberRequest{}

	errRequest := ctx.ShouldBindJSON(&request)
	if errRequest != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(errRequest))
		return
	}
	sprintId := utils.GetParamInt64(ctx, "sprint_id")
	if sprintId <= 0 {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("wrong id param"))
		return
	}
	request.SprintId = sprintId
	sprintMemberId, err := sprintmemberservice.Create(request.SprintId, request.UserId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(err))
		return
	}

	ctx.JSON(http.StatusOK, projectV1.CreateSprintMemberResponse{
		Id: sprintMemberId,
	})
}

func ListSprintMember(ctx *gin.Context) {
	request := projectV1.ListSprintMemberRequest{}
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
	sprintId := utils.GetParamInt64(ctx, "sprint_id")
	if sprintId <= 0 {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("wrong id param"))
		return
	}
	request.SprintId = sprintId
	sprintMembers, count, err := sprintmemberservice.ListSprintMemberBySprint(request.SprintId, request.Page, request.PageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(err))
		return
	}

	ctx.JSON(http.StatusOK, projectV1.ListProjectMemberResponse{
		Items: sprintMembers,
		Pagination: &generalV1.Pagination{
			Page:     request.Page,
			PageSize: request.PageSize,
			Total:    int32(count),
		},
	})
}

func ListUserSprint(ctx *gin.Context) {
	request := projectV1.ListUserSprintRequest{}
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
	userId := utils.GetParamInt64(ctx, "user_id")
	if userId <= 0 {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("wrong id param"))
		return
	}
	request.UserId = userId
	sprints, count, err := sprintmemberservice.ListSprintByUser(request.UserId, request.Page, request.PageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(err))
		return
	}

	ctx.JSON(http.StatusOK, projectV1.ListUserSprintResponse{
		Items: sprints,
		Pagination: &generalV1.Pagination{
			Page:     request.Page,
			PageSize: request.PageSize,
			Total:    int32(count),
		},
	})
}
