package utils

import (
	"encoding/json"

	projectv1 "github.com/letscrum/letscrum/api/project/v1"
	"github.com/letscrum/letscrum/internal/model"
)

func IsProjectAdmin(project model.Project, user model.User) bool {
	if project.CreatedBy == user.Id {
		return true
	}
	var projectMembers []*projectv1.ProjectMember
	err := json.Unmarshal([]byte(project.Members), &projectMembers)
	if err != nil {
		return false
	}
	for _, m := range projectMembers {
		if m.UserId == user.Id.String() && m.IsAdmin == true {
			return true
		}
	}
	return false
}

func IsProjectMember(project model.Project, user model.User) bool {
	if project.CreatedBy == user.Id {
		return true
	}
	var projectMembers []*projectv1.ProjectMember
	err := json.Unmarshal([]byte(project.Members), &projectMembers)
	if err != nil {
		return false
	}
	// check claims.UserId in projectMembers
	for _, m := range projectMembers {
		if m.UserId == user.Id.String() {
			return true
		}
	}
	return false
}

func IsOrgAdmin(org model.Org, orgUser []model.OrgUser, user model.User) bool {
	if org.CreatedBy == user.Id {
		return true
	}
	for _, m := range orgUser {
		if m.UserId == user.Id && m.IsAdmin == true {
			return true
		}
	}
	return false
}

func IsOrgMember(org model.Org, orgUser []model.OrgUser, user model.User) bool {
	if org.CreatedBy == user.Id {
		return true
	}
	for _, m := range orgUser {
		if m.UserId == user.Id {
			return true
		}
	}
	return false
}
