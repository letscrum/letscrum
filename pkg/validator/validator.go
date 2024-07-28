package validator

import (
	"encoding/json"

	projectv1 "github.com/letscrum/letscrum/api/project/v1"
	"github.com/letscrum/letscrum/internal/model"
)

func IsProjectAdmin(project model.Project, user model.User) bool {
	if user.IsSuperAdmin {
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
	if user.IsSuperAdmin {
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

func IsOrgAdmin(orgUser []*model.OrgUser, user model.User) bool {
	if user.IsSuperAdmin {
		return true
	}
	for _, m := range orgUser {
		if m.ForOrg.CreatedBy == user.Id {
			return true
		}
		if m.UserId == user.Id && m.IsAdmin == true {
			return true
		}
	}
	return false
}

func IsOrgMember(orgUser []model.OrgUser, user model.User) bool {
	if user.IsSuperAdmin {
		return true
	}
	for _, m := range orgUser {
		if m.ForOrg.CreatedBy == user.Id {
			return true
		}
		if m.UserId == user.Id {
			return true
		}
	}
	return false
}
