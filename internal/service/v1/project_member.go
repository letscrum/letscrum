package v1

import (
	"context"

	v1 "github.com/letscrum/letscrum/api/letscrum/v1"
	projectv1 "github.com/letscrum/letscrum/api/project/v1"
	"github.com/letscrum/letscrum/internal/dao"
)

type ProjectMemberService struct {
	v1.UnimplementedProjectMemberServer
	projectMemberDao dao.ProjectMemberDao
}

func NewProjectMemberService(dao dao.Interface) *ProjectMemberService {
	return &ProjectMemberService{
		projectMemberDao: dao.ProjectMemberDao(),
	}
}

func (s *ProjectMemberService) List(ctx context.Context, req *projectv1.ListProjectMemberRequest) (*projectv1.ListProjectMemberResponse, error) {
	return nil, nil
}
