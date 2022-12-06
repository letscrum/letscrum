package dao

import "github.com/letscrum/letscrum/internal/model"

type SprintDao interface {
	Get(id int64) (*model.Sprint, error)
	List(page, size int32, keyword string) ([]*model.Sprint, error)
	Count(keyword string) int64
	Create(sprint *model.Sprint) (int64, error)
	Update(sprint *model.Sprint) (bool, error)
	Delete(id int64) (bool, error)
}
