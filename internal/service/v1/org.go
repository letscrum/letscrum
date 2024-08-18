package v1

import (
	"context"

	"github.com/google/uuid"
	generalv1 "github.com/letscrum/letscrum/api/general/v1"
	v1 "github.com/letscrum/letscrum/api/letscrum/v1"
	orgv1 "github.com/letscrum/letscrum/api/org/v1"
	userv1 "github.com/letscrum/letscrum/api/user/v1"
	"github.com/letscrum/letscrum/internal/dao"
	"github.com/letscrum/letscrum/internal/model"
	"github.com/letscrum/letscrum/pkg/utils"
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
		return nil, status.Error(codes.PermissionDenied, utils.ErrReachOrgLimit)
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
	members, err := s.orgDao.ListMember(reqOrg)
	if org.CreatedBy != reqUser.Id {
		if utils.IsOrgMember(org, members, reqUser) == false {
			return nil, status.Error(codes.PermissionDenied, utils.ErrNotOrgMember)
		}
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var memberItems []*orgv1.OrgMember
	for _, m := range members {
		memberItems = append(memberItems, &orgv1.OrgMember{
			Member: &userv1.User{
				Id:    m.UserId.String(),
				Name:  m.Member.Name,
				Email: m.Member.Email,
			},
			IsAdmin: m.IsAdmin,
		})
	}

	return &orgv1.OrgResponse{
		Item: &orgv1.Org{
			Id:          org.Id.String(),
			Name:        org.Name,
			CreatedBy:   org.CreatedUser.Name,
			MemberCount: int32(len(memberItems)),
			Members:     memberItems,
			MyRole:      utils.GetOrgRole(org, members, reqUser),
		},
	}, nil
}

func (s OrgService) Update(ctx context.Context, req *orgv1.UpdateOrgRequest) (*orgv1.OrgResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s OrgService) Delete(ctx context.Context, req *orgv1.DeleteOrgRequest) (*orgv1.DeleteOrgResponse, error) {
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
		if utils.IsOrgAdmin(org, orgUsers, reqUser) == false {
			return nil, status.Error(codes.PermissionDenied, utils.ErrNotOrgAdmin)
		}
	}

	success, err := s.orgDao.Delete(reqOrg)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &orgv1.DeleteOrgResponse{
		Success: success,
	}, nil
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
			Id:          org.Id.String(),
			Name:        org.Name,
			CreatedBy:   org.CreatedUser.Name,
			MemberCount: 0,
			Members:     nil,
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

func (s OrgService) AddMembers(ctx context.Context, req *orgv1.AddMembersRequest) (*orgv1.ListMemberResponse, error) {
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
		if utils.IsOrgAdmin(org, orgUsers, reqUser) == false {
			return nil, status.Error(codes.PermissionDenied, utils.ErrNotOrgAdmin)
		}
	}
	var members []model.OrgUser
	for _, m := range req.Members {
		if m.UserId != org.CreatedBy.String() {
			members = append(members, model.OrgUser{
				OrgId:   org.Id,
				UserId:  uuid.MustParse(m.UserId),
				IsAdmin: m.IsAdmin,
			})
		}
	}
	if len(members) > 0 {
		_, err := s.orgDao.AddMembers(members)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	orgUsers, err := s.orgDao.ListMember(org)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var memberItems []*orgv1.OrgMember
	for _, m := range orgUsers {
		memberItems = append(memberItems, &orgv1.OrgMember{
			Member: &userv1.User{
				Id:    m.UserId.String(),
				Name:  m.Member.Name,
				Email: m.Member.Email,
			},
			IsAdmin: m.IsAdmin,
		})
	}
	return &orgv1.ListMemberResponse{
		Items: memberItems,
	}, nil
}

func (s OrgService) RemoveMember(ctx context.Context, req *orgv1.RemoveMemberRequest) (*orgv1.ListMemberResponse, error) {
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
		if utils.IsOrgAdmin(org, orgUsers, reqUser) == false {
			return nil, status.Error(codes.PermissionDenied, utils.ErrNotOrgAdmin)
		}
	}
	var member model.OrgUser
	uId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	member.OrgId = org.Id
	member.UserId = uId
	success, err := s.orgDao.RemoveMember(member)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !success {
		return nil, status.Error(codes.Internal, "failed to remove member")
	}
	orgUsers, err := s.orgDao.ListMember(org)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var memberItems []*orgv1.OrgMember
	for _, m := range orgUsers {
		memberItems = append(memberItems, &orgv1.OrgMember{
			Member: &userv1.User{
				Id:    m.UserId.String(),
				Name:  m.Member.Name,
				Email: m.Member.Email,
			},
			IsAdmin: m.IsAdmin,
		})
	}
	return &orgv1.ListMemberResponse{
		Items: memberItems,
	}, nil
}

func (s OrgService) SetAdmin(ctx context.Context, req *orgv1.SetAdminRequest) (*orgv1.ListMemberResponse, error) {
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
		if utils.IsOrgAdmin(org, orgUsers, reqUser) == false {
			return nil, status.Error(codes.PermissionDenied, utils.ErrNotOrgAdmin)
		}
	}
	var member model.OrgUser
	uId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	member.OrgId = org.Id
	member.UserId = uId
	member, err = s.orgDao.SetAdmin(member, req.IsAdmin)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	orgUsers, err := s.orgDao.ListMember(org)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var memberItems []*orgv1.OrgMember
	for _, m := range orgUsers {
		memberItems = append(memberItems, &orgv1.OrgMember{
			Member: &userv1.User{
				Id:    m.UserId.String(),
				Name:  m.Member.Name,
				Email: m.Member.Email,
			},
			IsAdmin: m.IsAdmin,
		})
	}
	return &orgv1.ListMemberResponse{
		Items: memberItems,
	}, nil
}

func (s OrgService) ListMember(ctx context.Context, req *orgv1.ListMemberRequest) (*orgv1.ListMemberResponse, error) {
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
		if utils.IsOrgMember(org, orgUsers, reqUser) == false {
			return nil, status.Error(codes.PermissionDenied, utils.ErrNotOrgMember)
		}
	}
	orgUsers, err := s.orgDao.ListMember(org)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var memberItems []*orgv1.OrgMember
	memberItems = append(memberItems, &orgv1.OrgMember{
		Member: &userv1.User{
			Id:    org.CreatedBy.String(),
			Name:  org.CreatedUser.Name,
			Email: org.CreatedUser.Email,
		},
		IsAdmin: true,
	})
	for _, m := range orgUsers {
		memberItems = append(memberItems, &orgv1.OrgMember{
			Member: &userv1.User{
				Id:    m.UserId.String(),
				Name:  m.Member.Name,
				Email: m.Member.Email,
			},
			IsAdmin: m.IsAdmin,
		})
	}
	return &orgv1.ListMemberResponse{
		Items: memberItems,
	}, nil
}
