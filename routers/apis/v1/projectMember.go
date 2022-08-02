package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	projectV1 "github.com/letscrum/letscrum/apis/project/v1"
	"github.com/letscrum/letscrum/pkg/errors"
	"github.com/letscrum/letscrum/pkg/utils"
	"github.com/letscrum/letscrum/services/projectMemberService"
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
	projectMemberId, err := projectMemberService.Create(request.ProjectId, request.UserId, request.IsAdmin)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.HandleErr(err))
		return
	}

	ctx.JSON(http.StatusOK, projectV1.CreateProjectMemberResponse{
		Id: projectMemberId,
	})
}

func ListProjectMembers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "")
}

func ListUserProject() {

}
