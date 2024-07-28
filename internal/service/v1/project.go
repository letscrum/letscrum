package v1

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	generalv1 "github.com/letscrum/letscrum/api/general/v1"
	v1 "github.com/letscrum/letscrum/api/letscrum/v1"
	projectv1 "github.com/letscrum/letscrum/api/project/v1"
	userv1 "github.com/letscrum/letscrum/api/user/v1"
	"github.com/letscrum/letscrum/internal/dao"
	"github.com/letscrum/letscrum/internal/model"
	"github.com/letscrum/letscrum/pkg/utils"
	"github.com/letscrum/letscrum/pkg/validator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProjectService struct {
	v1.UnimplementedProjectServer
	projectDao dao.ProjectDao
	userDao    dao.UserDao
	sprintDao  dao.SprintDao
	orgDao     dao.OrgDao
}

func NewProjectService(dao dao.Interface) *ProjectService {
	return &ProjectService{
		projectDao: dao.ProjectDao(),
		userDao:    dao.UserDao(),
		sprintDao:  dao.SprintDao(),
		orgDao:     dao.OrgDao(),
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
	orgUsers, err := s.orgDao.ListMember(reqOrg)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if validator.IsOrgMember(orgUsers, user) == false {
		return nil, status.Error(codes.PermissionDenied, "You are not a member of this organization")
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
	if validator.IsProjectMember(*project, user) == false {
		return nil, status.Error(codes.PermissionDenied, "You are not a member of this project")
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
		return nil, status.Error(codes.PermissionDenied, "You have reached the maximum number of projects.")
	}

	if org.CreatedBy != user.Id {
		orgUsers, err := s.orgDao.ListMember(reqOrg)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		if validator.IsOrgMember(orgUsers, *user) == false {
			return nil, status.Error(codes.PermissionDenied, "You are not a member of this organization")
		}
	}

	if req.DisplayName == "" {
		return nil, status.Error(codes.InvalidArgument, "project display name can't be empty.")
	}
	var members []*projectv1.ProjectMember
	// add current user as project admin
	members = append(members, &projectv1.ProjectMember{
		UserId:   user.Id.String(),
		UserName: user.Name,
		IsAdmin:  user.IsSuperAdmin,
	})
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
			return nil, status.Error(codes.Unknown, err.Error())
		}
		// add members to project members
		for _, u := range users {
			isAdmin := false
			if u.IsSuperAdmin == true {
				isAdmin = true
			} else {
				for _, m := range req.Members {
					if m.UserId == u.Id.String() {
						isAdmin = m.IsAdmin
						break
					}
				}
			}
			member := &projectv1.ProjectMember{
				UserId:   u.Id.String(),
				UserName: u.Name,
				IsAdmin:  isAdmin,
			}
			members = append(members, member)
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
	newProject.Name = req.DisplayName
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
	orgUsers, err := s.orgDao.ListMember(reqOrg)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if validator.IsOrgMember(orgUsers, user) == false {
		return nil, status.Error(codes.PermissionDenied, "You are not a member of this organization")
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
	if validator.IsProjectAdmin(*project, user) == false {
		return nil, status.Error(codes.PermissionDenied, "You are not a admin of this project")
	}
	project.DisplayName = req.DisplayName
	project.Description = req.Description
	var members []*projectv1.ProjectMember
	if req.Members != nil && len(req.Members) > 0 {
		for _, m := range req.Members {
			member := &projectv1.ProjectMember{
				UserId:   m.UserId,
				UserName: m.UserName,
				IsAdmin:  m.IsAdmin,
			}
			if m.UserId == user.Id.String() && user.IsSuperAdmin == true {
				member.IsAdmin = true
			}
			members = append(members, member)
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
	orgUsers, err := s.orgDao.ListMember(reqOrg)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if validator.IsOrgMember(orgUsers, user) == false {
		return nil, status.Error(codes.PermissionDenied, "You are not a member of this organization")
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
	if validator.IsProjectAdmin(*project, user) == false {
		return nil, status.Error(codes.PermissionDenied, "You are not a admin of this project")
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
