package user

import (
	generalV1 "github.com/letscrum/letscrum/apis/general/v1"
	userV1 "github.com/letscrum/letscrum/apis/user/v1"
	"github.com/letscrum/letscrum/models"
	"github.com/letscrum/letscrum/pkg/utils"
	"strconv"
)

func Create(user *userV1.User) error {
	if err := models.CreateUser(user); err != nil {
		return err
	}
	return nil
}

func List(pagination *generalV1.Pagination) ([]*userV1.User, int64, error) {
	users, err := models.ListUser(pagination)
	if err != nil {
		return nil, 0, err
	}
	var list []*userV1.User
	for _, u := range users {
		list = append(list, &userV1.User{
			Id:        u.Id,
			Name:      u.Name,
			Email:     u.Email,
			CreatedAt: u.CreatedAt.Unix(),
			UpdatedAt: u.UpdatedAt.Unix(),
		})
	}
	count := models.CountProject()
	return list, count, nil
}

func Update(user *userV1.User) error {
	if err := models.UpdateUser(user.Name, user); err != nil {
		return err
	}
	return nil
}

func Delete(name string) error {
	if err := models.DeleteUser(name); err != nil {
		return err
	}
	return nil
}

func Get(name string) (*userV1.User, error) {
	u, err := models.GetUser(name)
	if err != nil {
		return nil, err
	}
	user := &userV1.User{
		Id:        u.Id,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Unix(),
		UpdatedAt: u.UpdatedAt.Unix(),
	}
	return user, nil
}

func SignIn(name string, password string) (*userV1.User, error) {
	u, err := models.GetUserWithPassword(name, password)
	if err != nil {
		return nil, err
	}
	accessToken, refreshToken, errGenTokens := utils.GenerateTokens(strconv.FormatInt(u.Id, 10))
	if errGenTokens != nil {
		return nil, errGenTokens
	}
	user := &userV1.User{
		Id:        u.Id,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Unix(),
		UpdatedAt: u.UpdatedAt.Unix(),
		Token: &userV1.Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}
	return user, nil
}
