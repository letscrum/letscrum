package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	generalV1 "github.com/letscrum/letscrum/apis/general/v1"
	projectV1 "github.com/letscrum/letscrum/apis/project/v1"
	projectService "github.com/letscrum/letscrum/services/project"
	"net/http"
)

func CreateProject(ctx *gin.Context) {
	request := projectV1.CreateProjectRequest{}

	errRequest := ctx.ShouldBindJSON(&request)
	if errRequest != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
	}

	err := projectService.Create(&projectV1.Project{
		Name:        request.Name,
		DisplayName: request.DisplayName,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
	}

	ctx.JSON(http.StatusOK, projectV1.CreateProjectResponse{
		Item: &projectV1.Project{
			Name:        request.Name,
			DisplayName: request.DisplayName,
		},
	})
}

func ListProject(ctx *gin.Context) {
	request := projectV1.ListProjectRequest{}

	errRequest := ctx.ShouldBindWith(&request, binding.Form)
	if errRequest != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
	}

	if request.Page <= 0 {
		request.Page = 1
	}
	if request.PageSize <= 0 {
		request.PageSize = 10
	}

	list, err := projectService.List(&generalV1.Pagination{
		Page:     request.Page,
		PageSize: request.PageSize,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
	}

	ctx.JSON(http.StatusOK, projectV1.ListProjectResponse{
		Items: list,
		Pagination: &generalV1.Pagination{
			Page:     request.Page,
			PageSize: request.PageSize,
		},
	})
}
