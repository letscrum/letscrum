package utils

import (
	"encoding/json"

	generalv1 "github.com/letscrum/letscrum/api/general/v1"
	projectv1 "github.com/letscrum/letscrum/api/project/v1"
	"github.com/letscrum/letscrum/internal/model"
)

func GetOrgRole(org model.Org, orgUser []model.OrgUser, user model.User) generalv1.Role {
	if org.CreatedBy == user.Id {
		return generalv1.Role_OWNER
	}
	for _, m := range orgUser {
		if m.UserId == user.Id {
			if m.IsAdmin {
				return generalv1.Role_ADMIN
			}
			return generalv1.Role_MEMBER
		}
	}
	return generalv1.Role_NONE
}

func GetProjectRole(project model.Project, user model.User) generalv1.Role {
	if project.CreatedBy == user.Id {
		return generalv1.Role_OWNER
	}
	var projectMembers []*projectv1.ProjectMember
	err := json.Unmarshal([]byte(project.Members), &projectMembers)
	if err != nil {
		return generalv1.Role_NONE
	}
	for _, m := range projectMembers {
		if m.UserId == user.Id.String() {
			if m.IsAdmin {
				return generalv1.Role_ADMIN
			}
			return generalv1.Role_MEMBER
		}
	}
	return generalv1.Role_NONE
}
