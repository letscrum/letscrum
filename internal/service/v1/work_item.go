package v1

import (
	"context"

	generalv1 "github.com/letscrum/letscrum/api/general/v1"
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

type WorkItemService struct {
	v1.UnimplementedWorkItemServer
	workItemDao dao.WorkItemDao
	projectDao  dao.ProjectDao
	taskDao     dao.TaskDao
}

func (s WorkItemService) Create(ctx context.Context, req *itemv1.CreateWorkItemRequest) (*itemv1.CreateWorkItemResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var user model.User
	user.Id = int64(claims.Id)
	user.IsSuperAdmin = claims.IsSuperAdmin
	var reqProject model.Project
	reqProject.Id = req.ProjectId
	project, err := s.projectDao.Get(reqProject)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if validator.IsProjectMember(*project, user) == false {
		return nil, status.Error(codes.PermissionDenied, "You are not a member of this project")
	}
	newWorkItem := model.WorkItem{
		ProjectId:   req.ProjectId,
		SprintId:    req.SprintId,
		FeatureId:   req.FeatureId,
		Title:       req.Title,
		Type:        req.Type.String(),
		Description: req.Description,
		Status:      itemv1.WorkItem_New.String(),
		CreatedBy:   int64(claims.Id),
	}
	workItem, err := s.workItemDao.Create(newWorkItem)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &itemv1.CreateWorkItemResponse{
		Success: workItem.Id > 0,
		Item: &itemv1.WorkItem{
			Id:          workItem.Id,
			ProjectId:   workItem.ProjectId,
			SprintId:    workItem.SprintId,
			FeatureId:   workItem.FeatureId,
			Title:       workItem.Title,
			Type:        itemv1.WorkItemType(itemv1.WorkItemType_value[workItem.Type]),
			Description: workItem.Description,
			Status:      itemv1.WorkItem_WorkItemStatus(itemv1.WorkItem_WorkItemStatus_value[workItem.Status]),
			AssignUser: &userv1.User{
				Id:    0,
				Name:  "",
				Email: "",
			},
			CreatedUser: &userv1.User{
				Id:    workItem.CreatedUser.Id,
				Name:  workItem.CreatedUser.Name,
				Email: workItem.CreatedUser.Email,
			},
			CreatedAt: workItem.CreatedAt.Unix(),
			UpdatedAt: workItem.UpdatedAt.Unix(),
		},
	}, nil
}

func (s WorkItemService) List(ctx context.Context, req *itemv1.ListWorkItemRequest) (*itemv1.ListWorkItemResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var user model.User
	user.Id = int64(claims.Id)
	user.IsSuperAdmin = claims.IsSuperAdmin
	var reqProject model.Project
	reqProject.Id = req.ProjectId
	project, err := s.projectDao.Get(reqProject)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if validator.IsProjectMember(*project, user) == false {
		return nil, status.Error(codes.PermissionDenied, "You are not a member of this project")
	}
	req.Page, req.Size = utils.Pagination(req.Page, req.Size)
	var workItems []*model.WorkItem
	count := int64(0)
	if req.ProjectId > 0 {
		if req.SprintId > 0 {
			workItems, err = s.workItemDao.ListBySprint(req.SprintId, req.Page, req.Size, req.Keyword)
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			count = s.workItemDao.CountBySprint(req.SprintId, req.Keyword)
		} else {
			workItems, err = s.workItemDao.ListByProject(req.ProjectId, req.Page, req.Size, req.Keyword)
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			count = s.workItemDao.CountByProject(req.ProjectId, req.Keyword)
		}
	}
	// get workitemIds list by workItems
	var workItemIds []int64
	for _, w := range workItems {
		workItemIds = append(workItemIds, w.Id)
	}
	// get tasks by workItemIds
	tasks, err := s.taskDao.ListByWorkItemIds(workItemIds)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	var items []*itemv1.WorkItem
	for _, w := range workItems {
		// convert w.AssignUser to userv1.User
		assignUser := &userv1.User{
			Id:    w.AssignUser.Id,
			Name:  w.AssignUser.Name,
			Email: w.AssignUser.Email,
		}
		// convert w.CreatedUser to userv1.User
		createdUser := &userv1.User{
			Id:    w.CreatedUser.Id,
			Name:  w.CreatedUser.Name,
			Email: w.CreatedUser.Email,
		}
		// get tasks by workItemId from tasks
		var tasksAll []*itemv1.Task
		var tasksUNKNOWN []*itemv1.Task
		var tasksToDo []*itemv1.Task
		var tasksInProgress []*itemv1.Task
		var tasksDone []*itemv1.Task
		var tasksRemoved []*itemv1.Task
		for _, t := range tasks {
			if t.WorkItemId == w.Id {
				resTask := &itemv1.Task{
					Id:          t.Id,
					WorkItemId:  t.WorkItemId,
					Title:       t.Title,
					Description: t.Description,
					Status:      itemv1.Task_TaskStatus(itemv1.Task_TaskStatus_value[t.Status]),
					AssignUser: &userv1.User{
						Id:    t.AssignUser.Id,
						Name:  t.AssignUser.Name,
						Email: t.AssignUser.Email,
					},
					CreatedUser: &userv1.User{
						Id:    t.CreatedUser.Id,
						Name:  t.CreatedUser.Name,
						Email: t.CreatedUser.Email,
					},
				}
				tasksAll = append(tasksAll, resTask)
				if resTask.Status == itemv1.Task_UNKNOWN {
					tasksUNKNOWN = append(tasksUNKNOWN, resTask)
				}
				if resTask.Status == itemv1.Task_ToDo {
					tasksToDo = append(tasksToDo, resTask)
				}
				if resTask.Status == itemv1.Task_InProgress {
					tasksInProgress = append(tasksInProgress, resTask)
				}
				if resTask.Status == itemv1.Task_Done {
					tasksDone = append(tasksDone, resTask)
				}
				if resTask.Status == itemv1.Task_Removed {
					tasksRemoved = append(tasksRemoved, resTask)
				}
			}
		}
		items = append(items, &itemv1.WorkItem{
			Id:              w.Id,
			ProjectId:       w.ProjectId,
			SprintId:        w.SprintId,
			FeatureId:       w.FeatureId,
			Title:           w.Title,
			Type:            itemv1.WorkItemType(itemv1.WorkItemType_value[w.Type]),
			Description:     w.Description,
			Status:          itemv1.WorkItem_WorkItemStatus(itemv1.WorkItem_WorkItemStatus_value[w.Status]),
			AssignUser:      assignUser,
			CreatedUser:     createdUser,
			TasksAll:        tasksAll,
			TasksUnknown:    tasksUNKNOWN,
			TasksToDo:       tasksToDo,
			TasksInProgress: tasksInProgress,
			TasksDone:       tasksDone,
			TasksRemoved:    tasksRemoved,
		})
	}
	// order items by its id desc
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[i].Id < items[j].Id {
				items[i], items[j] = items[j], items[i]
			}
		}
	}

	return &itemv1.ListWorkItemResponse{
		Items: items,
		Pagination: &generalv1.Pagination{
			Page:  req.Page,
			Size:  req.Size,
			Total: int32(count),
		},
	}, nil
}

func (s WorkItemService) Get(ctx context.Context, req *itemv1.GetWorkItemRequest) (*itemv1.GetWorkItemResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewWorkItemService(dao dao.Interface) *WorkItemService {
	return &WorkItemService{
		workItemDao: dao.WorkItemDao(),
		projectDao:  dao.ProjectDao(),
		taskDao:     dao.TaskDao(),
	}
}
