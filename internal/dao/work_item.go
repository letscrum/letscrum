package dao

import (
	"github.com/google/uuid"
	"github.com/letscrum/letscrum/internal/model"
)

type WorkItemDao interface {
	Get(workItem model.WorkItem) (*model.WorkItem, []*model.ItemLog, error)
	ListByProject(projectId uuid.UUID, page, size int32, keyword string) ([]*model.WorkItem, error)
	ListBySprint(sprintId uuid.UUID, page, size int32, keyword string) ([]*model.WorkItem, error)
	CountByProject(projectId uuid.UUID, keyword string) int64
	CountBySprint(sprintId uuid.UUID, keyword string) int64
	CountBySprints(sprintIds []uuid.UUID) []int64
	Create(workItem model.WorkItem) (*model.WorkItem, error)
	Update(workItem model.WorkItem, userId uuid.UUID) (*model.WorkItem, error)
	UpdateAssignUser(workItem model.WorkItem, userId uuid.UUID) (*model.WorkItem, error)
	UpdateStatus(workItem model.WorkItem, userId uuid.UUID) (*model.WorkItem, error)
	UpdateSprintWithTasks(workItem model.WorkItem, userId uuid.UUID) (*model.WorkItem, error)
	Delete(workItem model.WorkItem, userId uuid.UUID) (bool, error)
	ReOrder(workItemIds []int64) ([]int64, error)
}
