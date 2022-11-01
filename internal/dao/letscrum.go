package dao

import (
	"github.com/letscrum/letscrum/internal/model"
)

type LetscrumDao interface {
	SignIn(name, password string) (*model.User, error)
}
