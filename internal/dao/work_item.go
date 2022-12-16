package dao

import "github.com/letscrum/letscrum/internal/model"

type WorkItemDao interface {
	Get(id int64) (*model.WorkItem, error)
	List(page, size int32, keyword string) ([]*model.WorkItem, error)
	Count(keyword string) int64
	Create(project *model.WorkItem) (int64, error)
	Update(project *model.WorkItem) (bool, error)
	Delete(id int64) (bool, error)
}
