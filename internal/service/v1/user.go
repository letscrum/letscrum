package v1

import (
	"context"

	"github.com/google/uuid"
	generalv1 "github.com/letscrum/letscrum/api/general/v1"
	letscrumv1 "github.com/letscrum/letscrum/api/letscrum/v1"
	userv1 "github.com/letscrum/letscrum/api/user/v1"
	"github.com/letscrum/letscrum/internal/dao"
	"github.com/letscrum/letscrum/internal/model"
	"github.com/letscrum/letscrum/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	letscrumv1.UnimplementedUserServer
	userDao dao.UserDao
}

func NewUserService(dao dao.Interface) *UserService {
	return &UserService{userDao: dao.UserDao()}
}

func (s *UserService) List(_ context.Context, req *userv1.ListUserRequest) (*userv1.ListUserResponse, error) {
	users, err := s.userDao.List(req.Page, req.Size, req.Keyword)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var list []*userv1.User
	for _, u := range users {
		list = append(list, &userv1.User{
			Id:           u.Id.String(),
			Name:         u.Name,
			Email:        u.Email,
			IsSuperAdmin: u.IsSuperAdmin,
			CreatedAt:    u.CreatedAt.Unix(),
			UpdatedAt:    u.UpdatedAt.Unix(),
		})
	}
	count := s.userDao.Count(req.Keyword)
	return &userv1.ListUserResponse{
		Items: list,
		Pagination: &generalv1.Pagination{
			Page:  req.Page,
			Size:  req.Size,
			Total: int32(count),
		},
	}, nil
}

func (s *UserService) Create(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var user model.User
	user.Id = claims.Id
	user.IsSuperAdmin = claims.IsSuperAdmin
	if user.IsSuperAdmin == false {
		return nil, status.Error(codes.PermissionDenied, "You are not a super admin")
	}
	var reqUser model.User
	reqUser.Id = uuid.New()
	reqUser.Name = req.Name
	reqUser.Email = req.Email
	reqUser.Password = req.Password
	reqUser.IsSuperAdmin = req.IsSuperAdmin
	newUser, err := s.userDao.Create(reqUser)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &userv1.CreateUserResponse{
		Success: newUser.Id != uuid.Nil,
		Item: &userv1.User{
			Id:           newUser.Id.String(),
			Name:         newUser.Name,
			Email:        newUser.Email,
			IsSuperAdmin: newUser.IsSuperAdmin,
			CreatedAt:    newUser.CreatedAt.Unix(),
			UpdatedAt:    newUser.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *UserService) Get(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *UserService) SetSuperAdmin(ctx context.Context, req *userv1.SetSuperAdminRequest) (*userv1.UpdateUserResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var user model.User
	user.Id = claims.Id
	user.IsSuperAdmin = claims.IsSuperAdmin
	if user.IsSuperAdmin == false {
		return nil, status.Error(codes.PermissionDenied, "You are not a super admin")
	}
	// if current user is only super admin can't set to false
	admins, err := s.userDao.ListSuperAdmins()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if len(admins) == 1 && admins[0].Id.String() == req.UserId && req.IsSuperAdmin == false {
		return nil, status.Error(codes.PermissionDenied, "Can't set the only super admin to false")
	}
	var reqUser model.User
	reqUser.Id = uuid.MustParse(req.UserId)
	updatedUser, err := s.userDao.SetSuperAdmin(reqUser, req.IsSuperAdmin)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &userv1.UpdateUserResponse{
		Success: updatedUser.Id != uuid.Nil,
		Item: &userv1.User{
			Id:           updatedUser.Id.String(),
			Name:         updatedUser.Name,
			Email:        updatedUser.Email,
			IsSuperAdmin: updatedUser.IsSuperAdmin,
			CreatedAt:    updatedUser.CreatedAt.Unix(),
			UpdatedAt:    updatedUser.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *UserService) UpdatePassword(ctx context.Context, req *userv1.UpdatePasswordRequest) (*userv1.UpdateUserResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var user model.User
	user.Id = claims.Id
	user.IsSuperAdmin = claims.IsSuperAdmin
	if user.Id.String() != req.UserId {
		return nil, status.Error(codes.PermissionDenied, "You are not allowed to update password for other user")
	}
	var reqUser model.User
	reqUser.Id = uuid.MustParse(req.UserId)
	reqUser.Password = req.OldPassword
	updatedUser, err := s.userDao.UpdatePassword(reqUser, req.NewPassword)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &userv1.UpdateUserResponse{
		Success: updatedUser.Id != uuid.Nil,
		Item: &userv1.User{
			Id:           updatedUser.Id.String(),
			Name:         updatedUser.Name,
			Email:        updatedUser.Email,
			IsSuperAdmin: updatedUser.IsSuperAdmin,
			CreatedAt:    updatedUser.CreatedAt.Unix(),
			UpdatedAt:    updatedUser.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *UserService) ResetPassword(ctx context.Context, req *userv1.ResetPasswordRequest) (*userv1.UpdateUserResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var user model.User
	user.Id = claims.Id
	user.IsSuperAdmin = claims.IsSuperAdmin
	if user.IsSuperAdmin == false {
		return nil, status.Error(codes.PermissionDenied, "You are not a super admin")
	}
	var reqUser model.User
	reqUser.Id = uuid.MustParse(req.UserId)
	updatedUser, err := s.userDao.ResetPassword(reqUser, req.NewPassword)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &userv1.UpdateUserResponse{
		Success: updatedUser.Id != uuid.Nil,
		Item: &userv1.User{
			Id:           updatedUser.Id.String(),
			Name:         updatedUser.Name,
			Email:        updatedUser.Email,
			IsSuperAdmin: updatedUser.IsSuperAdmin,
			CreatedAt:    updatedUser.CreatedAt.Unix(),
			UpdatedAt:    updatedUser.UpdatedAt.Unix(),
		},
	}, nil
}
