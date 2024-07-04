package mysql

import (
	"github.com/letscrum/letscrum/internal/model"
	"gorm.io/gorm"
)

type SprintDao struct {
	DB *gorm.DB
}

func (s SprintDao) Get(sprint model.Sprint) (*model.Sprint, error) {
	var currentSprint *model.Sprint
	err := s.DB.Where("id = ?", sprint.Id).Find(&currentSprint).Error
	if err != nil {
		return nil, err
	}
	return currentSprint, nil
}

func (s SprintDao) ListByProject(project model.Project, page, size int32, keyword string) ([]*model.Sprint, error) {
	var sprints []*model.Sprint
	err := s.DB.Where("project_id = ?", project.Id).Where("name LIKE ?", "%"+keyword+"%").Limit(int(size)).Offset(int((page - 1) * size)).Order("start_date, name").Find(&sprints).Error
	if err != nil {
		return nil, err
	}
	return sprints, nil
}

func (s SprintDao) CountByProject(project model.Project, keyword string) int64 {
	count := int64(0)
	s.DB.Where("project_id = ?", project.Id).Where("name LIKE ?", "%"+keyword+"%").Model(&model.Sprint{}).Count(&count)
	return count
}

func (s SprintDao) Create(sprint model.Sprint) (*model.Sprint, error) {
	if err := s.DB.Create(&sprint).Error; err != nil {
		return nil, err
	}
	return &sprint, nil
}

func (s SprintDao) Update(sprint model.Sprint) (*model.Sprint, error) {
	if err := s.DB.Model(&model.Sprint{}).Where("id = ?", sprint.Id).Update("name", sprint.Name).Update("start_date", sprint.StartDate).Update("end_date", sprint.EndDate).Error; err != nil {
		return nil, err
	}
	return &sprint, nil
}

func (s SprintDao) Delete(sprint model.Sprint) (bool, error) {
	if err := s.DB.Where("id = ?", sprint.Id).Delete(&model.Sprint{}).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (s SprintDao) DeleteByProject(project model.Project) (bool, error) {
	if err := s.DB.Where("project_id = ?", project.Id).Delete(&model.Sprint{}).Error; err != nil {
		return false, err
	}
	return true, nil
}

func NewSprintDao(d *gorm.DB) *SprintDao {
	return &SprintDao{d}
}
