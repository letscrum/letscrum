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
	if utils.IsProjectMember(*project, user) == false {
		return nil, status.Error(codes.PermissionDenied, utils.ErrNotProjectMember)
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
	if utils.IsProjectMember(*project, user) == false {
		return nil, status.Error(codes.PermissionDenied, utils.ErrNotProjectMember)
	}
	var task model.Task
	task.Id = req.TaskId
	getTask, getLogs, err := t.taskDao.Get(task)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	// convert task to itemv1.Task
	var resTask itemv1.Task
	if getTask.Id > 0 {
		// convert getLogs to itemv1.ItemLog
		var resLogs []*itemv1.Log
		for _, log := range getLogs {
			resLogs = append(resLogs, &itemv1.Log{
				Id:        log.Id.String(),
				ItemId:    log.ItemId,
				ItemType:  itemv1.Log_ItemType(itemv1.Log_ItemType_value[log.ItemType]),
				Action:    itemv1.Log_Action(itemv1.Log_Action_value[log.Action]),
				Log:       log.Log,
				Changes:   log.Changes,
				CreatedBy: log.CreatedBy.String(),
				CreatedAt: log.CreatedAt.Unix(),
			})
		}
		resTask = itemv1.Task{
			Id:          getTask.Id,
			ProjectId:   getTask.ProjectId.String(),
			SprintId:    getTask.SprintId.String(),
			WorkItemId:  getTask.WorkItemId,
			Title:       getTask.Title,
			Description: getTask.Description,
			Status:      itemv1.Task_TaskStatus(itemv1.Task_TaskStatus_value[getTask.Status]),
			Remaining:   getTask.Remaining,
			AssignUser: &userv1.User{
				Id:    getTask.AssignUser.Id.String(),
				Name:  getTask.AssignUser.Name,
				Email: getTask.AssignUser.Email,
			},
			CreatedUser: &userv1.User{
				Id:    getTask.CreatedUser.Id.String(),
				Name:  getTask.CreatedUser.Name,
				Email: getTask.CreatedUser.Email,
			},
			Logs:      resLogs,
			CreatedAt: getTask.CreatedAt.Unix(),
			UpdatedAt: getTask.UpdatedAt.Unix(),
		}
	}
	return &itemv1.GetTaskResponse{
		Item: &resTask,
	}, nil
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
	if utils.IsProjectMember(*project, user) == false {
		return nil, status.Error(codes.PermissionDenied, utils.ErrNotProjectMember)
	}
	var task model.Task
	task.Id = req.TaskId
	task.Status = req.Status.String()
	updateTask, err := t.taskDao.Move(task, user.Id)
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
			Remaining:   updateTask.Remaining,
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
	if utils.IsProjectMember(*project, user) == false {
		return nil, status.Error(codes.PermissionDenied, utils.ErrNotProjectMember)
	}
	var task model.Task
	task.Id = req.TaskId
	auId, err := uuid.Parse(req.AssignUserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	task.AssignTo = auId
	updateTask, err := t.taskDao.UpdateAssignUser(task, user.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
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
			Remaining:   updateTask.Remaining,
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
	if utils.IsProjectMember(*project, user) == false {
		return nil, status.Error(codes.PermissionDenied, utils.ErrNotProjectMember)
	}
	var task model.Task
	task.Id = req.TaskId
	task.Status = req.Status.String()
	task.WorkItemId = req.ToWorkItemId
	updateTask, err := t.taskDao.Move(task, user.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
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
			Remaining:   updateTask.Remaining,
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

func (t TaskService) ReOrder(ctx context.Context, req *itemv1.ReOrderTasksRequest) (*itemv1.ReOrderTasksResponse, error) {
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
	if utils.IsProjectMember(*project, user) == false {
		return nil, status.Error(codes.PermissionDenied, utils.ErrNotProjectMember)
	}
	_, err = t.taskDao.ReOrder(req.TaskIds)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &itemv1.ReOrderTasksResponse{
		Success: true,
	}, nil
}

func (t TaskService) UpdateWorkHours(ctx context.Context, req *itemv1.UpdateWorkHoursRequest) (*itemv1.UpdateTaskResponse, error) {
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
	if utils.IsProjectMember(*project, user) == false {
		return nil, status.Error(codes.PermissionDenied, utils.ErrNotProjectMember)
	}
	var task model.Task
	task.ProjectId = project.Id
	task.Id = req.TaskId
	task.Remaining = float32(req.Remaining)

	updateTask, err := t.taskDao.UpdateWorkHours(task, user.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
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
			Remaining:   updateTask.Remaining,
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

func (t TaskService) Delete(ctx context.Context, req *itemv1.DeleteTaskRequest) (*itemv1.DeleteTaskResponse, error) {
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
	if utils.IsProjectMember(*project, user) == false {
		return nil, status.Error(codes.PermissionDenied, utils.ErrNotProjectMember)
	}
	var task model.Task
	task.Id = req.TaskId
	success := false
	success, err = t.taskDao.Delete(task, user.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &itemv1.DeleteTaskResponse{
		Success: success,
	}, nil
}
