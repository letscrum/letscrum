package service

import (
	v1 "github.com/letscrum/letscrum/api/letscrum/v1"
	"github.com/letscrum/letscrum/internal/dao"
)

type TaskService struct {
	v1.UnimplementedTaskServer
	taskDao dao.TaskDao
}

func NewTaskService(dao dao.Interface) *TaskService {
	return &TaskService{
		taskDao: dao.TaskDao(),
	}
}
