package store

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/letscrum/letscrum/internal/model"
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
		// get now time
		createdDate := time.Now()
		// get current sprint statuses ordered by date
		var currentSprintStatuses []*model.SprintStatus
		if err := tx.Where("sprint_id = ?", task.SprintId).Order("sprint_date").Find(&currentSprintStatuses).Error; err != nil {
			tx.Rollback()
			return err
		}
		lastSprintStatusIndex := len(currentSprintStatuses) - 1

		// if createdDate before or equal the last sprint status date
		if createdDate.Before(currentSprintStatuses[lastSprintStatusIndex].SprintDate) || createdDate.Equal(currentSprintStatuses[lastSprintStatusIndex].SprintDate) {
			// if the first sprint status date is after createdDate, set createdDate to the first sprint status date
			if currentSprintStatuses[0].SprintDate.After(createdDate) {
				createdDate = currentSprintStatuses[0].SprintDate
			}
			// format createdDate to date
			correctDate := time.Date(createdDate.Year(), createdDate.Month(), createdDate.Day(), 0, 0, 0, 0, createdDate.Location())

			// update sprint status record and set work item count + 1
			if err := tx.Model(&model.SprintStatus{}).Where("sprint_id = ?", task.SprintId).Where("sprint_date = ?", correctDate).Update("task_count", gorm.Expr("task_count + ?", 1)).Error; err != nil {
				tx.Rollback()
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

func (t TaskDao) UpdateStatus(task model.Task, userId uuid.UUID) (*model.Task, error) {
	// make transaction of update task status and task log
	err := t.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Task{}).Where("id = ?", task.Id).Update("status", task.Status).Error; err != nil {
			tx.Rollback()
			return err
		}
		var log model.ItemLog
		log.Id = uuid.New()
		log.ItemId = task.Id
		log.ItemType = "TASK"
		log.Action = "UPDATE"
		log.Log = "Update task status to: " + task.Status
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
		if err := tx.Model(&model.Task{}).Where("id = ?", task.Id).Update("status", task.Status).Update("work_item_id", task.WorkItemId).Error; err != nil {
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
		if err := tx.Where("id = ?", task.Id).Delete(&model.Task{}).Error; err != nil {
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
		// get now time
		createdDate := time.Now()
		// get current sprint statuses ordered by date
		var currentSprintStatuses []*model.SprintStatus
		if err := tx.Where("sprint_id = ?", task.SprintId).Order("sprint_date").Find(&currentSprintStatuses).Error; err != nil {
			tx.Rollback()
			return err
		}
		lastSprintStatusIndex := len(currentSprintStatuses) - 1

		// if createdDate before or equal the last sprint status date
		if createdDate.Before(currentSprintStatuses[lastSprintStatusIndex].SprintDate) || createdDate.Equal(currentSprintStatuses[lastSprintStatusIndex].SprintDate) {
			// if the first sprint status date is after createdDate, set createdDate to the first sprint status date
			if currentSprintStatuses[0].SprintDate.After(createdDate) {
				createdDate = currentSprintStatuses[0].SprintDate
			}
			// format createdDate to date
			correctDate := time.Date(createdDate.Year(), createdDate.Month(), createdDate.Day(), 0, 0, 0, 0, createdDate.Location())

			// update sprint status record and set work item count + 1
			if err := tx.Model(&model.SprintStatus{}).Where("sprint_id = ?", task.SprintId).Where("sprint_date = ?", correctDate).Update("task_count", gorm.Expr("task_count - ?", 1)).Error; err != nil {
				tx.Rollback()
				return err
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
