package v1

import (
	"context"
	"encoding/json"

	generalv1 "github.com/letscrum/letscrum/api/general/v1"
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

type WorkItemService struct {
	v1.UnimplementedWorkItemServer
	workItemDao dao.WorkItemDao
	projectDao  dao.ProjectDao
}

func (s WorkItemService) Create(ctx context.Context, req *itemv1.CreateWorkItemRequest) (*itemv1.CreateWorkItemResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var reqProject model.Project
	reqProject.Id = req.ProjectId
	project, err := s.projectDao.Get(reqProject)
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
	newWorkItem := model.WorkItem{
		ProjectId:   req.ProjectId,
		SprintId:    req.SprintId,
		FeatureId:   req.FeatureId,
		Title:       req.Title,
		Type:        req.Type.String(),
		Description: req.Description,
		Status:      itemv1.WorkItemStatus_New.String(),
		CreatedBy:   int64(claims.Id),
	}
	workItem, err := s.workItemDao.Create(newWorkItem)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	success := false
	if workItem.Id > 0 {
		success = true
	}
	return &itemv1.CreateWorkItemResponse{
		Success: success,
		Id:      workItem.Id,
	}, nil
}

func (s WorkItemService) ListByProject(ctx context.Context, req *itemv1.ListWorkItemRequest) (*itemv1.ListWorkItemResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var reqProject model.Project
	reqProject.Id = req.ProjectId
	project, err := s.projectDao.Get(reqProject)
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
	req.Page, req.Size = utils.Pagination(req.Page, req.Size)
	workItems, err := s.workItemDao.ListByProject(req.ProjectId, req.Page, req.Size, req.Keyword)
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
		items = append(items, &itemv1.WorkItem{
			Id:          w.Id,
			ProjectId:   w.ProjectId,
			SprintId:    w.SprintId,
			FeatureId:   w.FeatureId,
			Title:       w.Title,
			Type:        itemv1.WorkItemType(itemv1.WorkItemType_value[w.Type]),
			Description: w.Description,
			Status:      itemv1.WorkItemStatus(itemv1.WorkItemStatus_value[w.Status]),
			AssignUser:  assignUser,
			CreatedUser: createdUser,
		})
	}
	count := s.workItemDao.CountByProject(req.ProjectId, req.Keyword)
	return &itemv1.ListWorkItemResponse{
		Items: items,
		Pagination: &generalv1.Pagination{
			Page:  req.Page,
			Size:  req.Size,
			Total: int32(count),
		},
	}, nil
}

func (s WorkItemService) ListBySprint(ctx context.Context, req *itemv1.ListWorkItemRequest) (*itemv1.ListWorkItemResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var reqProject model.Project
	reqProject.Id = req.ProjectId
	project, err := s.projectDao.Get(reqProject)
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
	req.Page, req.Size = utils.Pagination(req.Page, req.Size)
	workItems, err := s.workItemDao.ListBySprint(req.SprintId, req.Page, req.Size, req.Keyword)
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
		items = append(items, &itemv1.WorkItem{
			Id:          w.Id,
			ProjectId:   w.ProjectId,
			SprintId:    w.SprintId,
			FeatureId:   w.FeatureId,
			Title:       w.Title,
			Type:        itemv1.WorkItemType(itemv1.WorkItemType_value[w.Type]),
			Description: w.Description,
			Status:      itemv1.WorkItemStatus(itemv1.WorkItemStatus_value[w.Status]),
			AssignUser:  assignUser,
			CreatedUser: createdUser,
		})
	}
	count := s.workItemDao.CountBySprint(req.SprintId, req.Keyword)
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
	}
}
