package dao

import (
	"context"
	"github.com/letscrum/letscrum/internal/model"
)

type ProjectDao interface {
	Get(ctx context.Context, id int64) (*model.Project, error)
	//List(ctx context.Context, page, pageSize int32) (*[]model.Project, error)
	//Create(ctx context.Context, project *model.Project) (bool, error)
	//Update(ctx context.Context, project *model.Project) (bool, error)
	//Delete(ctx context.Context, id int64) (bool, error)
}
