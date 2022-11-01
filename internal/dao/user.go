package dao

import "github.com/letscrum/letscrum/internal/model"

type UserDao interface {
	Get(id int64) (*model.User, error)
}
