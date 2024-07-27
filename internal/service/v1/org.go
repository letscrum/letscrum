package v1

import (
	"context"

	v1 "github.com/letscrum/letscrum/api/letscrum/v1"
	orgv1 "github.com/letscrum/letscrum/api/org/v1"
	"github.com/letscrum/letscrum/internal/dao"
)

type OrgService struct {
	v1.UnimplementedOrgServer
	orgDao dao.OrgDao
}

func NewOrgService(dao dao.Interface) *OrgService {
	return &OrgService{
		orgDao: dao.OrgDao(),
	}
}

func (s OrgService) Create(ctx context.Context, req *orgv1.CreateOrgRequest) (*orgv1.CreateOrgResponse, error) {
	//TODO implement me
	panic("implement me")
}
