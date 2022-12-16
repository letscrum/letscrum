package mysql

import (
	"github.com/letscrum/letscrum/internal/model"
	"gorm.io/gorm"
)

type SprintDao struct {
	DB *gorm.DB
}

func (s SprintDao) Get(id int64) (*model.Sprint, error) {
	var sprint *model.Sprint
	err := s.DB.Where("id = ?", id).Find(&sprint).Error
	if err != nil {
		return nil, err
	}
	return sprint, nil
}

func (s SprintDao) List(projectID int64, page, size int32, keyword string) ([]*model.Sprint, error) {
	var sprints []*model.Sprint
	err := s.DB.Where("project_id = ?", projectID).Where("name LIKE ?", "%"+keyword+"%").Limit(int(size)).Offset(int((page - 1) * size)).Order("start_date, name").Find(&sprints).Error
	if err != nil {
		return nil, err
	}
	return sprints, nil
}

func (s SprintDao) Count(projectID int64, keyword string) int64 {
	count := int64(0)
	s.DB.Where("project_id = ?", projectID).Where("name LIKE ?", "%"+keyword+"%").Model(&model.Sprint{}).Count(&count)
	return count
}

func (s SprintDao) Create(sprint *model.Sprint) (int64, error) {
	if err := s.DB.Create(&sprint).Error; err != nil {
		return 0, err
	}
	var projectMembers []*model.ProjectMember
	err := s.DB.Where("project_id = ?", sprint.ProjectID).Find(&projectMembers).Error
	if err != nil {
		return 0, err
	}
	var sprintMembers []*model.SprintMember
	for _, pm := range projectMembers {
		sprintMember := model.SprintMember{
			SprintID: sprint.ID,
			UserID:   pm.UserID,
		}
		sprintMembers = append(sprintMembers, &sprintMember)
	}
	if err := s.DB.Create(&sprintMembers).Error; err != nil {
		return 0, err
	}
	return sprint.ID, nil
}

func (s SprintDao) Update(sprint *model.Sprint) (bool, error) {
	if err := s.DB.Model(&model.Sprint{}).Where("id = ?", sprint.ID).Update("name", sprint.Name).Update("start_date", sprint.StartDate).Update("end_date", sprint.EndDate).Error; err != nil {
		return false, nil
	}
	return true, nil
}

func (s SprintDao) Delete(id int64) (bool, error) {
	if err := s.DB.Where("id = ?", id).Delete(&model.Sprint{}).Error; err != nil {
		return false, nil
	}
	return true, nil
}

func NewSprintDao(d *gorm.DB) *SprintDao {
	return &SprintDao{d}
}
