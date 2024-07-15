package dao

import "github.com/letscrum/letscrum/internal/model"

type TaskDao interface {
	Get(task model.Task) (*model.Task, error)
	List(page, size int32, keyword string) ([]*model.Task, error)
	ListByWorkItem(workItemId int64, page, size int32, keyword string) ([]*model.Task, error)
	ListByWorkItemIds(workItemIds []int64) ([]*model.Task, error)
	Count(keyword string) int64
	CountByWorkItem(workItemId int64, keyword string) int64
	Create(task model.Task) (*model.Task, error)
	Update(task model.Task) (*model.Task, error)
	UpdateStatus(task model.Task) (*model.Task, error)
	UpdateAssignUser(task model.Task) (*model.Task, error)
	Move(task model.Task) (*model.Task, error)
	Delete(task model.Task) (bool, error)
}
