package mysql

import (
	"github.com/letscrum/letscrum/internal/model"
	"gorm.io/gorm"
)

type SprintDao struct {
	DB *gorm.DB
}

func (s SprintDao) Get(id int64) (*model.Sprint, error) {
	//TODO implement me
	panic("implement me")
}

func (s SprintDao) List(page, size int32, keyword string) ([]*model.Sprint, error) {
	//TODO implement me
	panic("implement me")
}

func (s SprintDao) Count(keyword string) int64 {
	//TODO implement me
	panic("implement me")
}

func (s SprintDao) Create(sprint *model.Sprint) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (s SprintDao) Update(sprint *model.Sprint) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (s SprintDao) Delete(id int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func NewSprintDao(d *gorm.DB) *SprintDao {
	return &SprintDao{d}
}
