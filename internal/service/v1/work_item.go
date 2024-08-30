package v1

import (
	"context"

	"github.com/google/uuid"
	generalv1 "github.com/letscrum/letscrum/api/general/v1"
	itemv1 "github.com/letscrum/letscrum/api/item/v1"
	v1 "github.com/letscrum/letscrum/api/letscrum/v1"
	userv1 "github.com/letscrum/letscrum/api/user/v1"
	"github.com/letscrum/letscrum/internal/dao"
	"github.com/letscrum/letscrum/internal/model"
	"github.com/letscrum/letscrum/pkg/utils"
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
	project, err := s.projectDao.Get(reqProject)
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
	newWorkItem := model.WorkItem{
		ProjectId:   reqProject.Id,
		SprintId:    sId,
		FeatureId:   req.FeatureId,
		Title:       req.Title,
		Type:        req.Type.String(),
		Description: req.Description,
		Status:      itemv1.WorkItem_New.String(),
		CreatedBy:   claims.Id,
	}
	workItem, err := s.workItemDao.Create(newWorkItem)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &itemv1.CreateWorkItemResponse{
		Success: workItem.Id > 0,
		Item: &itemv1.WorkItem{
			Id:          workItem.Id,
			ProjectId:   workItem.ProjectId.String(),
			SprintId:    workItem.SprintId.String(),
			FeatureId:   workItem.FeatureId,
			Title:       workItem.Title,
			Type:        itemv1.WorkItemType(itemv1.WorkItemType_value[workItem.Type]),
			Description: workItem.Description,
			Status:      itemv1.WorkItem_WorkItemStatus(itemv1.WorkItem_WorkItemStatus_value[workItem.Status]),
			AssignUser: &userv1.User{
				Id:    uuid.Nil.String(),
				Name:  "",
				Email: "",
			},
			CreatedUser: &userv1.User{
				Id:    workItem.CreatedUser.Id.String(),
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
	project, err := s.projectDao.Get(reqProject)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if utils.IsProjectMember(*project, user) == false {
		return nil, status.Error(codes.PermissionDenied, utils.ErrNotProjectMember)
	}
	req.Page, req.Size = utils.Pagination(req.Page, req.Size)
	var workItems []*model.WorkItem
	count := int64(0)
	// if req.ProjectId is not empty uuid string "00000000-0000-0000-0000-000000000000"
	if req.ProjectId != uuid.Nil.String() {
		if req.SprintId != uuid.Nil.String() {
			sId, err := uuid.Parse(req.SprintId)
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			workItems, err = s.workItemDao.ListBySprint(sId, req.Page, req.Size, req.Keyword)
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			count = s.workItemDao.CountBySprint(sId, req.Keyword)
		} else {
			workItems, err = s.workItemDao.ListByProject(project.Id, req.Page, req.Size, req.Keyword)
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			count = s.workItemDao.CountByProject(project.Id, req.Keyword)
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
			Id:    w.AssignUser.Id.String(),
			Name:  w.AssignUser.Name,
			Email: w.AssignUser.Email,
		}
		// convert w.CreatedUser to userv1.User
		createdUser := &userv1.User{
			Id:    w.CreatedUser.Id.String(),
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
						Id:    t.AssignUser.Id.String(),
						Name:  t.AssignUser.Name,
						Email: t.AssignUser.Email,
					},
					CreatedUser: &userv1.User{
						Id:    t.CreatedUser.Id.String(),
						Name:  t.CreatedUser.Name,
						Email: t.CreatedUser.Email,
					},
					Order: t.Order,
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
			ProjectId:       w.ProjectId.String(),
			SprintId:        w.SprintId.String(),
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
			Order:           w.Order,
		})
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
	project, err := s.projectDao.Get(reqProject)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if utils.IsProjectMember(*project, user) == false {
		return nil, status.Error(codes.PermissionDenied, utils.ErrNotProjectMember)
	}
	var workItem model.WorkItem
	workItem.Id = req.WorkItemId
	getWorkItem, getLogs, err := s.workItemDao.Get(workItem)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert w.AssignUser to userv1.User
	assignUser := &userv1.User{
		Id:    getWorkItem.AssignUser.Id.String(),
		Name:  getWorkItem.AssignUser.Name,
		Email: getWorkItem.AssignUser.Email,
	}
	// convert w.CreatedUser to userv1.User
	createdUser := &userv1.User{
		Id:    getWorkItem.CreatedUser.Id.String(),
		Name:  getWorkItem.CreatedUser.Name,
		Email: getWorkItem.CreatedUser.Email,
	}
	// get tasks by workItemId from tasks
	var tasksAll []*itemv1.Task
	var resLogs []*itemv1.Log
	for _, l := range getLogs {
		resLog := &itemv1.Log{
			Id:        l.Id.String(),
			ItemId:    l.ItemId,
			ItemType:  itemv1.Log_ItemType(itemv1.Log_ItemType_value[l.ItemType]),
			Action:    itemv1.Log_Action(itemv1.Log_Action_value[l.Action]),
			Log:       l.Log,
			Changes:   l.Changes,
			CreatedBy: l.CreatedBy.String(),
			CreatedAt: l.CreatedAt.Unix(),
		}
		resLogs = append(resLogs, resLog)
	}
	resWorkItem := &itemv1.WorkItem{
		Id:          getWorkItem.Id,
		ProjectId:   getWorkItem.ProjectId.String(),
		SprintId:    getWorkItem.SprintId.String(),
		FeatureId:   getWorkItem.FeatureId,
		Title:       getWorkItem.Title,
		Type:        itemv1.WorkItemType(itemv1.WorkItemType_value[getWorkItem.Type]),
		Description: getWorkItem.Description,
		Status:      itemv1.WorkItem_WorkItemStatus(itemv1.WorkItem_WorkItemStatus_value[getWorkItem.Status]),
		AssignUser:  assignUser,
		CreatedUser: createdUser,
		TasksAll:    tasksAll,
		Logs:        resLogs,
	}
	return &itemv1.GetWorkItemResponse{
		Item: resWorkItem,
	}, nil
}

func (s WorkItemService) Assign(ctx context.Context, req *itemv1.AssignWorkItemRequest) (*itemv1.UpdateWorkItemResponse, error) {
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
	project, err := s.projectDao.Get(reqProject)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if utils.IsProjectMember(*project, user) == false {
		return nil, status.Error(codes.PermissionDenied, utils.ErrNotProjectMember)
	}
	uId, err := uuid.Parse(req.AssignUserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var workItem model.WorkItem
	workItem.Id = req.WorkItemId
	workItem.AssignTo = uId
	_, err = s.workItemDao.UpdateAssignUser(workItem, user.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &itemv1.UpdateWorkItemResponse{
		Success: true,
		Item:    &itemv1.WorkItem{},
	}, nil
}

func (s WorkItemService) UpdateStatus(ctx context.Context, req *itemv1.UpdateWorkItemStatusRequest) (*itemv1.UpdateWorkItemResponse, error) {
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
	project, err := s.projectDao.Get(reqProject)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if utils.IsProjectMember(*project, user) == false {
		return nil, status.Error(codes.PermissionDenied, utils.ErrNotProjectMember)
	}
	var workItem model.WorkItem
	workItem.Id = req.WorkItemId
	workItem.Status = req.Status.String()
	_, err = s.workItemDao.UpdateStatus(workItem, user.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &itemv1.UpdateWorkItemResponse{
		Success: true,
		Item:    &itemv1.WorkItem{},
	}, nil
}

func (s WorkItemService) Move(ctx context.Context, req *itemv1.MoveWorkItemRequest) (*itemv1.UpdateWorkItemResponse, error) {
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
	project, err := s.projectDao.Get(reqProject)
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
	var workItem model.WorkItem
	workItem.Id = req.WorkItemId
	workItem.SprintId = sId
	_, err = s.workItemDao.UpdateSprintWithTasks(workItem, user.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &itemv1.UpdateWorkItemResponse{
		Success: true,
		Item:    &itemv1.WorkItem{},
	}, nil
}

func (s WorkItemService) ReOrder(ctx context.Context, req *itemv1.ReOrderWorkItemsRequest) (*itemv1.ReOrderWorkItemsResponse, error) {
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
	project, err := s.projectDao.Get(reqProject)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if utils.IsProjectMember(*project, user) == false {
		return nil, status.Error(codes.PermissionDenied, utils.ErrNotProjectMember)
	}
	_, err = s.workItemDao.ReOrder(req.WorkItemIds)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &itemv1.ReOrderWorkItemsResponse{
		Success: true,
	}, nil
}

func NewWorkItemService(dao dao.Interface) *WorkItemService {
	return &WorkItemService{
		workItemDao: dao.WorkItemDao(),
		projectDao:  dao.ProjectDao(),
		taskDao:     dao.TaskDao(),
	}
}
