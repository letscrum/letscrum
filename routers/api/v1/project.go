package v1

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/letscrum/letscrum/apis/project/v1"
	projectService "github.com/letscrum/letscrum/services/project"
	"net/http"
)

func CreateProject(ctx *gin.Context) {
	request := v1.CreateProjectRequest{}
	ctx.ShouldBindJSON(&request)
	projectService.CreateProject(&v1.Project{
		Name:        request.Name,
		DisplayName: request.DisplayName,
	})

	ctx.JSON(http.StatusOK, v1.CreateProjectResponse{
		Project: &v1.Project{
			Name:        request.Name,
			DisplayName: request.DisplayName,
		},
	})
}
