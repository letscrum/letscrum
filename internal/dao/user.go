package dao

import "github.com/letscrum/letscrum/internal/model"

type UserDao interface {
	SignIn(name, password string) (*model.User, error)
	Get(id int64) (*model.User, error)
}
