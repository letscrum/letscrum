package service

import (
	"context"
	v1 "github.com/letscrum/letscrum/api/letscrum/v1"
	projectV1 "github.com/letscrum/letscrum/api/project/v1"
	userV1 "github.com/letscrum/letscrum/api/user/v1"
	"github.com/letscrum/letscrum/internal/dao"
	"github.com/letscrum/letscrum/internal/model"
	models2 "github.com/letscrum/letscrum/internal/model/projectmembermodel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProjectServiceInterface interface {
	Get(context.Context, *projectV1.GetProjectRequest) (*projectV1.GetProjectResponse, error)
}

type ProjectService struct {
	v1.UnimplementedProjectServer
	dao dao.ProjectDao
}

func NewProjectService(dao dao.Interface) *ProjectService {
	return &ProjectService{dao: dao.ProjectDao()}
}

func (s *ProjectService) Get(ctx context.Context, req *projectV1.GetProjectRequest) (*projectV1.GetProjectResponse, error) {
	project, err := s.dao.Get(ctx, req.ProjectId)
	if err != nil {
		result := status.Convert(err)
		if result.Code() == codes.NotFound {
			return nil, status.Errorf(codes.NotFound, "get book err: %s not found", req.ProjectId)
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &projectV1.GetProjectResponse{
		Item: &projectV1.Project{
			Id: project.Id,
		},
	}, nil
}

func Create(name string, displayName string, createdUserId int64) (int64, error) {
	projectId, err := model.CreateProject(name, displayName, createdUserId)
	if err != nil {
		return 0, err
	}
	_, errCreateMember := models2.CreateProjectMember(projectId, createdUserId, true)
	if errCreateMember != nil {
		errDeleteProject := model.DeleteProject(projectId)
		if errDeleteProject != nil {
			return 0, errDeleteProject
		}
		return 0, errCreateMember
	}
	return projectId, nil
}

func List(page int32, pageSize int32) ([]*projectV1.Project, int64, error) {
	projects, err := model.ListProject(page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	var list []*projectV1.Project
	for _, p := range projects {
		list = append(list, &projectV1.Project{
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
	count := model.CountProject()
	return list, count, nil
}

func Update(id int64, displayName string) error {
	if err := model.UpdateProject(id, displayName); err != nil {
		return err
	}
	return nil
}

func Delete(id int64) error {
	if err := model.DeleteProject(id); err != nil {
		return err
	}
	return nil
}

func HardDelete(id int64) error {
	if err := models2.HardDeleteProjectMemberByProject(id); err != nil {
		return err
	}
	if err := model.HardDeleteProject(id); err != nil {
		return err
	}
	return nil
}

func Get(id int64) (*projectV1.Project, error) {
	p, err := model.GetProject(id)
	if err != nil {
		return nil, err
	}
	project := &projectV1.Project{
		Id:          p.Id,
		Name:        p.Name,
		DisplayName: p.DisplayName,
		CreatedUser: &userV1.User{
			Id:   p.CreatedUser.Id,
			Name: p.CreatedUser.Name,
		},
		CreatedAt: p.CreatedAt.Unix(),
		UpdatedAt: p.UpdatedAt.Unix(),
	}
	return project, nil
}
