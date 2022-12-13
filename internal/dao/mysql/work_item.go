package mysql

import (
	"github.com/letscrum/letscrum/internal/model"
	"gorm.io/gorm"
)

type WorkItemDao struct {
	DB *gorm.DB
}

func (w WorkItemDao) Get(id int64) (*model.WorkItem, error) {
	//TODO implement me
	panic("implement me")
}

func (w WorkItemDao) List(page, size int32, keyword string) ([]*model.WorkItem, error) {
	//TODO implement me
	panic("implement me")
}

func (w WorkItemDao) Count(keyword string) int64 {
	//TODO implement me
	panic("implement me")
}

func (w WorkItemDao) Create(project *model.WorkItem) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (w WorkItemDao) Update(project *model.WorkItem) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (w WorkItemDao) Delete(id int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func NewWorkItemDao(d *gorm.DB) *WorkItemDao {
	return &WorkItemDao{d}
}
