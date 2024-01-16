package dao

import (
	"github.com/letscrum/letscrum/internal/model"
)

type ProjectDao interface {
	Get(project model.Project) (*model.Project, error)
	List(page, size int32, keyword string) ([]*model.Project, error)
	Count(keyword string) int64
	Create(project model.Project) (*model.Project, error)
	Update(project model.Project) (*model.Project, error)
	Delete(project model.Project) (*model.Project, error)
}
