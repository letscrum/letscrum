package store

import (
	"github.com/letscrum/letscrum/internal/model"
	"gorm.io/gorm"
)

type OrgDao struct {
	DB *gorm.DB
}

func (d OrgDao) Get(org model.Org) (*model.Org, error) {
	//TODO implement me
	panic("implement me")
}

func (d OrgDao) List(page, size int32, keyword string) ([]*model.Org, error) {
	//TODO implement me
	panic("implement me")
}

func (d OrgDao) Count(keyword string) int64 {
	//TODO implement me
	panic("implement me")
}

func (d OrgDao) Create(org model.Org) (*model.Org, error) {
	//TODO implement me
	panic("implement me")
}

func (d OrgDao) Update(org model.Org) (*model.Org, error) {
	//TODO implement me
	panic("implement me")
}

func (d OrgDao) Delete(org model.Org) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func NewOrgDao(d *gorm.DB) *OrgDao {
	return &OrgDao{d}
}
