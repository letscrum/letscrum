package v1

import (
	"context"

	v1 "github.com/letscrum/letscrum/api/letscrum/v1"
	projectv1 "github.com/letscrum/letscrum/api/project/v1"
	"github.com/letscrum/letscrum/internal/dao"
	"github.com/letscrum/letscrum/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	_, err := utils.AuthJWT(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	return nil, nil
}
