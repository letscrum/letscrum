package dao

import "github.com/letscrum/letscrum/internal/model"

type TaskDao interface {
	Get(id int64) (*model.Task, error)
	List(page, size int32, keyword string) ([]*model.Task, error)
	Count(keyword string) int64
	Create(project *model.Task) (int64, error)
	Update(project *model.Task) (bool, error)
	Delete(id int64) (bool, error)
}
