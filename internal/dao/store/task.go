package store

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/letscrum/letscrum/internal/model"
	"github.com/letscrum/letscrum/pkg/utils"
	"gorm.io/gorm"
)

type TaskDao struct {
	DB *gorm.DB
}

func (t TaskDao) ListByWorkItemIds(workItemIds []int64) ([]*model.Task, error) {
	var tasks []*model.Task
	err := t.DB.Where("work_item_id IN ?", workItemIds).Preload("CreatedUser").Preload("AssignUser").Order("`order`").Order("created_at").Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (t TaskDao) Get(task model.Task) (*model.Task, []*model.ItemLog, error) {
	if err := t.DB.Where("id = ?", task.Id).Find(&task).Error; err != nil {
		return nil, nil, err
	}
	var logs []*model.ItemLog
	err := t.DB.Where("item_id = ?", task.Id).Where("item_type = ?", "TASK").Find(&logs).Error
	if err != nil {
		return nil, nil, err
	}
	return &task, logs, nil
}

func (t TaskDao) List(page, size int32, keyword string) ([]*model.Task, error) {
	var tasks []*model.Task
	err := t.DB.Where("title LIKE ?", "%"+keyword+"%").Limit(int(size)).Offset(int((page - 1) * size)).Preload("CreatedUser").Preload("AssignUser").Order("`order`").Order("created_at").Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (t TaskDao) ListByWorkItem(workItemId int64, page, size int32, keyword string) ([]*model.Task, error) {
	var tasks []*model.Task
	err := t.DB.Where("work_item_id = ?", workItemId).Where("title LIKE ?", "%"+keyword+"%").Limit(int(size)).Offset(int((page - 1) * size)).Preload("CreatedUser").Preload("AssignUser").Order("`order`").Order("created_at").Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (t TaskDao) ListByProject(projectId uuid.UUID, page, size int32, keyword string) ([]*model.Task, error) {
	var tasks []*model.Task
	err := t.DB.Where("project_id = ?", projectId).Where("title LIKE ?", "%"+keyword+"%").Limit(int(size)).Offset(int((page - 1) * size)).Preload("CreatedUser").Preload("AssignUser").Order("`order`").Order("created_at").Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (t TaskDao) CountByProject(projectId uuid.UUID, keyword string) int64 {
	count := int64(0)
	t.DB.Where("project_id = ?", projectId).Where("title LIKE ?", "%"+keyword+"%").Model(&model.Task{}).Count(&count)
	return count
}

func (t TaskDao) ListByProjectNotInSprint(projectId uuid.UUID, page, size int32, keyword string) ([]*model.Task, error) {
	var tasks []*model.Task
	err := t.DB.Where("project_id = ?", projectId).Where("sprint_id IS NULL").Where("title LIKE ?", "%"+keyword+"%").Limit(int(size)).Offset(int((page - 1) * size)).Preload("CreatedUser").Preload("AssignUser").Order("`order`").Order("created_at").Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (t TaskDao) CountByProjectNotInSprint(projectId uuid.UUID, keyword string) int64 {
	count := int64(0)
	t.DB.Where("project_id = ?", projectId).Where("sprint_id IS NULL").Where("title LIKE ?", "%"+keyword+"%").Model(&model.Task{}).Count(&count)
	return count
}

func (t TaskDao) Count(keyword string) int64 {
	count := int64(0)
	t.DB.Where("title LIKE ?", "%"+keyword+"%").Model(&model.Task{}).Count(&count)
	return count
}

func (t TaskDao) CountByWorkItem(workItemId int64, keyword string) int64 {
	count := int64(0)
	t.DB.Where("work_item_id = ?", workItemId).Where("title LIKE ?", "%"+keyword+"%").Model(&model.Task{}).Count(&count)
	return count
}

func (t TaskDao) CountBySprint(sprintId uuid.UUID, keyword string) int64 {
	count := int64(0)
	t.DB.Where("sprint_id = ?", sprintId).Where("title LIKE ?", "%"+keyword+"%").Model(&model.Task{}).Count(&count)
	return count
}

func (t TaskDao) CountBySprints(sprintIds []uuid.UUID) []int64 {
	var counts []int64
	for _, sprintId := range sprintIds {
		count := int64(0)
		t.DB.Where("sprint_id = ?", sprintId).Model(&model.Task{}).Count(&count)
		counts = append(counts, count)
	}
	return counts
}

func (t TaskDao) WorkHoursBySprint(sprintId uuid.UUID) float32 {
	var workHours float32
	t.DB.Model(&model.Task{}).Where("sprint_id = ?", sprintId).Select("COALESCE(SUM(remaining), 0)").Scan(&workHours)
	return workHours
}

func (t TaskDao) Create(task model.Task) (*model.Task, error) {
	// make transaction of create task and task log
	err := t.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&task).Error; err != nil {
			tx.Rollback()
			return err
		}
		var log model.ItemLog
		log.Id = uuid.New()
		log.ItemId = task.Id
		log.ItemType = "TASK"
		log.Action = "CREATE"
		log.Log = "Create task"
		log.Changes = "projectId: " + task.ProjectId.String() + ", sprintId: " + task.SprintId.String()
		log.CreatedBy = task.CreatedBy
		if err := tx.Create(&log).Error; err != nil {
			tx.Rollback()
			return err
		}

		if task.SprintId != uuid.Nil {
			// get now time
			createdDate := time.Now()
			currentSprintStatuses, lastSprintStatusIndex, err := utils.GetBurndown(tx, task.SprintId)
			if err != nil {
				return err
			}

			// if createdDate before or equal the last sprint status date
			if createdDate.Before(currentSprintStatuses[lastSprintStatusIndex].SprintDate) || createdDate.Equal(currentSprintStatuses[lastSprintStatusIndex].SprintDate) {
				// if the first sprint status date is after createdDate, set createdDate to the first sprint status date
				if currentSprintStatuses[0].SprintDate.After(createdDate) {
					createdDate = currentSprintStatuses[0].SprintDate
				}
				// format createdDate to date
				correctDate := time.Date(createdDate.Year(), createdDate.Month(), createdDate.Day(), 0, 0, 0, 0, createdDate.Location())

				err := utils.UpdateBurndown(tx, task.SprintId, correctDate, 1, task.Remaining)

				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (t TaskDao) Update(task model.Task, userId uuid.UUID) (*model.Task, error) {
	// make transaction of update task and task log
	err := t.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Task{}).Where("id = ?", task.Id).Updates(task).Error; err != nil {
			tx.Rollback()
			return err
		}
		var log model.ItemLog
		log.Id = uuid.New()
		log.ItemId = task.Id
		log.ItemType = "TASK"
		log.Action = "UPDATE"
		log.Log = "Update task"
		log.CreatedBy = userId
		if err := tx.Create(&log).Error; err != nil {
			tx.Rollback()
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (t TaskDao) UpdateAssignUser(task model.Task, userId uuid.UUID) (*model.Task, error) {
	// make transaction of update task assign user and task log
	err := t.DB.Transaction(func(tx *gorm.DB) error {
		if task.AssignTo == uuid.Nil {
			if err := tx.Model(&model.Task{}).Where("id = ?", task.Id).Update("assign_to", sql.NullString{}).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			if err := tx.Model(&model.Task{}).Where("id = ?", task.Id).Update("assign_to", task.AssignTo).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
		var log model.ItemLog
		log.Id = uuid.New()
		log.ItemId = task.Id
		log.ItemType = "TASK"
		log.Action = "UPDATE"
		log.Log = "Update task, assign task to the user, id is: " + task.AssignTo.String()
		log.CreatedBy = userId
		if err := tx.Create(&log).Error; err != nil {
			tx.Rollback()
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (t TaskDao) Move(task model.Task, userId uuid.UUID) (*model.Task, error) {
	// make transaction of update task status and move task to work item
	err := t.DB.Transaction(func(tx *gorm.DB) error {
		originalTask := model.Task{}
		if err := tx.Where("id = ?", task.Id).Find(&originalTask).Error; err != nil {
			tx.Rollback()
			return err
		}
		var workItemId int64
		if task.WorkItemId == 0 {
			workItemId = originalTask.WorkItemId
		}
		workItemId = task.WorkItemId
		var remaining float32
		if task.Status == "Done" || task.Status == "Removed" {
			remaining = 0
		} else {
			remaining = originalTask.Remaining
		}
		if err := tx.Model(&model.Task{}).Where("id = ?", task.Id).Find(&originalTask).Update("status", task.Status).Update("work_item_id", workItemId).Update("remaining", remaining).Error; err != nil {
			tx.Rollback()
			return err
		}
		// get the updated task
		if err := tx.Where("id = ?", task.Id).Find(&task).Error; err != nil {
			tx.Rollback()
			return err
		}
		var log model.ItemLog
		log.Id = uuid.New()
		log.ItemId = task.Id
		log.ItemType = "TASK"
		log.Action = "UPDATE"
		log.Log = "Update task status and move task to work item"
		log.Changes = "projectId: " + task.ProjectId.String() + ", sprintId: " + task.SprintId.String()
		log.CreatedBy = userId
		if err := tx.Create(&log).Error; err != nil {
			tx.Rollback()
			return err
		}
		// get now time
		createdDate := time.Now()
		currentSprintStatuses, lastSprintStatusIndex, err := utils.GetBurndown(tx, originalTask.SprintId)
		if err != nil {
			return err
		}

		// if createdDate before or equal the last sprint status date
		if createdDate.Before(currentSprintStatuses[lastSprintStatusIndex].SprintDate) || createdDate.Equal(currentSprintStatuses[lastSprintStatusIndex].SprintDate) {
			// if the first sprint status date is after createdDate, set createdDate to the first sprint status date
			if currentSprintStatuses[0].SprintDate.After(createdDate) {
				createdDate = currentSprintStatuses[0].SprintDate
			}
			// format createdDate to date
			correctDate := time.Date(createdDate.Year(), createdDate.Month(), createdDate.Day(), 0, 0, 0, 0, createdDate.Location())

			var taskCount int32
			var workHours float32
			if originalTask.Status == "Done" || originalTask.Status == "Removed" {
				// check task status is Done or Removed
				if task.Status == "Done" || task.Status == "Removed" {
					taskCount = 0
					workHours = 0
				} else {
					taskCount = 1
					workHours = 0
				}
			} else {
				if task.Status == "Done" || task.Status == "Removed" {
					taskCount = -1
					workHours = -originalTask.Remaining
				} else {
					taskCount = 0
					workHours = 0
				}
			}

			err := utils.UpdateBurndown(tx, originalTask.SprintId, correctDate, taskCount, workHours)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (t TaskDao) UpdateWorkHours(task model.Task, userId uuid.UUID) (*model.Task, error) {
	// make transaction of update task work hours and task log
	err := t.DB.Transaction(func(tx *gorm.DB) error {
		originalTask := model.Task{}
		if err := tx.Where("id = ?", task.Id).Find(&originalTask).Error; err != nil {
			tx.Rollback()
			return err
		}
		var remaining float32
		if task.Status == "Done" || task.Status == "Removed" {
			remaining = 0
		} else {
			remaining = task.Remaining
		}
		if err := tx.Model(&model.Task{}).Where("id = ?", task.Id).Update("remaining", remaining).Error; err != nil {
			tx.Rollback()
			return err
		}
		var log model.ItemLog
		log.Id = uuid.New()
		log.ItemId = task.Id
		log.ItemType = "TASK"
		log.Action = "UPDATE"
		log.Log = "Update task work hours"
		log.Changes = "projectId: " + task.ProjectId.String() + ", sprintId: " + task.SprintId.String()
		log.CreatedBy = userId
		if err := tx.Create(&log).Error; err != nil {
			tx.Rollback()
			return err
		}
		if task.Status != "Done" && task.Status != "Removed" {
			// get now time
			createdDate := time.Now()
			currentSprintStatuses, lastSprintStatusIndex, err := utils.GetBurndown(tx, originalTask.SprintId)
			if err != nil {
				return err
			}

			// if createdDate before or equal the last sprint status date
			if createdDate.Before(currentSprintStatuses[lastSprintStatusIndex].SprintDate) || createdDate.Equal(currentSprintStatuses[lastSprintStatusIndex].SprintDate) {
				// if the first sprint status date is after createdDate, set createdDate to the first sprint status date
				if currentSprintStatuses[0].SprintDate.After(createdDate) {
					createdDate = currentSprintStatuses[0].SprintDate
				}
				// format createdDate to date
				correctDate := time.Date(createdDate.Year(), createdDate.Month(), createdDate.Day(), 0, 0, 0, 0, createdDate.Location())
				updateRemaining := task.Remaining - originalTask.Remaining
				err := utils.UpdateBurndown(tx, originalTask.SprintId, correctDate, 0, updateRemaining)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (t TaskDao) Delete(task model.Task, userId uuid.UUID) (bool, error) {
	// make transaction of delete task and task log
	err := t.DB.Transaction(func(tx *gorm.DB) error {
		// get task and delete from database
		if err := tx.Where("id = ?", task.Id).Find(&task).Delete(&task).Error; err != nil {
			tx.Rollback()
			return err
		}

		var log model.ItemLog
		log.Id = uuid.New()
		log.ItemId = task.Id
		log.ItemType = "TASK"
		log.Action = "DELETE"
		log.Log = "Delete task"
		log.Changes = "projectId: " + task.ProjectId.String() + ", sprintId: " + task.SprintId.String()
		log.CreatedBy = userId
		if err := tx.Create(&log).Error; err != nil {
			tx.Rollback()
			return err
		}
		if task.Status != "Done" && task.Status != "Removed" {
			// get now time
			createdDate := time.Now()
			currentSprintStatuses, lastSprintStatusIndex, err := utils.GetBurndown(tx, task.SprintId)
			if err != nil {
				return err
			}

			// if createdDate before or equal the last sprint status date
			if createdDate.Before(currentSprintStatuses[lastSprintStatusIndex].SprintDate) || createdDate.Equal(currentSprintStatuses[lastSprintStatusIndex].SprintDate) {
				// if the first sprint status date is after createdDate, set createdDate to the first sprint status date
				if currentSprintStatuses[0].SprintDate.After(createdDate) {
					createdDate = currentSprintStatuses[0].SprintDate
				}
				// format createdDate to date
				correctDate := time.Date(createdDate.Year(), createdDate.Month(), createdDate.Day(), 0, 0, 0, 0, createdDate.Location())

				err := utils.UpdateBurndown(tx, task.SprintId, correctDate, -1, -task.Remaining)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (t TaskDao) ReOrder(taskIds []int64) ([]int64, error) {
	// update task order
	err := t.DB.Transaction(func(tx *gorm.DB) error {
		for i, id := range taskIds {
			if err := tx.Model(&model.Task{}).Where("id = ?", id).Update("order", i+1).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return taskIds, nil
}

func NewTaskDao(d *gorm.DB) *TaskDao {
	return &TaskDao{d}
}
