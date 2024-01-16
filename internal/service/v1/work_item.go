package v1

import (
	v1 "github.com/letscrum/letscrum/api/letscrum/v1"
	"github.com/letscrum/letscrum/internal/dao"
)

type WorkItemService struct {
	v1.UnimplementedWorkItemServer
	workItemDao dao.WorkItemDao
}

func NewWorkItemService(dao dao.Interface) *WorkItemService {
	return &WorkItemService{
		workItemDao: dao.WorkItemDao(),
	}
}
