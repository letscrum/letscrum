package v1

import (
	"context"
	"encoding/json"

	itemv1 "github.com/letscrum/letscrum/api/item/v1"
	v1 "github.com/letscrum/letscrum/api/letscrum/v1"
	projectv1 "github.com/letscrum/letscrum/api/project/v1"
	userv1 "github.com/letscrum/letscrum/api/user/v1"
	"github.com/letscrum/letscrum/internal/dao"
	"github.com/letscrum/letscrum/internal/model"
	"github.com/letscrum/letscrum/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TaskService struct {
	v1.UnimplementedTaskServer
	taskDao    dao.TaskDao
	projectDao dao.ProjectDao
}

func NewTaskService(dao dao.Interface) *TaskService {
	return &TaskService{
		taskDao:    dao.TaskDao(),
		projectDao: dao.ProjectDao(),
	}
}

func (t TaskService) Create(ctx context.Context, req *itemv1.CreateTaskRequest) (*itemv1.CreateTaskResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var reqProject model.Project
	reqProject.Id = req.ProjectId
	project, err := t.projectDao.Get(reqProject)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	var projectMembers []*projectv1.ProjectMember
	err = json.Unmarshal([]byte(project.Members), &projectMembers)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	// check claims.UserId in projectMembers
	var isMember bool
	for _, m := range projectMembers {
		if m.UserId == int64(claims.Id) {
			isMember = true
			break
		}
	}
	if !isMember {
		return nil, status.Error(codes.PermissionDenied, "You are not a member of this project")
	}
	newTask := model.Task{
		ProjectId:  req.ProjectId,
		SprintId:   req.SprintId,
		WorkItemId: req.WorkItemId,
		Title:      req.Title,
		Status:     itemv1.Task_ToDo.String(),
		CreatedBy:  int64(claims.Id),
	}
	task, err := t.taskDao.Create(newTask)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert task to itemv1.Task
	var resTask itemv1.Task

	success := false
	if task.Id > 0 {
		success = true
		resTask = itemv1.Task{
			Id:          task.Id,
			ProjectId:   task.ProjectId,
			SprintId:    task.SprintId,
			WorkItemId:  task.WorkItemId,
			Title:       task.Title,
			Description: task.Description,
			Status:      itemv1.Task_UNKNOWN,
			AssignUser:  nil,
			CreatedUser: &userv1.User{
				Id:    task.CreatedUser.Id,
				Name:  task.CreatedUser.Name,
				Email: task.CreatedUser.Email,
			},
			CreatedAt: task.CreatedAt.Unix(),
			UpdatedAt: task.UpdatedAt.Unix(),
		}
	}
	return &itemv1.CreateTaskResponse{
		Success: success,
		Item:    &resTask,
	}, nil
}

func (t TaskService) List(ctx context.Context, req *itemv1.ListTaskRequest) (*itemv1.ListTaskResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (t TaskService) Get(ctx context.Context, req *itemv1.GetTaskRequest) (*itemv1.GetTaskResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (t TaskService) Update(ctx context.Context, req *itemv1.UpdateTaskRequest) (*itemv1.UpdateTaskResponse, error) {
	//TODO implement me
	panic("implement me")
}
