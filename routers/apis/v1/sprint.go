package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	projectV1 "github.com/letscrum/letscrum/apis/project/v1"
	"github.com/letscrum/letscrum/internal/service/sprintservice"
	"github.com/letscrum/letscrum/pkg/errors"
	"github.com/letscrum/letscrum/pkg/utils"
	"net/http"
	"time"
)

func CreateSprint(ctx *gin.Context) {
	request := projectV1.CreateSprintRequest{}

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
	sprintId, err := sprintservice.Create(request.ProjectId, request.Sprint.Name, time.Unix(request.Sprint.StartDate, 0), time.Unix(request.Sprint.EndDate, 0))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(err))
		return
	}

	ctx.JSON(http.StatusOK, projectV1.CreateSprintResponse{
		Id: sprintId,
	})
}
