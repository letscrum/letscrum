package dao

import (
	"context"
	"github.com/letscrum/letscrum/internal/model"
)

type LetscrumDao interface {
	GetVersion(ctx context.Context) (*model.Project, error)
	SignIn(ctx context.Context, name, password string) (*model.User, error)
}
