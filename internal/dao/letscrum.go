package dao

import (
	"context"
	"github.com/letscrum/letscrum/internal/model"
)

type LetscrumDao interface {
	SignIn(ctx context.Context, name, password string) (*model.User, error)
}
