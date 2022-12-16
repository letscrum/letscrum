package mysql

import (
	"github.com/letscrum/letscrum/internal/model"
	"gorm.io/gorm"
)

type SprintMemberDao struct {
	DB *gorm.DB
}

func (s SprintMemberDao) List(sprintID int64, page, size int32) ([]*model.SprintMember, error) {
	var sprintMembers []*model.SprintMember
	err := s.DB.Where("sprint_id = ?", sprintID).Limit(int(size)).Offset(int((page - 1) * size)).Preload("User").Find(&sprintMembers).Error
	if err != nil {
		return nil, err
	}
	return sprintMembers, nil
}

func (s SprintMemberDao) Count(sprintID int64) int64 {
	count := int64(0)
	s.DB.Where("sprint_id = ?", sprintID).Model(&model.SprintMember{}).Count(&count)
	return count
}

func (s SprintMemberDao) Add(sprintID int64, userID int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (s SprintMemberDao) Update(sprintMembers []*model.SprintMember) (bool, error) {
	for _, m := range sprintMembers {
		if err := s.DB.Model(&m).Updates(model.SprintMember{
			Role:     m.Role,
			Capacity: m.Capacity,
		}).Error; err != nil {
			return false, nil
		}
	}
	return true, nil
}

func (s SprintMemberDao) Remove(sprintID, userID int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func NewSprintMemberDao(d *gorm.DB) *SprintMemberDao {
	return &SprintMemberDao{d}
}
