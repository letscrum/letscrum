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
)

func CreateProject(ctx *gin.Context) {
	request := projectV1.CreateProjectRequest{}

	errRequest := ctx.ShouldBindJSON(&request)
	if errRequest != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(errRequest))
		return
	}
	err := projectService.Create(&projectV1.Project{
		Name:        request.Name,
		DisplayName: request.DisplayName,
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
			Name:        request.Name,
			DisplayName: request.DisplayName,
			CreatedUser: &userV1.User{
				Id: ctx.GetInt64("userId"),
			},
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
	request.ProjectName = ctx.Param("project_name")

	err := projectService.Update(&projectV1.Project{
		Name:        request.ProjectName,
		DisplayName: request.Project.DisplayName,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(err))
		return
	}

	ctx.JSON(http.StatusOK, projectV1.UpdateProjectResponse{
		Project: &projectV1.Project{
			Name:        request.ProjectName,
			DisplayName: request.Project.DisplayName,
		},
	})
}

func DeleteProject(ctx *gin.Context) {
	request := projectV1.DeleteProjectRequest{}
	request.ProjectName = ctx.Param("name")

	err := projectService.Delete(request.ProjectName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(err))
		return
	}

	ctx.JSON(http.StatusOK, projectV1.DeleteProjectResponse{
		Project: &projectV1.Project{
			Name: request.ProjectName,
		},
	})
}

func GetProject(ctx *gin.Context) {
	request := projectV1.GetProjectRequest{}
	request.ProjectName = ctx.Param("project_name")

	project, err := projectService.Get(request.ProjectName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(err))
		return
	}

	ctx.JSON(http.StatusOK, projectV1.GetProjectResponse{
		Project: project,
	})
}

func ListProjectMembers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "")
}
