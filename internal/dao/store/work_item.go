package store

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/letscrum/letscrum/internal/model"
	"gorm.io/gorm"
)

type WorkItemDao struct {
	DB *gorm.DB
}

func (w WorkItemDao) Get(workItem model.WorkItem) (*model.WorkItem, []*model.ItemLog, error) {
	// get workitem from database
	if err := w.DB.Where("id = ?", workItem.Id).Find(&workItem).Error; err != nil {
		return nil, nil, err
	}
	// get logs from database
	var logs []*model.ItemLog
	err := w.DB.Where("item_id = ?", workItem.Id).Where("item_type = ?", "WORKITEM").Find(&logs).Error
	if err != nil {
		return nil, nil, err
	}
	return &workItem, logs, nil
}

func (w WorkItemDao) ListByProject(projectId uuid.UUID, page, size int32, keyword string) ([]*model.WorkItem, error) {
	// get workitems from database
	var workItems []*model.WorkItem
	err := w.DB.Where("project_id = ?", projectId).Where("sprint_id IS NULL").Where("title LIKE ?", "%"+keyword+"%").Limit(int(size)).Offset(int((page - 1) * size)).Preload("CreatedUser").Preload("AssignUser").Order("`order`").Order("created_at").Find(&workItems).Error
	if err != nil {
		return nil, err
	}
	return workItems, nil
}

func (w WorkItemDao) ListByProjectNotInSprint(projectId uuid.UUID, page, size int32, keyword string) ([]*model.WorkItem, error) {
	// get workitems from database
	var workItems []*model.WorkItem
	err := w.DB.Where("project_id = ?", projectId).Where("sprint_id IS NULL").Where("title LIKE ?", "%"+keyword+"%").Limit(int(size)).Offset(int((page - 1) * size)).Preload("CreatedUser").Preload("AssignUser").Order("`order`").Order("created_at").Find(&workItems).Error
	if err != nil {
		return nil, err
	}
	return workItems, nil
}

func (w WorkItemDao) ListBySprint(sprintId uuid.UUID, page, size int32, keyword string) ([]*model.WorkItem, error) {
	// get workitems from database
	var workItems []*model.WorkItem
	err := w.DB.Where("sprint_id = ?", sprintId).Where("title LIKE ?", "%"+keyword+"%").Limit(int(size)).Offset(int((page - 1) * size)).Preload("CreatedUser").Preload("AssignUser").Order("`order`").Order("created_at").Find(&workItems).Error
	if err != nil {
		return nil, err
	}
	return workItems, nil
}

func (w WorkItemDao) CountByProject(projectId uuid.UUID, keyword string) int64 {
	count := int64(0)
	w.DB.Where("project_id = ?", projectId).Where("title LIKE ?", "%"+keyword+"%").Model(&model.WorkItem{}).Count(&count)
	return count
}

func (w WorkItemDao) CountByProjectNotInSprint(projectId uuid.UUID, keyword string) int64 {
	count := int64(0)
	w.DB.Where("project_id = ?", projectId).Where("sprint_id IS NULL").Where("title LIKE ?", "%"+keyword+"%").Model(&model.WorkItem{}).Count(&count)
	return count
}

func (w WorkItemDao) CountBySprint(sprintId uuid.UUID, keyword string) int64 {
	count := int64(0)
	w.DB.Where("sprint_id = ?", sprintId).Where("title LIKE ?", "%"+keyword+"%").Model(&model.WorkItem{}).Count(&count)
	return count
}

func (w WorkItemDao) CountBySprints(sprintIds []uuid.UUID) []int64 {
	var counts []int64
	for _, sprintId := range sprintIds {
		count := int64(0)
		w.DB.Where("sprint_id = ?", sprintId).Model(&model.WorkItem{}).Count(&count)
		counts = append(counts, count)
	}
	return counts
}

func (w WorkItemDao) Create(workItem model.WorkItem) (*model.WorkItem, error) {
	// make transaction of create work item and work item log
	err := w.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&workItem).Error; err != nil {
			tx.Rollback()
			return err
		}
		var log model.ItemLog
		log.Id = uuid.New()
		log.ItemId = workItem.Id
		log.ItemType = "WORKITEM"
		log.Action = "CREATE"
		log.Log = "Create work item"
		log.Changes = "projectId: " + workItem.ProjectId.String() + ", sprintId: " + workItem.SprintId.String()
		log.CreatedBy = workItem.CreatedBy
		if err := tx.Create(&log).Error; err != nil {
			tx.Rollback()
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &workItem, nil
}

func (w WorkItemDao) Update(workItem model.WorkItem, userId uuid.UUID) (*model.WorkItem, error) {
	// make transaction of update work item and work item log
	err := w.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.WorkItem{}).Where("id = ?", workItem.Id).Updates(workItem).Error; err != nil {
			tx.Rollback()
			return err
		}
		var log model.ItemLog
		log.Id = uuid.New()
		log.ItemId = workItem.Id
		log.ItemType = "WORKITEM"
		log.Action = "UPDATE"
		log.Log = "Update work item"
		log.Changes = "projectId: " + workItem.ProjectId.String() + ", sprintId: " + workItem.SprintId.String()
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
	return &workItem, nil
}

func (w WorkItemDao) UpdateAssignUser(workItem model.WorkItem, userId uuid.UUID) (*model.WorkItem, error) {
	// make transaction of update work item assign user and work item log
	err := w.DB.Transaction(func(tx *gorm.DB) error {
		if workItem.AssignTo == uuid.Nil {
			if err := tx.Model(&model.WorkItem{}).Where("id = ?", workItem.Id).Update("assign_to", sql.NullString{}).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			if err := tx.Model(&model.WorkItem{}).Where("id = ?", workItem.Id).Update("assign_to", workItem.AssignTo).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
		var log model.ItemLog
		log.Id = uuid.New()
		log.ItemId = workItem.Id
		log.ItemType = "WORKITEM"
		log.Action = "UPDATE"
		log.Log = "Update work item, assign work item to the user, id is: " + workItem.AssignTo.String()
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

	return &workItem, nil
}

func (w WorkItemDao) UpdateStatus(workItem model.WorkItem, userId uuid.UUID) (*model.WorkItem, error) {
	// make transaction of update work item status and work item log
	err := w.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.WorkItem{}).Where("id = ?", workItem.Id).Update("status", workItem.Status).Error; err != nil {
			tx.Rollback()
			return err
		}
		var log model.ItemLog
		log.Id = uuid.New()
		log.ItemId = workItem.Id
		log.ItemType = "WORKITEM"
		log.Action = "UPDATE"
		log.Log = "Update work item status to: " + workItem.Status
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
	return &workItem, nil
}

func (w WorkItemDao) UpdateSprintWithTasks(workItem model.WorkItem, userId uuid.UUID) (*model.WorkItem, error) {
	// make transaction of update work item sprint and work item log
	err := w.DB.Transaction(func(tx *gorm.DB) error {
		// get workitem from database
		var oldWorkItem model.WorkItem
		if err := tx.Where("id = ?", workItem.Id).Find(&oldWorkItem).Error; err != nil {
			tx.Rollback()
			return err
		}
		if workItem.SprintId == uuid.Nil {
			if err := tx.Model(&model.WorkItem{}).Where("id = ?", workItem.Id).Update("sprint_id", sql.NullString{}).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			if err := tx.Model(&model.WorkItem{}).Where("id = ?", workItem.Id).Update("sprint_id", workItem.SprintId).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
		var log model.ItemLog
		log.Id = uuid.New()
		log.ItemId = workItem.Id
		log.ItemType = "WORKITEM"
		log.Action = "UPDATE"
		log.Log = "Update work item sprint, move to sprint: " + workItem.SprintId.String()
		log.Changes = "projectId: " + workItem.ProjectId.String() + ", sprintId: " + workItem.SprintId.String()
		log.CreatedBy = userId
		if err := tx.Create(&log).Error; err != nil {
			tx.Rollback()
			return err
		}
		// get tasks of work item
		var tasks []*model.Task
		if err := tx.Where("work_item_id = ?", workItem.Id).Find(&tasks).Error; err != nil {
			tx.Rollback()
			return err
		}
		if len(tasks) > 0 {
			// update tasks to database
			if workItem.SprintId == uuid.Nil {
				if err := tx.Model(&model.Task{}).Where("work_item_id = ?", workItem.Id).Update("sprint_id", sql.NullString{}).Error; err != nil {
					tx.Rollback()
					return err
				}
			} else {
				if err := tx.Model(&model.Task{}).Where("work_item_id = ?", workItem.Id).Update("sprint_id", workItem.SprintId).Error; err != nil {
					tx.Rollback()
					return err
				}
			}
			// add item log for each task
			var taskLogs []*model.ItemLog
			for _, task := range tasks {
				var log model.ItemLog
				log.Id = uuid.New()
				log.ItemId = task.Id
				log.ItemType = "TASK"
				log.Action = "UPDATE"
				log.Log = "Update task sprint, move to sprint: " + workItem.SprintId.String()
				log.Changes = "projectId: " + task.ProjectId.String() + ", sprintId: " + workItem.SprintId.String()
				log.CreatedBy = userId
				taskLogs = append(taskLogs, &log)
			}
			if err := tx.Create(&taskLogs).Error; err != nil {
				tx.Rollback()
				return err
			}
			// get task count that status is not Done or Removed
			needUpdatetaskCount := 0
			for _, task := range tasks {
				if task.Status != "Done" && task.Status != "Removed" {
					needUpdatetaskCount++
				}
			}
			// get now time
			createdOldDate := time.Now()
			// get old sprint statuses ordered by date
			var oldSprintStatuses []*model.SprintStatus
			if err := tx.Where("sprint_id = ?", oldWorkItem.SprintId).Order("sprint_date").Find(&oldSprintStatuses).Error; err != nil {
				tx.Rollback()
				return err
			}
			lastOldSprintStatusIndex := len(oldSprintStatuses) - 1

			// if createdDate before or equal the last sprint status date
			if createdOldDate.Before(oldSprintStatuses[lastOldSprintStatusIndex].SprintDate) || createdOldDate.Equal(oldSprintStatuses[lastOldSprintStatusIndex].SprintDate) {
				// if the first sprint status date is after createdDate, set createdDate to the first sprint status date
				if oldSprintStatuses[0].SprintDate.After(createdOldDate) {
					createdOldDate = oldSprintStatuses[0].SprintDate
				}
				// format createdDate to date
				correctOldDate := time.Date(createdOldDate.Year(), createdOldDate.Month(), createdOldDate.Day(), 0, 0, 0, 0, createdOldDate.Location())
				// update sprint status record and set work item count + 1
				if err := tx.Model(&model.SprintStatus{}).Where("sprint_id = ?", oldWorkItem.SprintId).Where("sprint_date = ?", correctOldDate).Update("task_count", gorm.Expr("task_count - ?", needUpdatetaskCount)).Error; err != nil {
					tx.Rollback()
					return err
				}
			}

			if workItem.SprintId != uuid.Nil {
				// get now time
				createdDate := time.Now()
				// get current sprint statuses ordered by date
				var currentSprintStatuses []*model.SprintStatus
				if err := tx.Where("sprint_id = ?", workItem.SprintId).Order("sprint_date").Find(&currentSprintStatuses).Error; err != nil {
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
					if err := tx.Model(&model.SprintStatus{}).Where("sprint_id = ?", workItem.SprintId).Where("sprint_date = ?", correctDate).Update("task_count", gorm.Expr("task_count + ?", needUpdatetaskCount)).Error; err != nil {
						tx.Rollback()
						return err
					}
				}
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &workItem, nil
}

func (w WorkItemDao) Delete(workItem model.WorkItem, userId uuid.UUID) (bool, error) {
	// make transaction of delete work item and work item log
	err := w.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", workItem.Id).Delete(&model.WorkItem{}).Error; err != nil {
			tx.Rollback()
			return err
		}
		var log model.ItemLog
		log.Id = uuid.New()
		log.ItemId = workItem.Id
		log.ItemType = "WORKITEM"
		log.Action = "DELETE"
		log.Log = "Delete work item"
		log.Changes = "projectId: " + workItem.ProjectId.String() + ", sprintId: " + workItem.SprintId.String()
		log.CreatedBy = userId
		if err := tx.Create(&log).Error; err != nil {
			tx.Rollback()
			return err
		}
		// get now time
		createdDate := time.Now()
		// get current sprint statuses ordered by date
		var currentSprintStatuses []*model.SprintStatus
		if err := tx.Where("sprint_id = ?", workItem.SprintId).Order("sprint_date").Find(&currentSprintStatuses).Error; err != nil {
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
			if err := tx.Model(&model.SprintStatus{}).Where("sprint_id = ?", workItem.SprintId).Where("sprint_date = ?", correctDate).Update("work_item_count", gorm.Expr("work_item_count - ?", 1)).Error; err != nil {
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

func (w WorkItemDao) ReOrder(workItemIds []int64) ([]int64, error) {
	// make transaction of reorder work items
	err := w.DB.Transaction(func(tx *gorm.DB) error {
		for i, id := range workItemIds {
			if err := tx.Model(&model.WorkItem{}).Where("id = ?", id).Update("order", i+1).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return workItemIds, nil
}

func NewWorkItemDao(d *gorm.DB) *WorkItemDao {
	return &WorkItemDao{d}
}
