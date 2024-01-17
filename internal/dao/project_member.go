package dao

import "github.com/letscrum/letscrum/internal/model"

type ProjectMemberDao interface {
	Get(projectMember model.ProjectMember) (*model.ProjectMember, error)
	GetByProject(project model.Project) (*model.ProjectMember, error)
	ListByProject(project model.Project, page, size int32) ([]*model.ProjectMember, error)
	Count() int64
	BatchAdd(projectMembers []model.ProjectMember) ([]*model.ProjectMember, error)
	BatchUpdate(projectMembers []model.ProjectMember) ([]*model.ProjectMember, error)
	BatchRemove(projectMembers []model.ProjectMember) ([]*model.ProjectMember, error)
}
