package utils

import (
	"encoding/json"
	"regexp"

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

func IsSprintMember(sprintMembers []*projectv1.SprintMember, user model.User) bool {
	for _, m := range sprintMembers {
		if m.UserId == user.Id.String() {
			return true
		}
	}
	return false
}

func IsLegalName(name string) bool {
	if len(name) < 5 || len(name) > 50 {
		return false
	}
	result, err := regexp.Match(`^[a-z][a-z0-9_-]+$`, []byte(name))
	if err != nil {
		return false
	}
	return result
}
