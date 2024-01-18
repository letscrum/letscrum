package mysql

import (
	"github.com/letscrum/letscrum/internal/model"
	"gorm.io/gorm"
)

type SprintMemberDao struct {
	DB *gorm.DB
}

func (s SprintMemberDao) ListBySprint(sprint model.Sprint, page, size int32) ([]*model.SprintMember, error) {
	var sprintMembers []*model.SprintMember
	err := s.DB.Where("sprint_id = ?", sprint.ID).Limit(int(size)).Offset(int((page - 1) * size)).Preload("MemberUser").Find(&sprintMembers).Error
	if err != nil {
		return nil, err
	}
	return sprintMembers, nil
}

func (s SprintMemberDao) CountBySprint(sprint model.Sprint) int64 {
	count := int64(0)
	s.DB.Where("sprint_id = ?", sprint.ID).Model(&model.SprintMember{}).Count(&count)
	return count
}

func (s SprintMemberDao) BatchAdd(sprintMembers []model.SprintMember) ([]*model.SprintMember, error) {
	var sMembers []*model.SprintMember
	for _, sm := range sprintMembers {
		sprintMember := model.SprintMember{
			SprintID: sm.SprintID,
			UserID:   sm.UserID,
			Role:     sm.Role,
			Capacity: sm.Capacity,
		}
		sMembers = append(sMembers, &sprintMember)
	}
	if len(sMembers) > 0 {
		if err := s.DB.Create(&sMembers).Error; err != nil {
			return nil, err
		}
	}
	return sMembers, nil
}

func (s SprintMemberDao) BatchUpdate(sprintMembers []model.SprintMember) ([]*model.SprintMember, error) {
	var sMembers []*model.SprintMember
	for _, m := range sprintMembers {
		if err := s.DB.Model(&m).Updates(model.SprintMember{
			Role:     m.Role,
			Capacity: m.Capacity,
		}).Error; err != nil {
			return nil, err
		}
		sMembers = append(sMembers, &m)
	}
	return sMembers, nil
}

func (s SprintMemberDao) BatchRemove(sprintMembers []model.SprintMember) ([]*model.SprintMember, error) {
	//TODO implement me
	panic("implement me")
}

func NewSprintMemberDao(d *gorm.DB) *SprintMemberDao {
	return &SprintMemberDao{d}
}
