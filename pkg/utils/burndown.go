package utils

import (
	"time"

	"github.com/google/uuid"
	"github.com/letscrum/letscrum/internal/model"
	"gorm.io/gorm"
)

func UpdateBurndown(tx *gorm.DB, sprintId uuid.UUID, correctDate time.Time, taskCount int32, hours float32) error {
	// update sprint status record and set work item count + 1
	if err := tx.Model(&model.SprintBurndown{}).Where("sprint_id = ?", sprintId).Where("sprint_date = ?", correctDate).Update("task_count", gorm.Expr("task_count + ?", taskCount)).Update("work_hours", gorm.Expr("work_hours + ?", hours)).Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func GetBurndown(tx *gorm.DB, sprintId uuid.UUID) ([]*model.SprintBurndown, int, error) {
	// get current sprint statuses ordered by date
	var currentSprintBurndown []*model.SprintBurndown
	if err := tx.Where("sprint_id = ?", sprintId).Order("sprint_date").Find(&currentSprintBurndown).Error; err != nil {
		tx.Rollback()
		return nil, 0, err
	}
	lastSprintStatusIndex := len(currentSprintBurndown) - 1
	return currentSprintBurndown, lastSprintStatusIndex, nil
}
