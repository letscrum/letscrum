package dao

import (
	"github.com/google/uuid"
	"github.com/letscrum/letscrum/internal/model"
)

type UserDao interface {
	SignIn(name, password string) (*model.User, error)
	Get(user model.User) (*model.User, error)
	GetByName(user model.User) (*model.User, error)
	List(page, size int32, keyword string) ([]*model.User, error)
	ListSuperAdmins() ([]*model.User, error)
	ListByIds(page, size int32, ids []uuid.UUID) ([]*model.User, error)
	Count(keyword string) int64
	Create(user model.User) (*model.User, error)
	SetSuperAdmin(user model.User, isSuperAdmin bool) (*model.User, error)
	UpdatePassword(user model.User, newPassword string) (*model.User, error)
	ResetPassword(user model.User, newPassword string) (*model.User, error)
}
