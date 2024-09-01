package store

import (
	"time"

	"github.com/letscrum/letscrum/internal/model"
	"gorm.io/gorm"
)

type SprintDao struct {
	DB *gorm.DB
}

func (s SprintDao) Get(sprint model.Sprint) (*model.Sprint, error) {
	var currentSprint *model.Sprint
	err := s.DB.Where("project_id", sprint.ProjectId).Where("id = ?", sprint.Id).Find(&currentSprint).Error
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
	// make transaction of create sprint and sprint status
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&sprint).Error; err != nil {
			return err
		}
		var sprintStatuses []*model.SprintStatus
		for d := sprint.StartDate; d.Before(sprint.EndDate) || d.Equal(sprint.EndDate); d = d.AddDate(0, 0, 1) {
			sprintStatus := model.SprintStatus{
				SprintId:   sprint.Id,
				SprintDate: d,
			}
			sprintStatuses = append(sprintStatuses, &sprintStatus)
		}
		if err := tx.Create(&sprintStatuses).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &sprint, nil
}

func (s SprintDao) Update(sprint model.Sprint) (*model.Sprint, error) {
	// make transaction of update sprint and sprint status
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Sprint{}).Where("id = ?", sprint.Id).Update("name", sprint.Name).Update("start_date", sprint.StartDate).Update("end_date", sprint.EndDate).Error; err != nil {
			tx.Rollback()
			return err
		}
		// get current sprint statuses ordered by date
		var currentSprintStatuses []*model.SprintStatus
		if err := tx.Where("sprint_id = ?", sprint.Id).Order("sprint_date").Find(&currentSprintStatuses).Error; err != nil {
			tx.Rollback()
			return err
		}
		lastSprintStatusIndex := len(currentSprintStatuses) - 1
		// convert sprint end date to unix date
		sprintEndDate := sprint.EndDate.Unix() + 1
		// convert unix date to time
		sprint.EndDate = time.Unix(sprintEndDate, 0).AddDate(0, 0, -1)

		// if current sprint statuses is empty, create new sprint statuses
		if len(currentSprintStatuses) == 0 {
			var sprintStatuses []*model.SprintStatus
			for d := sprint.StartDate; d.Before(sprint.EndDate) || d.Equal(sprint.EndDate); d = d.AddDate(0, 0, 1) {
				sprintStatus := model.SprintStatus{
					SprintId:   sprint.Id,
					SprintDate: d,
				}
				sprintStatuses = append(sprintStatuses, &sprintStatus)
			}
			if err := tx.Create(&sprintStatuses).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
		if currentSprintStatuses[0].SprintDate.Equal(sprint.StartDate) && currentSprintStatuses[lastSprintStatusIndex].SprintDate.Equal(sprint.EndDate) {
			return nil
		} else {
			// if updated sprint start date is before current sprint status
			if sprint.StartDate.Before(currentSprintStatuses[0].SprintDate) {
				// create new sprint status from updated sprint start date to current sprint status 0
				var beforeSprintStatuses []*model.SprintStatus
				for d := sprint.StartDate; d.Before(currentSprintStatuses[0].SprintDate); d = d.AddDate(0, 0, 1) {
					sprintStatus := model.SprintStatus{
						SprintId:   sprint.Id,
						SprintDate: d,
					}
					beforeSprintStatuses = append(beforeSprintStatuses, &sprintStatus)
				}
				if err := tx.Create(&beforeSprintStatuses).Error; err != nil {
					tx.Rollback()
					return err
				}
			} else {
				// delete sprint status from current sprint status 0 to updated sprint start date
				if err := tx.Where("sprint_id = ?", sprint.Id).Where("sprint_date < ?", sprint.StartDate).Delete(&model.SprintStatus{}).Error; err != nil {
					tx.Rollback()
					return err
				}
			}
			// if updated sprint end date is after the last current sprint status
			if sprint.EndDate.After(currentSprintStatuses[len(currentSprintStatuses)-1].SprintDate) {
				// create new sprint status from the last current sprint status to updated sprint end date
				var afterSprintStatuses []*model.SprintStatus
				for d := currentSprintStatuses[len(currentSprintStatuses)-1].SprintDate.AddDate(0, 0, 1); d.Before(sprint.EndDate) || d.Equal(sprint.EndDate); d = d.AddDate(0, 0, 1) {
					sprintStatus := model.SprintStatus{
						SprintId:   sprint.Id,
						SprintDate: d,
					}
					afterSprintStatuses = append(afterSprintStatuses, &sprintStatus)
				}
				if err := tx.Create(&afterSprintStatuses).Error; err != nil {
					tx.Rollback()
					return err
				}
			} else {
				// delete sprint status from updated sprint end date to
				if err := tx.Where("sprint_id = ?", sprint.Id).Where("sprint_date > ?", sprint.EndDate).Delete(&model.SprintStatus{}).Error; err != nil {
					tx.Rollback()
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &sprint, nil
}

func (s SprintDao) UpdateMembers(sprint model.Sprint) (*model.Sprint, error) {
	if err := s.DB.Model(&model.Sprint{}).Where("id = ?", sprint.Id).Update("members", sprint.Members).Error; err != nil {
		return nil, err
	}
	return &sprint, nil
}

func (s SprintDao) Delete(sprint model.Sprint) (bool, error) {
	// make transaction of delete sprint and sprint status
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", sprint.Id).Delete(&model.Sprint{}).Error; err != nil {
			return err
		}
		if err := tx.Where("sprint_id = ?", sprint.Id).Delete(&model.SprintStatus{}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
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
