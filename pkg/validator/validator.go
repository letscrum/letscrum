package validator

import (
	"encoding/json"

	projectv1 "github.com/letscrum/letscrum/api/project/v1"
	"github.com/letscrum/letscrum/internal/model"
)

func IsProjectAdmin(project model.Project, user model.User) bool {
	var projectMembers []*projectv1.ProjectMember
	err := json.Unmarshal([]byte(project.Members), &projectMembers)
	if err != nil {
		return false
	}
	for _, m := range projectMembers {
		if m.UserId == user.Id && m.IsAdmin == false {
			if user.IsSuperAdmin == false {
				return false
			}
		}
	}
	return true
}

func IsProjectMember(project model.Project, userId int64, isSuperAdmin bool) bool {
	if isSuperAdmin {
		return true
	}
	var projectMembers []*projectv1.ProjectMember
	err := json.Unmarshal([]byte(project.Members), &projectMembers)
	if err != nil {
		return false
	}
	// check claims.UserId in projectMembers
	for _, m := range projectMembers {
		if m.UserId == userId {
			return true
		}
	}
	return false
}
