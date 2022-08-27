package dao

import (
	"context"
	"github.com/letscrum/letscrum/internal/model"
)

type DemoDbDao interface {
	DemoDb(ctx context.Context, demoDb string) (*model.DemoDb, error)
}
