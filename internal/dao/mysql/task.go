package mysql

import (
	"github.com/letscrum/letscrum/internal/model"
	"gorm.io/gorm"
)

type TaskDao struct {
	DB *gorm.DB
}

func (t TaskDao) Get(id int64) (*model.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (t TaskDao) List(page, size int32, keyword string) ([]*model.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (t TaskDao) Count(keyword string) int64 {
	//TODO implement me
	panic("implement me")
}

func (t TaskDao) Create(project *model.Task) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (t TaskDao) Update(project *model.Task) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (t TaskDao) Delete(id int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func NewTaskDao(d *gorm.DB) *TaskDao {
	return &TaskDao{d}
}
