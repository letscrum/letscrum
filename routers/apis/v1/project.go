package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	generalV1 "github.com/letscrum/letscrum/apis/general/v1"
	projectV1 "github.com/letscrum/letscrum/apis/project/v1"
	"github.com/letscrum/letscrum/pkg/errors"
	"github.com/letscrum/letscrum/pkg/utils"
	"github.com/letscrum/letscrum/services/projectService"
	"net/http"
)

func CreateProject(ctx *gin.Context) {
	request := projectV1.CreateProjectRequest{}

	errRequest := ctx.ShouldBindJSON(&request)
	if errRequest != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(errRequest))
		return
	}
	projectId, err := projectService.Create(request.Project.Name, request.Project.DisplayName, ctx.GetInt64("userId"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(err))
		return
	}

	ctx.JSON(http.StatusOK, projectV1.CreateProjectResponse{
		Id: projectId,
	})
}

func ListProject(ctx *gin.Context) {
	request := projectV1.ListProjectRequest{}

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

	list, count, err := projectService.List(request.Page, request.PageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(err))
		return
	}

	ctx.JSON(http.StatusOK, projectV1.ListProjectResponse{
		Items: list,
		Pagination: &generalV1.Pagination{
			Page:     request.Page,
			PageSize: request.PageSize,
			Total:    int32(count),
		},
	})
}

func UpdateProject(ctx *gin.Context) {
	request := projectV1.UpdateProjectRequest{}

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

	err := projectService.Update(request.ProjectId, request.Project.DisplayName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(err))
		return
	}

	ctx.JSON(http.StatusOK, projectV1.UpdateProjectResponse{
		Id: request.ProjectId,
	})
}

func DeleteProject(ctx *gin.Context) {
	request := projectV1.DeleteProjectRequest{}
	projectId := utils.GetParamInt64(ctx, "project_id")
	if projectId <= 0 {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("wrong id param"))
		return
	}
	request.ProjectId = projectId

	err := projectService.Delete(request.ProjectId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(err))
		return
	}

	ctx.JSON(http.StatusOK, projectV1.DeleteProjectResponse{
		Id: request.ProjectId,
	})
}

func HardDeleteProject(ctx *gin.Context) {
	request := projectV1.DeleteProjectRequest{}
	projectId := utils.GetParamInt64(ctx, "project_id")
	if projectId <= 0 {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("wrong id param"))
		return
	}
	request.ProjectId = projectId

	err := projectService.HardDelete(request.ProjectId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(err))
		return
	}

	ctx.JSON(http.StatusOK, projectV1.DeleteProjectResponse{
		Id: request.ProjectId,
	})
}

func GetProject(ctx *gin.Context) {
	request := projectV1.GetProjectRequest{}

	projectId := utils.GetParamInt64(ctx, "project_id")
	if projectId <= 0 {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("wrong id param"))
		return
	}
	request.ProjectId = projectId

	project, err := projectService.Get(request.ProjectId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(err))
		return
	}

	ctx.JSON(http.StatusOK, projectV1.GetProjectResponse{
		Item: project,
	})
}
