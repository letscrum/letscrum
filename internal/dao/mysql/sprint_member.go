package mysql

import (
	"github.com/letscrum/letscrum/internal/model"
	"gorm.io/gorm"
)

type SprintMemberDao struct {
	DB *gorm.DB
}

func (s SprintMemberDao) List(sprintID int64, page, size int32) ([]*model.SprintMember, error) {
	//TODO implement me
	panic("implement me")
}

func (s SprintMemberDao) Count() int64 {
	//TODO implement me
	panic("implement me")
}

func (s SprintMemberDao) Add(sprintID int64, userID int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (s SprintMemberDao) Update(sprintID, userID int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (s SprintMemberDao) Remove(sprintID, userID int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func NewSprintMemberDao(d *gorm.DB) *SprintMemberDao {
	return &SprintMemberDao{d}
}
