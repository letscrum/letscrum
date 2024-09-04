package v1

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
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

type ProjectService struct {
	v1.UnimplementedProjectServer
	projectDao  dao.ProjectDao
	userDao     dao.UserDao
	sprintDao   dao.SprintDao
	orgDao      dao.OrgDao
	workItemDao dao.WorkItemDao
	taskDao     dao.TaskDao
}

func NewProjectService(dao dao.Interface) *ProjectService {
	return &ProjectService{
		projectDao:  dao.ProjectDao(),
		userDao:     dao.UserDao(),
		sprintDao:   dao.SprintDao(),
		orgDao:      dao.OrgDao(),
		workItemDao: dao.WorkItemDao(),
		taskDao:     dao.TaskDao(),
	}
}

func (s ProjectService) Get(ctx context.Context, req *projectv1.GetProjectRequest) (*projectv1.GetProjectResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var user model.User
	user.Id = claims.Id
	user.IsSuperAdmin = claims.IsSuperAdmin
	var reqOrg model.Org
	oId, err := uuid.Parse(req.OrgId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	reqOrg.Id = oId
	org, err := s.orgDao.Get(reqOrg)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	orgUsers, err := s.orgDao.ListMember(reqOrg)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if utils.IsOrgMember(org, orgUsers, user) == false {
		return nil, status.Error(codes.PermissionDenied, utils.ErrNotOrgMember)
	}

	var reqProject model.Project
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
	if project.Id == uuid.Nil {
		return nil, status.Error(codes.NotFound, "project not fount.")
	}
	if utils.IsProjectMember(*project, user) == false {
		return nil, status.Error(codes.PermissionDenied, utils.ErrNotProjectMember)
	}
	sprints, err := s.sprintDao.ListByProject(reqProject, 1, 999, "")
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	sprint := projectv1.Sprint{}
	for i, s := range sprints {
		// Set the 1st Sprint as default
		if i == 0 {
			sprint = projectv1.Sprint{
				Id:        s.Id.String(),
				ProjectId: s.ProjectId.String(),
				Name:      s.Name,
				StartDate: s.StartDate.Unix(),
				EndDate:   s.EndDate.Unix(),
				Status:    projectv1.Sprint_UNKNOWN,
				CreatedAt: s.CreatedAt.Unix(),
				UpdatedAt: s.UpdatedAt.Unix(),
			}
		}
		// Set the real current sprint
		if time.Now().After(s.StartDate) && time.Now().Before(s.EndDate) {
			sprint = projectv1.Sprint{
				Id:        s.Id.String(),
				ProjectId: s.ProjectId.String(),
				Name:      s.Name,
				StartDate: s.StartDate.Unix(),
				EndDate:   s.EndDate.Unix(),
				Status:    projectv1.Sprint_Current,
				CreatedAt: s.CreatedAt.Unix(),
				UpdatedAt: s.UpdatedAt.Unix(),
			}
			break
		}
	}
	var members []*projectv1.ProjectMember
	err = json.Unmarshal([]byte(project.Members), &members)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	members = append(members, &projectv1.ProjectMember{
		UserId:   project.CreatedUser.Id.String(),
		UserName: project.CreatedUser.Name,
		IsAdmin:  true,
	})
	return &projectv1.GetProjectResponse{
		Item: &projectv1.Project{
			Id:          project.Id.String(),
			Name:        project.Name,
			DisplayName: project.DisplayName,
			Description: project.Description,
			CreatedUser: &userv1.User{
				Id:           project.CreatedUser.Id.String(),
				Name:         project.CreatedUser.Name,
				IsSuperAdmin: project.CreatedUser.IsSuperAdmin,
			},
			Members:       members,
			CurrentSprint: &sprint,
			CreatedAt:     project.CreatedAt.Unix(),
			UpdatedAt:     project.UpdatedAt.Unix(),
			MyRole:        utils.GetProjectRole(*project, user),
		},
	}, nil
}

func (s *ProjectService) List(ctx context.Context, req *projectv1.ListProjectRequest) (*projectv1.ListProjectResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var user model.User
	user.Id = claims.Id
	user.IsSuperAdmin = claims.IsSuperAdmin
	req.Page, req.Size = utils.Pagination(req.Page, req.Size)
	var reqOrg model.Org
	oId, err := uuid.Parse(req.OrgId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	reqOrg.Id = oId
	projects, err := s.projectDao.ListVisibleProject(reqOrg, req.Page, req.Size, req.Keyword, user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var list []*projectv1.Project
	for _, p := range projects {
		var members []*projectv1.ProjectMember
		if p.Members != "" {
			err = json.Unmarshal([]byte(p.Members), &members)
			if err != nil {
				return nil, status.Error(codes.Unknown, err.Error())
			}
		}
		members = append(members, &projectv1.ProjectMember{
			UserId:   p.CreatedUser.Id.String(),
			UserName: p.CreatedUser.Name,
			IsAdmin:  true,
		})
		var project = &projectv1.Project{
			Id:          p.Id.String(),
			Name:        p.Name,
			DisplayName: p.DisplayName,
			Description: p.Description,
			Members:     members,
			CreatedUser: &userv1.User{
				Id:           p.CreatedUser.Id.String(),
				Name:         p.CreatedUser.Name,
				IsSuperAdmin: p.CreatedUser.IsSuperAdmin,
			},
			CreatedAt: p.CreatedAt.Unix(),
			UpdatedAt: p.UpdatedAt.Unix(),
		}
		list = append(list, project)
	}
	count := s.projectDao.CountVisibleProject(reqOrg, req.Keyword, user)
	return &projectv1.ListProjectResponse{
		Items: list,
		Pagination: &generalv1.Pagination{
			Page:  req.Page,
			Size:  req.Size,
			Total: int32(count),
		},
	}, nil
}

func (s *ProjectService) Create(ctx context.Context, req *projectv1.CreateProjectRequest) (*projectv1.CreateProjectResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var reqUser model.User
	reqUser.Id = claims.Id
	user, err := s.userDao.Get(reqUser)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var reqOrg model.Org
	oId, err := uuid.Parse(req.OrgId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	reqOrg.Id = oId
	org, err := s.orgDao.Get(reqOrg)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	projectCount := s.projectDao.CountByOrg(org)
	if projectCount >= org.ProjectLimitation {
		return nil, status.Error(codes.PermissionDenied, utils.ErrReachProjectLimit)
	}

	orgUsers, err := s.orgDao.ListMember(reqOrg)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if utils.IsOrgAdmin(org, orgUsers, *user) == false {
		return nil, status.Error(codes.PermissionDenied, utils.ErrNotOrgAdmin)
	}
	if !utils.IsLegalName(req.Name) {
		return nil, status.Error(codes.InvalidArgument, "project name can't be less than 5, can only contain lower case letter, number, _ and -, only can start from lower case letter.")
	}
	var members []*projectv1.ProjectMember
	if req.Members != nil && len(req.Members) > 0 {
		// convert req.Members to id list
		var userIds []uuid.UUID
		for _, m := range req.Members {
			if m.UserId != user.Id.String() {
				uid, err := uuid.Parse(m.UserId)
				if err != nil {
					return nil, status.Error(codes.Internal, err.Error())
				}
				userIds = append(userIds, uid)
			}
		}
		users, err := s.userDao.ListByIds(1, 999, userIds)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		// add members to project members
		for _, u := range users {
			for _, rm := range req.Members {
				if u.Id.String() == rm.UserId {
					member := &projectv1.ProjectMember{
						UserId:   u.Id.String(),
						UserName: u.Name,
						IsAdmin:  rm.IsAdmin,
					}
					members = append(members, member)
					continue
				}
			}
		}
	}
	// convert members to json string
	membersJson, err := json.Marshal(members)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	var newProject model.Project
	newProject.OrgId = oId
	newProject.Id = uuid.New()
	newProject.Name = req.Name
	newProject.DisplayName = req.DisplayName
	newProject.Description = req.Description
	newProject.Members = string(membersJson)
	newProject.CreatedBy = user.Id

	project, err := s.projectDao.Create(newProject)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &projectv1.CreateProjectResponse{
		Success: project.Id != uuid.Nil,
		Id:      project.Id.String(),
	}, nil
}

func (s *ProjectService) Update(ctx context.Context, req *projectv1.UpdateProjectRequest) (*projectv1.UpdateProjectResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var user model.User
	user.Id = claims.Id
	user.IsSuperAdmin = claims.IsSuperAdmin
	var reqOrg model.Org
	oId, err := uuid.Parse(req.OrgId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	reqOrg.Id = oId
	org, err := s.orgDao.Get(reqOrg)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	orgUsers, err := s.orgDao.ListMember(reqOrg)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if utils.IsOrgMember(org, orgUsers, user) == false {
		return nil, status.Error(codes.PermissionDenied, utils.ErrNotOrgMember)
	}
	var reqProject model.Project
	pId, err := uuid.Parse(req.ProjectId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	reqProject.OrgId = oId
	reqProject.Id = pId
	reqProject.CreatedUser.Id = user.Id
	project, err := s.projectDao.Get(reqProject)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if project.Id == uuid.Nil {
		return nil, status.Error(codes.NotFound, "project not fount.")
	}
	if utils.IsProjectAdmin(*project, user) == false {
		if utils.IsOrgAdmin(org, orgUsers, user) == false {
			return nil, status.Error(codes.PermissionDenied, utils.ErrNoAdminPermissionForProject)
		}
	}

	project.DisplayName = req.DisplayName
	project.Description = req.Description
	var members []*projectv1.ProjectMember
	if req.Members != nil && len(req.Members) > 0 {
		for _, m := range req.Members {
			if project.CreatedUser.Id.String() != m.UserId {
				member := &projectv1.ProjectMember{
					UserId:   m.UserId,
					UserName: m.UserName,
					IsAdmin:  m.IsAdmin,
				}
				members = append(members, member)
			}
		}
	}
	membersJson, err := json.Marshal(members)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	project.Members = string(membersJson)

	updatedProject, err := s.projectDao.Update(*project)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &projectv1.UpdateProjectResponse{
		Success: updatedProject != nil,
		Id:      updatedProject.Id.String(),
	}, nil
}

func (s *ProjectService) Delete(ctx context.Context, req *projectv1.DeleteProjectRequest) (*projectv1.DeleteProjectResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var user model.User
	user.Id = claims.Id
	user.IsSuperAdmin = claims.IsSuperAdmin
	var reqOrg model.Org
	oId, err := uuid.Parse(req.OrgId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	reqOrg.Id = oId
	org, err := s.orgDao.Get(reqOrg)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	orgUsers, err := s.orgDao.ListMember(reqOrg)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if utils.IsOrgMember(org, orgUsers, user) == false {
		return nil, status.Error(codes.PermissionDenied, utils.ErrNotOrgMember)
	}
	var reqProject model.Project
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
	if project.Id == uuid.Nil {
		return nil, status.Error(codes.NotFound, "project not fount.")
	}
	if utils.IsProjectAdmin(*project, user) == false {
		if utils.IsOrgAdmin(org, orgUsers, user) == false {
			return nil, status.Error(codes.PermissionDenied, utils.ErrNoAdminPermissionForProject)
		}
	}
	deletedProject, err := s.projectDao.Delete(*project)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &projectv1.DeleteProjectResponse{
		Success: deletedProject,
		Id:      project.Id.String(),
	}, nil
}

func (s *ProjectService) SetAdmin(ctx context.Context, req *projectv1.SetAdminRequest) (*projectv1.UpdateProjectResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var user model.User
	user.Id = claims.Id
	user.IsSuperAdmin = claims.IsSuperAdmin
	var reqOrg model.Org
	oId, err := uuid.Parse(req.OrgId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	reqOrg.Id = oId
	org, err := s.orgDao.Get(reqOrg)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	orgUsers, err := s.orgDao.ListMember(reqOrg)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if utils.IsOrgMember(org, orgUsers, user) == false {
		return nil, status.Error(codes.PermissionDenied, utils.ErrNotOrgMember)
	}
	var reqProject model.Project
	pId, err := uuid.Parse(req.ProjectId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	reqProject.OrgId = oId
	reqProject.Id = pId
	reqProject.CreatedUser.Id = user.Id
	project, err := s.projectDao.Get(reqProject)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if project.Id == uuid.Nil {
		return nil, status.Error(codes.NotFound, "project not fount.")
	}
	if utils.IsProjectAdmin(*project, user) == false {
		if utils.IsOrgAdmin(org, orgUsers, user) == false {
			return nil, status.Error(codes.PermissionDenied, utils.ErrNoAdminPermissionForProject)
		}
	}
	var members []*projectv1.ProjectMember
	if project.Members != "" {
		err = json.Unmarshal([]byte(project.Members), &members)
		if err != nil {
			return nil, status.Error(codes.Unknown, err.Error())
		}
	}
	for _, m := range members {
		if m.UserId == req.UserId {
			m.IsAdmin = req.IsAdmin
			break
		}
	}
	membersJson, err := json.Marshal(members)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	project.Members = string(membersJson)
	updatedProject, err := s.projectDao.Update(*project)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &projectv1.UpdateProjectResponse{
		Success: true,
		Id:      updatedProject.Id.String(),
	}, nil
}

func (s *ProjectService) RemoveMember(ctx context.Context, req *projectv1.RemoveMemberRequest) (*projectv1.UpdateProjectResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var user model.User
	user.Id = claims.Id
	user.IsSuperAdmin = claims.IsSuperAdmin
	var reqOrg model.Org
	oId, err := uuid.Parse(req.OrgId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	reqOrg.Id = oId
	org, err := s.orgDao.Get(reqOrg)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	orgUsers, err := s.orgDao.ListMember(reqOrg)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if utils.IsOrgMember(org, orgUsers, user) == false {
		return nil, status.Error(codes.PermissionDenied, utils.ErrNotOrgMember)
	}
	var reqProject model.Project
	pId, err := uuid.Parse(req.ProjectId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	reqProject.OrgId = oId
	reqProject.Id = pId
	reqProject.CreatedUser.Id = user.Id
	project, err := s.projectDao.Get(reqProject)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if project.Id == uuid.Nil {
		return nil, status.Error(codes.NotFound, "project not fount.")
	}
	if utils.IsProjectAdmin(*project, user) == false {
		if utils.IsOrgAdmin(org, orgUsers, user) == false {
			return nil, status.Error(codes.PermissionDenied, utils.ErrNoAdminPermissionForProject)
		}
	}
	var members []*projectv1.ProjectMember
	if project.Members != "" {
		err = json.Unmarshal([]byte(project.Members), &members)
		if err != nil {
			return nil, status.Error(codes.Unknown, err.Error())
		}
	}
	var newMembers []*projectv1.ProjectMember
	for _, m := range members {
		if m.UserId != req.UserId {
			newMembers = append(newMembers, m)
		}
	}
	membersJson, err := json.Marshal(newMembers)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	project.Members = string(membersJson)
	updatedProject, err := s.projectDao.Update(*project)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &projectv1.UpdateProjectResponse{
		Success: true,
		Id:      updatedProject.Id.String(),
	}, nil
}

func (s *ProjectService) ListItem(ctx context.Context, req *itemv1.ListItemRequest) (*itemv1.ListItemResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var user model.User
	user.Id = claims.Id
	user.IsSuperAdmin = claims.IsSuperAdmin
	var reqOrg model.Org
	oId, err := uuid.Parse(req.OrgId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	reqOrg.Id = oId
	org, err := s.orgDao.Get(reqOrg)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	orgUsers, err := s.orgDao.ListMember(reqOrg)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if utils.IsOrgMember(org, orgUsers, user) == false {
		return nil, status.Error(codes.PermissionDenied, utils.ErrNotOrgMember)
	}
	var reqProject model.Project
	pId, err := uuid.Parse(req.ProjectId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	reqProject.OrgId = oId
	reqProject.Id = pId
	reqProject.CreatedUser.Id = user.Id
	project, err := s.projectDao.Get(reqProject)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if project.Id == uuid.Nil {
		return nil, status.Error(codes.NotFound, "project not fount.")
	}
	if utils.IsProjectMember(*project, user) == false {
		return nil, status.Error(codes.PermissionDenied, utils.ErrNotProjectMember)
	}
	var workItems []*model.WorkItem
	workItemCount := 0
	var tasks []*model.Task
	taskCount := 0
	if req.ShowAll == true {
		workItems, err = s.workItemDao.ListByProject(project.Id, req.Page, req.Size, req.Keyword)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		workItemCount = int(s.workItemDao.CountByProject(project.Id, req.Keyword))
		tasks, err = s.taskDao.ListByProject(project.Id, req.Page, req.Size, req.Keyword)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		taskCount = int(s.taskDao.CountByProject(project.Id, req.Keyword))
	} else {
		workItems, err = s.workItemDao.ListByProjectNotInSprint(project.Id, req.Page, req.Size, req.Keyword)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		workItemCount = int(s.workItemDao.CountByProjectNotInSprint(project.Id, req.Keyword))
		tasks, err = s.taskDao.ListByProjectNotInSprint(project.Id, req.Page, req.Size, req.Keyword)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		taskCount = int(s.taskDao.CountByProjectNotInSprint(project.Id, req.Keyword))
	}
	var items []*itemv1.Item
	for _, w := range workItems {
		items = append(items, &itemv1.Item{
			ItemType:    itemv1.Item_WORK_ITEM,
			Id:          w.Id,
			Title:       w.Title,
			Description: w.Description,
			ProjectId:   w.ProjectId.String(),
			SprintId:    w.SprintId.String(),
			Type:        w.Type,
			Status:      w.Status,
			AssignUser: &userv1.User{
				Id:    w.AssignUser.Id.String(),
				Name:  w.AssignUser.Name,
				Email: w.AssignUser.Email,
			},
			CreatedUser: &userv1.User{
				Id:    w.CreatedUser.Id.String(),
				Name:  w.CreatedUser.Name,
				Email: w.CreatedUser.Email,
			},
			CreatedAt: w.CreatedAt.Unix(),
			UpdatedAt: w.UpdatedAt.Unix(),
		})
	}
	for _, t := range tasks {
		items = append(items, &itemv1.Item{
			ItemType:    itemv1.Item_TASK,
			Id:          t.Id,
			Title:       t.Title,
			Description: t.Description,
			ProjectId:   t.ProjectId.String(),
			SprintId:    t.SprintId.String(),
			Type:        itemv1.Item_TASK.String(),
			Status:      t.Status,
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
			CreatedAt: t.CreatedAt.Unix(),
			UpdatedAt: t.UpdatedAt.Unix(),
		})
	}
	return &itemv1.ListItemResponse{
		Items: items,
		Pagination: &generalv1.Pagination{
			Page:  req.Page,
			Size:  req.Size,
			Total: int32(workItemCount + taskCount),
		},
	}, nil
}
