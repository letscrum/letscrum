package dao

import (
	"github.com/google/uuid"
	"github.com/letscrum/letscrum/internal/model"
)

type TaskDao interface {
	Get(task model.Task) (*model.Task, []*model.ItemLog, error)
	List(page, size int32, keyword string) ([]*model.Task, error)
	ListByWorkItem(workItemId int64, page, size int32, keyword string) ([]*model.Task, error)
	ListByWorkItemIds(workItemIds []int64) ([]*model.Task, error)
	Count(keyword string) int64
	CountByWorkItem(workItemId int64, keyword string) int64
	CountBySprint(sprintId uuid.UUID, keyword string) int64
	CountBySprints(sprintIds []uuid.UUID) []int64
	Create(task model.Task) (*model.Task, error)
	Update(task model.Task, userId uuid.UUID) (*model.Task, error)
	UpdateStatus(task model.Task, userId uuid.UUID) (*model.Task, error)
	UpdateAssignUser(task model.Task, userId uuid.UUID) (*model.Task, error)
	Move(task model.Task, userId uuid.UUID) (*model.Task, error)
	Delete(task model.Task, userId uuid.UUID) (bool, error)
	ReOrder(taskIds []int64) ([]int64, error)
}
