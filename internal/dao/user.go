package dao

import "github.com/letscrum/letscrum/internal/model"

type UserDao interface {
	SignIn(name, password string) (*model.User, error)
	Get(id int64) (*model.User, error)
	List(page, size int32, keyword string) ([]*model.User, error)
	ListSuperAdmins() ([]*model.User, error)
	ListByIds(page, size int32, ids []int64) ([]*model.User, error)
	Count(keyword string) int64
	Create(name, email, password string, isSuperAdmin bool) (*model.User, error)
	SetSuperAdmin(id int64, isSuperAdmin bool) (*model.User, error)
	UpdatePassword(id int64, oldPassword, newPassword string) (*model.User, error)
	ResetPassword(id int64, newPassword string) (*model.User, error)
}
