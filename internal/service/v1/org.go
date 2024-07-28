package v1

import (
	"context"

	"github.com/google/uuid"
	generalv1 "github.com/letscrum/letscrum/api/general/v1"
	v1 "github.com/letscrum/letscrum/api/letscrum/v1"
	orgv1 "github.com/letscrum/letscrum/api/org/v1"
	"github.com/letscrum/letscrum/internal/dao"
	"github.com/letscrum/letscrum/internal/model"
	"github.com/letscrum/letscrum/pkg/utils"
	"github.com/letscrum/letscrum/pkg/validator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrgService struct {
	v1.UnimplementedOrgServer
	orgDao  dao.OrgDao
	userDao dao.UserDao
}

func NewOrgService(dao dao.Interface) *OrgService {
	return &OrgService{
		orgDao:  dao.OrgDao(),
		userDao: dao.UserDao(),
	}
}

func (s OrgService) Create(ctx context.Context, req *orgv1.CreateOrgRequest) (*orgv1.CreateOrgResponse, error) {
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
	orgCount := s.orgDao.CountByUser(*user)
	if orgCount >= user.OrgLimitation {
		return nil, status.Error(codes.PermissionDenied, "You have reached the maximum number of organizations.")
	}

	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "organization name can't be empty.")
	}

	var newOrg model.Org
	newOrg.Id = uuid.New()
	newOrg.Name = req.Name
	newOrg.CreatedBy = user.Id

	org, err := s.orgDao.Create(newOrg)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &orgv1.CreateOrgResponse{
		Success: org.Id != uuid.Nil,
		Id:      org.Id.String(),
	}, nil
}

func (s OrgService) Get(ctx context.Context, req *orgv1.GetOrgRequest) (*orgv1.OrgResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var reqUser model.User
	reqUser.Id = claims.Id
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

	if org.CreatedBy != reqUser.Id {
		orgUsers, err := s.orgDao.ListMember(reqOrg)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		if validator.IsOrgMember(orgUsers, reqUser) == false {
			return nil, status.Error(codes.PermissionDenied, "You are not a member of this organization")
		}
	}

	return &orgv1.OrgResponse{
		Item: &orgv1.Org{
			Id:   org.Id.String(),
			Name: org.Name,
		},
	}, nil
}

func (s OrgService) Update(ctx context.Context, req *orgv1.UpdateOrgRequest) (*orgv1.OrgResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s OrgService) List(ctx context.Context, req *orgv1.ListOrgRequest) (*orgv1.ListOrgResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var reqUser model.User
	reqUser.Id = claims.Id
	reqUser.IsSuperAdmin = claims.IsSuperAdmin
	req.Page, req.Size = utils.Pagination(req.Page, req.Size)
	orgs, err := s.orgDao.ListVisibleOrg(req.Page, req.Size, req.Keyword, reqUser)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var orgItems []*orgv1.Org
	for _, org := range orgs {
		orgItems = append(orgItems, &orgv1.Org{
			Id:   org.Id.String(),
			Name: org.Name,
		})
	}
	count := s.orgDao.CountVisibleOrg(req.Keyword, reqUser)
	return &orgv1.ListOrgResponse{
		Items: orgItems,
		Pagination: &generalv1.Pagination{
			Page:  req.Page,
			Size:  req.Size,
			Total: int32(count),
		},
	}, nil
}

func (s OrgService) AddMembers(ctx context.Context, req *orgv1.AddMembersRequest) (*orgv1.MemberResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s OrgService) RemoveMember(ctx context.Context, req *orgv1.RemoveMemberRequest) (*orgv1.MemberResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s OrgService) SetAdmin(ctx context.Context, req *orgv1.SetAdminRequest) (*orgv1.MemberResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s OrgService) ListMember(ctx context.Context, req *orgv1.ListMemberRequest) (*orgv1.ListMemberResponse, error) {
	//TODO implement me
	panic("implement me")
}
