package v1

import (
	"context"

	"github.com/google/uuid"
	itemv1 "github.com/letscrum/letscrum/api/item/v1"
	v1 "github.com/letscrum/letscrum/api/letscrum/v1"
	userv1 "github.com/letscrum/letscrum/api/user/v1"
	"github.com/letscrum/letscrum/internal/dao"
	"github.com/letscrum/letscrum/internal/model"
	"github.com/letscrum/letscrum/pkg/utils"
	"github.com/letscrum/letscrum/pkg/validator"
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
	var user model.User
	user.Id = claims.Id
	user.IsSuperAdmin = claims.IsSuperAdmin
	var reqProject model.Project
	oId, err := uuid.Parse(req.OrgId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	pId, err := uuid.Parse(req.ProjectId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	reqProject.OrgId = oId
	reqProject.Id = pId
	project, err := t.projectDao.Get(reqProject)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if validator.IsProjectMember(*project, user) == false {
		return nil, status.Error(codes.PermissionDenied, "You are not a member of this project")
	}
	sId, err := uuid.Parse(req.SprintId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	newTask := model.Task{
		ProjectId:  project.Id,
		SprintId:   sId,
		WorkItemId: req.WorkItemId,
		Title:      req.Title,
		Status:     itemv1.Task_ToDo.String(),
		CreatedBy:  claims.Id,
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
			ProjectId:   task.ProjectId.String(),
			SprintId:    task.SprintId.String(),
			WorkItemId:  task.WorkItemId,
			Title:       task.Title,
			Description: task.Description,
			Status:      itemv1.Task_TaskStatus(itemv1.Task_TaskStatus_value[task.Status]),
			AssignUser: &userv1.User{
				Id:    uuid.Nil.String(),
				Name:  "",
				Email: "",
			},
			CreatedUser: &userv1.User{
				Id:    task.CreatedUser.Id.String(),
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

func (t TaskService) UpdateStatus(ctx context.Context, req *itemv1.UpdateTaskStatusRequest) (*itemv1.UpdateTaskResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var user model.User
	user.Id = claims.Id
	user.IsSuperAdmin = claims.IsSuperAdmin
	var reqProject model.Project
	oId, err := uuid.Parse(req.OrgId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	pId, err := uuid.Parse(req.ProjectId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	reqProject.OrgId = oId
	reqProject.Id = pId
	project, err := t.projectDao.Get(reqProject)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if validator.IsProjectMember(*project, user) == false {
		return nil, status.Error(codes.PermissionDenied, "You are not a member of this project")
	}
	var task model.Task
	task.Id = req.TaskId
	task.Status = req.Status.String()
	updateTask, err := t.taskDao.UpdateStatus(task)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &itemv1.UpdateTaskResponse{
		Success: true,
		Item: &itemv1.Task{
			Id:          updateTask.Id,
			ProjectId:   updateTask.ProjectId.String(),
			SprintId:    updateTask.SprintId.String(),
			WorkItemId:  updateTask.WorkItemId,
			Title:       updateTask.Title,
			Description: updateTask.Description,
			Status:      itemv1.Task_TaskStatus(itemv1.Task_TaskStatus_value[updateTask.Status]),
			AssignUser: &userv1.User{
				Id:    updateTask.AssignUser.Id.String(),
				Name:  updateTask.AssignUser.Name,
				Email: updateTask.AssignUser.Email,
			},
			CreatedUser: &userv1.User{
				Id:    updateTask.CreatedUser.Id.String(),
				Name:  updateTask.CreatedUser.Name,
				Email: updateTask.CreatedUser.Email,
			},
			CreatedAt: updateTask.CreatedAt.Unix(),
			UpdatedAt: updateTask.UpdatedAt.Unix(),
		},
	}, nil
}

func (t TaskService) Assign(ctx context.Context, req *itemv1.AssignTaskRequest) (*itemv1.UpdateTaskResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var user model.User
	user.Id = claims.Id
	user.IsSuperAdmin = claims.IsSuperAdmin
	var reqProject model.Project
	oId, err := uuid.Parse(req.OrgId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	pId, err := uuid.Parse(req.ProjectId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	reqProject.OrgId = oId
	reqProject.Id = pId
	project, err := t.projectDao.Get(reqProject)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if validator.IsProjectMember(*project, user) == false {
		return nil, status.Error(codes.PermissionDenied, "You are not a member of this project")
	}
	var task model.Task
	task.Id = req.TaskId
	auId, err := uuid.Parse(req.AssignUserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	task.AssignTo = auId
	updateTask, err := t.taskDao.UpdateAssignUser(task)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &itemv1.UpdateTaskResponse{
		Success: true,
		Item: &itemv1.Task{
			Id:          updateTask.Id,
			ProjectId:   updateTask.ProjectId.String(),
			SprintId:    updateTask.SprintId.String(),
			WorkItemId:  updateTask.WorkItemId,
			Title:       updateTask.Title,
			Description: updateTask.Description,
			Status:      itemv1.Task_TaskStatus(itemv1.Task_TaskStatus_value[updateTask.Status]),
			AssignUser: &userv1.User{
				Id:    updateTask.AssignUser.Id.String(),
				Name:  updateTask.AssignUser.Name,
				Email: updateTask.AssignUser.Email,
			},
			CreatedUser: &userv1.User{
				Id:    updateTask.CreatedUser.Id.String(),
				Name:  updateTask.CreatedUser.Name,
				Email: updateTask.CreatedUser.Email,
			},
			CreatedAt: updateTask.CreatedAt.Unix(),
			UpdatedAt: updateTask.UpdatedAt.Unix(),
		},
	}, nil
}

func (t TaskService) Move(ctx context.Context, req *itemv1.MoveTaskRequest) (*itemv1.UpdateTaskResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var user model.User
	user.Id = claims.Id
	user.IsSuperAdmin = claims.IsSuperAdmin
	var reqProject model.Project
	oId, err := uuid.Parse(req.OrgId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	pId, err := uuid.Parse(req.ProjectId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	reqProject.OrgId = oId
	reqProject.Id = pId
	project, err := t.projectDao.Get(reqProject)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if validator.IsProjectMember(*project, user) == false {
		return nil, status.Error(codes.PermissionDenied, "You are not a member of this project")
	}
	var task model.Task
	task.Id = req.TaskId
	task.Status = req.Status.String()
	task.WorkItemId = req.ToWorkItemId
	updateTask, err := t.taskDao.Move(task)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &itemv1.UpdateTaskResponse{
		Success: true,
		Item: &itemv1.Task{
			Id:          updateTask.Id,
			ProjectId:   updateTask.ProjectId.String(),
			SprintId:    updateTask.SprintId.String(),
			WorkItemId:  updateTask.WorkItemId,
			Title:       updateTask.Title,
			Description: updateTask.Description,
			Status:      itemv1.Task_TaskStatus(itemv1.Task_TaskStatus_value[updateTask.Status]),
			AssignUser: &userv1.User{
				Id:    updateTask.AssignUser.Id.String(),
				Name:  updateTask.AssignUser.Name,
				Email: updateTask.AssignUser.Email,
			},
			CreatedUser: &userv1.User{
				Id:    updateTask.CreatedUser.Id.String(),
				Name:  updateTask.CreatedUser.Name,
				Email: updateTask.CreatedUser.Email,
			},
			CreatedAt: updateTask.CreatedAt.Unix(),
			UpdatedAt: updateTask.UpdatedAt.Unix(),
		},
	}, nil
}
