package dao

import "github.com/letscrum/letscrum/internal/model"

type WorkItemDao interface {
	Get(workItem model.WorkItem) (*model.WorkItem, error)
	ListByProject(projectId int64, page, size int32, keyword string) ([]*model.WorkItem, error)
	ListBySprint(sprintId int64, page, size int32, keyword string) ([]*model.WorkItem, error)
	CountByProject(projectId int64, keyword string) int64
	CountBySprint(sprintId int64, keyword string) int64
	Create(workItem model.WorkItem) (*model.WorkItem, error)
	Update(workItem model.WorkItem) (*model.WorkItem, error)
	Delete(workItem model.WorkItem) (bool, error)
}
