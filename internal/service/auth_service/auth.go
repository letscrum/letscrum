package auth_service

import (
	"github.com/letscrum/letscrum/internal/model"
)

type Auth struct {
	Username string
	Password string
}

func (a *Auth) Check() (bool, error) {
	return model.CheckAuth(a.Username, a.Password)
}
