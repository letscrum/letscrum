package dao

import (
	"github.com/letscrum/letscrum/internal/model"
)

type ProjectDao interface {
	Get(project model.Project) (*model.Project, error)
	List(org model.Org, page, size int32, keyword string) ([]*model.Project, error)
	Count(org model.Org, keyword string) int64
	ListVisibleProject(org model.Org, page, size int32, keyword string, user model.User) ([]*model.Project, error)
	CountVisibleProject(org model.Org, keyword string, user model.User) int64
	CountByOrg(org model.Org) int64
	Create(project model.Project) (*model.Project, error)
	Update(project model.Project) (*model.Project, error)
	Delete(project model.Project) (bool, error)
}
