package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	generalV1 "github.com/letscrum/letscrum/apis/general/v1"
	projectV1 "github.com/letscrum/letscrum/apis/project/v1"
	"github.com/letscrum/letscrum/internal/service/projectmemberservice"
	"github.com/letscrum/letscrum/pkg/errors"
	"github.com/letscrum/letscrum/pkg/utils"
	"net/http"
)

func CreateProjectMember(ctx *gin.Context) {
	request := projectV1.CreateProjectMemberRequest{}

	errRequest := ctx.ShouldBindJSON(&request)
	if errRequest != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(errRequest))
		return
	}
	projectId := utils.GetParamInt64(ctx, "project_id")
	if projectId <= 0 {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("wrong id param"))
		return
	}
	request.ProjectId = projectId
	projectMemberId, err := projectmemberservice.Create(request.ProjectId, request.UserId, request.IsAdmin)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(err))
		return
	}

	ctx.JSON(http.StatusOK, projectV1.CreateProjectMemberResponse{
		Id: projectMemberId,
	})
}

func ListProjectMember(ctx *gin.Context) {
	request := projectV1.ListProjectMemberRequest{}
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
	projectId := utils.GetParamInt64(ctx, "project_id")
	if projectId <= 0 {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("wrong id param"))
		return
	}
	request.ProjectId = projectId
	projectMembers, count, err := projectmemberservice.ListProjectMemberByProject(request.ProjectId, request.Page, request.PageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(err))
		return
	}

	ctx.JSON(http.StatusOK, projectV1.ListProjectMemberResponse{
		Items: projectMembers,
		Pagination: &generalV1.Pagination{
			Page:     request.Page,
			PageSize: request.PageSize,
			Total:    int32(count),
		},
	})
}

func ListUserProject(ctx *gin.Context) {
	request := projectV1.ListUserProjectRequest{}
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
	projects, count, err := projectmemberservice.ListProjectByUser(request.UserId, request.Page, request.PageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(err))
		return
	}

	ctx.JSON(http.StatusOK, projectV1.ListUserProjectResponse{
		Items: projects,
		Pagination: &generalV1.Pagination{
			Page:     request.Page,
			PageSize: request.PageSize,
			Total:    int32(count),
		},
	})
}
