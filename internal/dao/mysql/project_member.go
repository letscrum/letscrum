package mysql

import (
	"github.com/letscrum/letscrum/internal/model"
	"gorm.io/gorm"
)

type ProjectMemberDao struct {
	Db *gorm.DB
}

func (p ProjectMemberDao) Update(projectId, userId int64, isAdmin bool) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProjectMemberDao) List(projectId int64, page, size int32) ([]*model.ProjectMember, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProjectMemberDao) Count() int64 {
	//TODO implement me
	panic("implement me")
}

func (p ProjectMemberDao) Add(projectId int64, userId int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProjectMemberDao) Remove(projectId int64, userId int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func NewProjectMemberDao(d *gorm.DB) *ProjectMemberDao {
	return &ProjectMemberDao{d}
}
