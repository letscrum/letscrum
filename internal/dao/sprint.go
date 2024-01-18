package dao

import "github.com/letscrum/letscrum/internal/model"

type SprintDao interface {
	Get(sprint model.Sprint) (*model.Sprint, error)
	ListByProject(project model.Project, page, size int32, keyword string) ([]*model.Sprint, error)
	CountByProject(project model.Project, keyword string) int64
	Create(sprint model.Sprint) (*model.Sprint, error)
	Update(sprint model.Sprint) (*model.Sprint, error)
	Delete(sprint model.Sprint) (*model.Sprint, error)
}
