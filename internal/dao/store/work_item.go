package store

import (
	"database/sql"

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
	err := w.DB.Where("project_id = ?", projectId).Where("title LIKE ?", "%"+keyword+"%").Limit(int(size)).Offset(int((page - 1) * size)).Preload("CreatedUser").Preload("AssignUser").Order("`order`").Order("created_at").Find(&workItems).Error
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

func (w WorkItemDao) CountBySprint(sprintId uuid.UUID, keyword string) int64 {
	count := int64(0)
	w.DB.Where("sprint_id = ?", sprintId).Where("title LIKE ?", "%"+keyword+"%").Model(&model.WorkItem{}).Count(&count)
	return count
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
		if err := tx.Model(&model.WorkItem{}).Where("id = ?", workItem.Id).Update("sprint_id", workItem.SprintId).Error; err != nil {
			tx.Rollback()
			return err
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
		// update tasks to database
		if err := tx.Model(&model.Task{}).Where("work_item_id = ?", workItem.Id).Update("sprint_id", workItem.SprintId).Error; err != nil {
			tx.Rollback()
			return err
		}
		// TODO: add task logs
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
