package mysql

import (
	"github.com/letscrum/letscrum/internal/model"
	"gorm.io/gorm"
)

type ProjectMemberDao struct {
	DB *gorm.DB
}

func (p ProjectMemberDao) Update(projectID, userID int64, isAdmin bool) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (p *ProjectMemberDao) List(projectID int64, page, size int32) ([]*model.ProjectMember, error) {
	var projectMembers []*model.ProjectMember
	err := p.DB.Where("project_id = ?", projectID).Limit(int(size)).Offset(int((page - 1) * size)).Preload("User").Find(&projectMembers).Error
	if err != nil {
		return nil, err
	}
	return projectMembers, nil
}

func (p ProjectMemberDao) Count() int64 {
	//TODO implement me
	panic("implement me")
}

func (p ProjectMemberDao) Add(projectID int64, userID int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProjectMemberDao) Remove(projectID int64, userID int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func NewProjectMemberDao(d *gorm.DB) *ProjectMemberDao {
	return &ProjectMemberDao{d}
}
