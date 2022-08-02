package projectMemberService

import (
	"github.com/letscrum/letscrum/models"
)

func Create(projectId int64, userId int64, isAdmin bool) (int64, error) {
	projectMemberId, err := models.CreateProjectMember(projectId, userId, isAdmin)
	if err != nil {
		return 0, err
	}
	return projectMemberId, nil
}
