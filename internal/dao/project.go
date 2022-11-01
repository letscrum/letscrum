package dao

import (
	"github.com/letscrum/letscrum/internal/model"
)

type ProjectDao interface {
	Get(id int64) (*model.Project, error)
	List(page, size int32) ([]*model.Project, error)
	Count() int64
	Create(project *model.Project) (bool, error)
	Update(project *model.Project) (bool, error)
	Delete(id int64) (bool, error)
}
