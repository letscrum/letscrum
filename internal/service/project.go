package projectservice

import (
	"context"
	v1 "github.com/letscrum/letscrum/api/letscrum/v1"
	projectV1 "github.com/letscrum/letscrum/api/project/v1"
	userV1 "github.com/letscrum/letscrum/api/user/v1"
	"github.com/letscrum/letscrum/internal/dao"
	models2 "github.com/letscrum/letscrum/internal/model/projectmembermodel"
	"github.com/letscrum/letscrum/internal/model/projectmodel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProjectService interface {
	Get(context.Context, *projectV1.GetProjectRequest) (*projectV1.GetProjectResponse, error)
}

type ProjectServer struct {
	v1.UnimplementedProjectServer
	dao dao.LetscrumDao
}

func NewProjectService(letscrumService *LetscrumService) *ProjectService {
	return &ProjectServer{dao: hiveService.dao}
}

func (s *bookService) Get(ctx context.Context, req *v1alpha1.GetBookReq) (*v1alpha1.GetBookReply, error) {
	book, err := s.dao.BookDao().Get(ctx, req.Uid, metav1alpha1.GetOptions{})
	if err != nil {
		result := status.Convert(err)
		if result.Code() == codes.NotFound {
			return nil, status.Errorf(codes.NotFound, "get book err: %s not found", req.Uid)
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &v1alpha1.GetBookReply{
		Uid:         book.UID,
		Name:        book.Name,
		Author:      book.Author,
		Status:      book.Status,
		IsPublished: book.IsPublished,
		PublishedAt: timestamppb.New(book.PublishedAt),
	}, nil
}

func Create(name string, displayName string, createdUserId int64) (int64, error) {
	projectId, err := projectmodel.CreateProject(name, displayName, createdUserId)
	if err != nil {
		return 0, err
	}
	_, errCreateMember := models2.CreateProjectMember(projectId, createdUserId, true)
	if errCreateMember != nil {
		errDeleteProject := projectmodel.DeleteProject(projectId)
		if errDeleteProject != nil {
			return 0, errDeleteProject
		}
		return 0, errCreateMember
	}
	return projectId, nil
}

func List(page int32, pageSize int32) ([]*projectV1.Project, int64, error) {
	projects, err := projectmodel.ListProject(page, pageSize)
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
	count := projectmodel.CountProject()
	return list, count, nil
}

func Update(id int64, displayName string) error {
	if err := projectmodel.UpdateProject(id, displayName); err != nil {
		return err
	}
	return nil
}

func Delete(id int64) error {
	if err := projectmodel.DeleteProject(id); err != nil {
		return err
	}
	return nil
}

func HardDelete(id int64) error {
	if err := models2.HardDeleteProjectMemberByProject(id); err != nil {
		return err
	}
	if err := projectmodel.HardDeleteProject(id); err != nil {
		return err
	}
	return nil
}

func Get(id int64) (*projectV1.Project, error) {
	p, err := projectmodel.GetProject(id)
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
