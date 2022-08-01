package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	generalV1 "github.com/letscrum/letscrum/apis/general/v1"
	projectV1 "github.com/letscrum/letscrum/apis/project/v1"
	userV1 "github.com/letscrum/letscrum/apis/user/v1"
	"github.com/letscrum/letscrum/pkg/errors"
	projectService "github.com/letscrum/letscrum/services/project"
	"net/http"
	"strconv"
)

func CreateProject(ctx *gin.Context) {
	request := projectV1.CreateProjectRequest{}

	errRequest := ctx.ShouldBindJSON(&request)
	if errRequest != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(errRequest))
		return
	}
	projectId, err := projectService.Create(&projectV1.Project{
		Name:        request.Project.Name,
		DisplayName: request.Project.DisplayName,
		CreatedUser: &userV1.User{
			Id: ctx.GetInt64("userId"),
		},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(err))
		return
	}

	ctx.JSON(http.StatusOK, projectV1.CreateProjectResponse{
		Item: &projectV1.Project{
			Id: projectId,
		},
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

	list, count, err := projectService.List(&generalV1.Pagination{
		Page:     request.Page,
		PageSize: request.PageSize,
	})
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

	projectId, errGetParam := strconv.ParseInt(ctx.Param("project_id"), 10, 64)
	if errGetParam != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(errGetParam))
		return
	}
	request.ProjectId = projectId

	err := projectService.Update(&projectV1.Project{
		Id:          request.ProjectId,
		DisplayName: request.Project.DisplayName,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(err))
		return
	}

	ctx.JSON(http.StatusOK, projectV1.UpdateProjectResponse{
		Item: &projectV1.Project{
			Id: request.ProjectId,
		},
	})
}

func DeleteProject(ctx *gin.Context) {
	request := projectV1.DeleteProjectRequest{}
	projectId, errGetParam := strconv.ParseInt(ctx.Param("project_id"), 10, 64)
	if errGetParam != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(errGetParam))
		return
	}
	request.ProjectId = projectId

	err := projectService.Delete(request.ProjectId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(err))
		return
	}

	ctx.JSON(http.StatusOK, projectV1.DeleteProjectResponse{
		Item: &projectV1.Project{
			Id: request.ProjectId,
		},
	})
}

func GetProject(ctx *gin.Context) {
	request := projectV1.GetProjectRequest{}

	projectId, errGetParam := strconv.ParseInt(ctx.Param("project_id"), 10, 64)
	if errGetParam != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(errGetParam))
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

func ListProjectMembers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "")
}
