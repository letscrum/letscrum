package v1

import (
	"context"

	generalv1 "github.com/letscrum/letscrum/api/general/v1"
	letscrumv1 "github.com/letscrum/letscrum/api/letscrum/v1"
	userV1 "github.com/letscrum/letscrum/api/user/v1"
	"github.com/letscrum/letscrum/internal/dao"
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

func (s *UserService) List(ctx context.Context, req *userV1.ListUserRequest) (*userV1.ListUserResponse, error) {
	users, err := s.userDao.List(req.Page, req.Size, req.Keyword)
	if err != nil {
		result := status.Convert(err)
		if result.Code() == codes.NotFound {
			return nil, status.Errorf(codes.NotFound, "not found.")
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}
	var list []*userV1.User
	for _, u := range users {
		list = append(list, &userV1.User{
			Id:           u.Id,
			Name:         u.Name,
			Email:        u.Email,
			IsSuperAdmin: u.IsSuperAdmin,
			CreatedAt:    u.CreatedAt.Unix(),
			UpdatedAt:    u.UpdatedAt.Unix(),
		})
	}
	count := s.userDao.Count(req.Keyword)
	return &userV1.ListUserResponse{
		Items: list,
		Pagination: &generalv1.Pagination{
			Page:  req.Page,
			Size:  req.Size,
			Total: int32(count),
		},
	}, nil
}

//
//func Create2(name string, email string, password string, isSuperAdmin bool) (int64, error) {
//    id, err := model.CreateUser(name, email, password, isSuperAdmin)
//    if err != nil {
//        return 0, err
//    }
//    return id, nil
//}
//
//func List2(keyword string, page int32, pageSize int32) ([]*userV1.User, int64, error) {
//    users, err := model.ListUser(keyword, page, pageSize)
//    if err != nil {
//        return nil, 0, err
//    }
//    var list []*userV1.User
//    for _, u := range users {
//        list = append(list, &userV1.User{
//            Id:        u.Id,
//            Name:      u.Name,
//            Email:     u.Email,
//            CreatedAt: u.CreatedAt.Unix(),
//            UpdatedAt: u.UpdatedAt.Unix(),
//        })
//    }
//    count := model.CountUser(keyword)
//    return list, count, nil
//}
//
//func Update2(user *userV1.User) error {
//    if err := model.UpdateUser(user.Name, user.Email, user.Password); err != nil {
//        return err
//    }
//    return nil
//}
//
//func Delete2(name string) error {
//    if err := model.DeleteUser(name); err != nil {
//        return err
//    }
//    return nil
//}
//
//func Get2(id int64) (*userV1.User, error) {
//    u, err := model.GetUser(id)
//    if err != nil {
//        return nil, err
//    }
//    user := &userV1.User{
//        Id:        u.Id,
//        Name:      u.Name,
//        Email:     u.Email,
//        CreatedAt: u.CreatedAt.Unix(),
//        UpdatedAt: u.UpdatedAt.Unix(),
//    }
//    return user, nil
//}
//
//func SignIn2(name string, password string) (*userV1.User, error) {
//    u, err := model.GetUserWithPassword(name, password)
//    if err != nil {
//        return nil, err
//    }
//    accessToken, refreshToken, errGenTokens := utils.GenerateTokens(strconv.FormatInt(u.Id, 10))
//    if errGenTokens != nil {
//        return nil, errGenTokens
//    }
//    user := &userV1.User{
//        Id:           u.Id,
//        Name:         u.Name,
//        Email:        u.Email,
//        IsSuperAdmin: u.IsSuperAdmin,
//        CreatedAt:    u.CreatedAt.Unix(),
//        UpdatedAt:    u.UpdatedAt.Unix(),
//        Token: &userV1.Token{
//            AccessToken:  accessToken,
//            RefreshToken: refreshToken,
//        },
//    }
//    return user, nil
//}
