package dao

import (
	"context"
	"github.com/letscrum/letscrum/internal/model"
)

type ProjectDao interface {
	Get(ctx context.Context, id int64) (*model.Project, error)
	List(ctx context.Context, page, size int32) ([]*model.Project, error)
	Count(context.Context) int64
	Create(ctx context.Context, project *model.Project) (bool, error)
	Update(ctx context.Context, project *model.Project) (bool, error)
	Delete(ctx context.Context, id int64) (bool, error)
}
