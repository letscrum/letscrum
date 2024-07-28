package v1

import (
	"context"

	"github.com/google/uuid"
	generalv1 "github.com/letscrum/letscrum/api/general/v1"
	letscrumv1 "github.com/letscrum/letscrum/api/letscrum/v1"
	userv1 "github.com/letscrum/letscrum/api/user/v1"
	"github.com/letscrum/letscrum/internal/dao"
	"github.com/letscrum/letscrum/pkg/build"
	"github.com/letscrum/letscrum/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type LetscrumService struct {
	letscrumv1.UnimplementedLetscrumServer
	userDao dao.UserDao
}

func NewLetscrumService(dao dao.Interface) *LetscrumService {
	return &LetscrumService{userDao: dao.UserDao()}
}

func (s *LetscrumService) GetVersion(context.Context, *emptypb.Empty) (*generalv1.GetVersionResponse, error) {
	ver := build.Version()
	return &generalv1.GetVersionResponse{
		Version: &generalv1.Version{
			Version:   ver.Version,
			GitCommit: ver.GitCommit,
			BuildDate: ver.BuildDate,
			GoVersion: ver.GoVersion,
		},
	}, nil
}

func (s *LetscrumService) SignIn(_ context.Context, req *userv1.SignInRequest) (*userv1.SignInResponse, error) {
	user, err := s.userDao.SignIn(req.Name, req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to sign in")
	}
	if user.Id == uuid.Nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	accessToken, refreshToken, errGenTokens := utils.GenerateTokens(user.Id, user.IsSuperAdmin)
	if errGenTokens != nil {
		return nil, errGenTokens
	}
	return &userv1.SignInResponse{
		Item: &userv1.User{
			Id:           user.Id.String(),
			Name:         user.Name,
			Email:        user.Email,
			IsSuperAdmin: user.IsSuperAdmin,
			CreatedAt:    user.CreatedAt.Unix(),
			UpdatedAt:    user.UpdatedAt.Unix(),
			Token: &userv1.Token{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			},
		},
	}, nil
}

func (s *LetscrumService) RefreshToken(ctx context.Context, req *userv1.RefreshTokenRequest) (*userv1.RefreshTokenResponse, error) {
	token, err := utils.ParseToken(req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid token")
	}
	accessToken, refreshToken, err := utils.GenerateTokens(token["iss"].(uuid.UUID), token["aud"].(bool))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
	}
	return &userv1.RefreshTokenResponse{
		Item: &userv1.Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}
