package service

import (
	"context"
	generalv1 "github.com/letscrum/letscrum/api/general/v1"
	v1 "github.com/letscrum/letscrum/api/letscrum/v1"
	projectv1 "github.com/letscrum/letscrum/api/project/v1"
	userV1 "github.com/letscrum/letscrum/api/user/v1"
	"github.com/letscrum/letscrum/internal/dao"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProjectService struct {
	v1.UnimplementedProjectServer
	dao dao.ProjectDao
}

func NewProjectService(dao dao.Interface) *ProjectService {
	return &ProjectService{dao: dao.ProjectDao()}
}

func (s *ProjectService) Get(ctx context.Context, req *projectv1.GetProjectRequest) (*projectv1.GetProjectResponse, error) {
	project, err := s.dao.Get(ctx, req.ProjectId)
	if err != nil {
		result := status.Convert(err)
		if result.Code() == codes.NotFound {
			return nil, status.Errorf(codes.NotFound, "get book err: %s not found", req.ProjectId)
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &projectv1.GetProjectResponse{
		Item: &projectv1.Project{
			Id:          project.Id,
			Name:        project.Name,
			DisplayName: project.DisplayName,
			CreatedUser: &userV1.User{
				Id:   project.CreatedUser.Id,
				Name: project.CreatedUser.Name,
			},
			CreatedAt: project.CreatedAt.Unix(),
			UpdatedAt: project.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *ProjectService) List(ctx context.Context, req *projectv1.ListProjectRequest) (*projectv1.ListProjectResponse, error) {
	projects, err := s.dao.List(ctx, req.Page, req.PageSize)
	if err != nil {
		result := status.Convert(err)
		if result.Code() == codes.NotFound {
			return nil, status.Errorf(codes.NotFound, "get book err: not found")
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}
	var list []*projectv1.Project
	for _, p := range projects {
		list = append(list, &projectv1.Project{
			Id:          p.Id,
			Name:        p.Name,
			DisplayName: p.DisplayName,
			CreatedUser: &userV1.User{
				Id:   p.CreatedUser.Id,
				Name: p.CreatedUser.Name,
			},
			CreatedAt: p.CreatedAt.Unix(),
			UpdatedAt: p.UpdatedAt.Unix(),
		})
	}
	count := s.dao.Count(ctx)
	return &projectv1.ListProjectResponse{
		Items: list,
		Pagination: &generalv1.Pagination{
			Page:     req.Page,
			PageSize: req.PageSize,
			Total:    int32(count),
		},
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
