package service

import (
	"context"
	generalv1 "github.com/letscrum/letscrum/api/general/v1"
	v1 "github.com/letscrum/letscrum/api/letscrum/v1"
	projectv1 "github.com/letscrum/letscrum/api/project/v1"
	userV1 "github.com/letscrum/letscrum/api/user/v1"
	"github.com/letscrum/letscrum/internal/dao"
	"github.com/letscrum/letscrum/internal/model"
	"github.com/letscrum/letscrum/pkg/utils"
	"github.com/spf13/cast"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type ProjectService struct {
	v1.UnimplementedProjectServer
	projectDao       dao.ProjectDao
	userDao          dao.UserDao
	projectMemberDao dao.ProjectMemberDao
	sprintDao        dao.SprintDao
}

func NewProjectService(dao dao.Interface) *ProjectService {
	return &ProjectService{
		projectDao:       dao.ProjectDao(),
		userDao:          dao.UserDao(),
		projectMemberDao: dao.ProjectMemberDao(),
		sprintDao:        dao.SprintDao(),
	}
}

func (s *ProjectService) Get(ctx context.Context, req *projectv1.GetProjectRequest) (*projectv1.GetProjectResponse, error) {
	_, err := utils.AuthJWT(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	project, err := s.projectDao.Get(req.ProjectId)
	if err != nil {
		result := status.Convert(err)
		if result.Code() == codes.NotFound {
			return nil, status.Errorf(codes.NotFound, "get book err: %d not found", req.ProjectId)
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}
	if project.ID == 0 {
		return nil, status.Error(codes.NotFound, "project not fount.")
	}
	members, err := s.projectMemberDao.List(req.ProjectId, 1, 999)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	var memberList []*projectv1.ProjectMember
	for _, m := range members {
		var member = &projectv1.ProjectMember{
			Id:             m.ID,
			UserId:         m.UserID,
			ProjectId:      m.ProjectID,
			UserName:       m.User.Name,
			IsSuperAdmin:   m.User.IsSuperAdmin,
			IsProjectAdmin: m.IsAdmin,
		}
		memberList = append(memberList, member)
	}
	sprints, err := s.sprintDao.List(req.ProjectId, 1, 999, "")
	if err != nil {
		result := status.Convert(err)
		if result.Code() == codes.NotFound {
			return nil, status.Errorf(codes.NotFound, "not found.")
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}
	sprint := projectv1.Sprint{}
	for i, s := range sprints {
		if i == 0 {
			sprint = projectv1.Sprint{
				Id:        s.ID,
				ProjectId: s.ProjectID,
				Name:      s.Name,
				StartDate: s.StartDate.Unix(),
				EndDate:   s.EndDate.Unix(),
				Status:    projectv1.SprintStatus_UNKNOWN,
				CreatedAt: s.CreatedAt.Unix(),
				UpdatedAt: s.UpdatedAt.Unix(),
			}
		}
		if time.Now().After(s.StartDate) && time.Now().Before(s.EndDate) {
			sprint = projectv1.Sprint{
				Id:        s.ID,
				ProjectId: s.ProjectID,
				Name:      s.Name,
				StartDate: s.StartDate.Unix(),
				EndDate:   s.EndDate.Unix(),
				Status:    projectv1.SprintStatus_CURRENT,
				CreatedAt: s.CreatedAt.Unix(),
				UpdatedAt: s.UpdatedAt.Unix(),
			}
			break
		}
	}
	return &projectv1.GetProjectResponse{
		Item: &projectv1.Project{
			Id:          project.ID,
			Name:        project.Name,
			DisplayName: project.DisplayName,
			Description: project.Description,
			CreatedUser: &userV1.User{
				Id:           project.CreatedUser.ID,
				Name:         project.CreatedUser.Name,
				IsSuperAdmin: project.CreatedUser.IsSuperAdmin,
			},
			Members:   memberList,
			Sprint:    &sprint,
			CreatedAt: project.CreatedAt.Unix(),
			UpdatedAt: project.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *ProjectService) List(ctx context.Context, req *projectv1.ListProjectRequest) (*projectv1.ListProjectResponse, error) {
	_, err := utils.AuthJWT(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	projects, err := s.projectDao.List(req.Page, req.Size, req.Keyword)
	if err != nil {
		result := status.Convert(err)
		if result.Code() == codes.NotFound {
			return nil, status.Errorf(codes.NotFound, "not found.")
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}
	var list []*projectv1.Project
	for _, p := range projects {
		var project = &projectv1.Project{
			Id:          p.ID,
			Name:        p.Name,
			DisplayName: p.DisplayName,
			Description: p.Description,
			CreatedUser: &userV1.User{
				Id:           p.CreatedUser.ID,
				Name:         p.CreatedUser.Name,
				IsSuperAdmin: p.CreatedUser.IsSuperAdmin,
			},
			CreatedAt: p.CreatedAt.Unix(),
			UpdatedAt: p.UpdatedAt.Unix(),
		}
		members, err := s.projectMemberDao.List(p.ID, 1, 999)
		if err != nil {
			return nil, status.Error(codes.Unknown, err.Error())
		}
		for _, m := range members {
			var member = &projectv1.ProjectMember{
				Id:             m.ID,
				UserId:         m.UserID,
				ProjectId:      m.ProjectID,
				UserName:       m.User.Name,
				IsSuperAdmin:   m.User.IsSuperAdmin,
				IsProjectAdmin: m.IsAdmin,
			}
			project.Members = append(project.Members, member)
		}
		list = append(list, project)
	}
	count := s.projectDao.Count(req.Keyword)
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
	jwt, err := utils.AuthJWT(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	if !jwt.IsSuperAdmin {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	if req.DisplayName == "" {
		return nil, status.Error(codes.InvalidArgument, "project display name can't be empty.")
	}
	project := model.Project{
		Name:        req.DisplayName,
		DisplayName: req.DisplayName,
		Description: req.Description,
		CreatedBy:   cast.ToInt64(jwt.Id),
	}
	id, err := s.projectDao.Create(&project)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	success := false
	if id > 0 {
		success = true
	}
	var userIDs []int64
	for _, u := range req.Members {
		if u != project.CreatedBy {
			userIDs = append(userIDs, u)
		}
	}
	if len(userIDs) > 0 {
		successMembers, err := s.projectMemberDao.Add(id, userIDs)
		if err != nil {
			return nil, status.Error(codes.Unknown, err.Error())
		}
		success = successMembers
	}
	return &projectv1.CreateProjectResponse{
		Success: success,
	}, nil
}

func (s *ProjectService) Update(ctx context.Context, req *projectv1.UpdateProjectRequest) (*projectv1.UpdateProjectResponse, error) {
	jwt, err := utils.AuthJWT(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	myMember, err := s.projectMemberDao.Get(req.ProjectId, cast.ToInt64(jwt.Id))
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	if !myMember.IsAdmin || !jwt.IsSuperAdmin {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	project, err := s.projectDao.Get(req.ProjectId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	project.DisplayName = req.DisplayName
	project.Description = req.Description
	success, err := s.projectDao.Update(project)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &projectv1.UpdateProjectResponse{
		Success: success,
	}, nil
}

func (s *ProjectService) Delete(ctx context.Context, req *projectv1.DeleteProjectRequest) (*projectv1.DeleteProjectResponse, error) {
	jwt, err := utils.AuthJWT(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	user, err := s.userDao.Get(cast.ToInt64(jwt.Id))
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if !user.IsSuperAdmin {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	success, err := s.projectDao.Delete(req.ProjectId)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &projectv1.DeleteProjectResponse{
		Success: success,
	}, nil
}

//
//func Create1(name string, displayName string, createdUserId int64) (int64, error) {
//	projectId, err := model.CreateProject(name, displayName, createdUserId)
//	if err != nil {
//		return 0, err
//	}
//	_, errCreateMember := models2.CreateProjectMember(projectId, createdUserId, true)
//	if errCreateMember != nil {
//		errDeleteProject := model.DeleteProject(projectId)
//		if errDeleteProject != nil {
//			return 0, errDeleteProject
//		}
//		return 0, errCreateMember
//	}
//	return projectId, nil
//}
//
//func List1(page int32, pageSize int32) ([]*projectv1.Project, int64, error) {
//	projects, err := model.ListProject(page, pageSize)
//	if err != nil {
//		return nil, 0, err
//	}
//	var list []*projectv1.Project
//	for _, p := range projects {
//		list = append(list, &projectv1.Project{
//			Id:          p.Id,
//			Name:        p.Name,
//			DisplayName: p.DisplayName,
//			CreatedUser: &userV1.User{
//				Id:   p.CreatedUser.Id,
//				Name: p.CreatedUser.Name,
//			},
//			CreatedAt: p.CreatedAt.Unix(),
//			UpdatedAt: p.UpdatedAt.Unix(),
//		})
//	}
//	count := model.CountProject()
//	return list, count, nil
//}
//
//func Update1(id int64, displayName string) error {
//	if err := model.UpdateProject(id, displayName); err != nil {
//		return err
//	}
//	return nil
//}
//
//func Delete1(id int64) error {
//	if err := model.DeleteProject(id); err != nil {
//		return err
//	}
//	return nil
//}
//
//func HardDelete1(id int64) error {
//	if err := models2.HardDeleteProjectMemberByProject(id); err != nil {
//		return err
//	}
//	if err := model.HardDeleteProject(id); err != nil {
//		return err
//	}
//	return nil
//}
