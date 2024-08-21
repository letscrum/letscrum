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
	// create work item to database
	if err := w.DB.Create(&workItem).Error; err != nil {
		return nil, err
	}
	var log model.ItemLog
	log.Id = uuid.New()
	log.ItemId = workItem.Id
	log.ItemType = "WORKITEM"
	log.Action = "CREATE"
	log.Log = "Create work item"
	log.CreatedBy = workItem.CreatedBy
	if err := w.DB.Create(&log).Error; err != nil {
		return nil, err
	}
	return &workItem, nil
}

func (w WorkItemDao) Update(workItem model.WorkItem, userId uuid.UUID) (*model.WorkItem, error) {
	if err := w.DB.Model(&model.WorkItem{}).Where("id = ?", workItem.Id).Update("assign_to", workItem.AssignTo).Error; err != nil {
		return nil, err
	}
	var log model.ItemLog
	log.Id = uuid.New()
	log.ItemId = workItem.Id
	log.ItemType = "WORKITEM"
	log.Action = "UPDATE"
	log.Log = "Update work item"
	log.CreatedBy = userId
	if err := w.DB.Create(&log).Error; err != nil {
		return nil, err
	}
	return &workItem, nil
}

func (w WorkItemDao) UpdateAssignUser(workItem model.WorkItem, userId uuid.UUID) (*model.WorkItem, error) {
	if workItem.AssignTo == uuid.Nil {
		if err := w.DB.Model(&model.WorkItem{}).Where("id = ?", workItem.Id).Update("assign_to", sql.NullString{}).Error; err != nil {
			return nil, err
		}
	} else {
		if err := w.DB.Model(&model.WorkItem{}).Where("id = ?", workItem.Id).Update("assign_to", workItem.AssignTo).Error; err != nil {
			return nil, err
		}
	}
	var log model.ItemLog
	log.Id = uuid.New()
	log.ItemId = workItem.Id
	log.ItemType = "WORKITEM"
	log.Action = "UPDATE"
	log.Log = "Update work item, assign work item to the user, id is: " + workItem.AssignTo.String()
	log.CreatedBy = userId
	if err := w.DB.Create(&log).Error; err != nil {
		return nil, err
	}
	return &workItem, nil
}

func (w WorkItemDao) UpdateStatus(workItem model.WorkItem, userId uuid.UUID) (*model.WorkItem, error) {
	if err := w.DB.Model(&model.WorkItem{}).Where("id = ?", workItem.Id).Update("status", workItem.Status).Error; err != nil {
		return nil, err
	}
	var log model.ItemLog
	log.Id = uuid.New()
	log.ItemId = workItem.Id
	log.ItemType = "WORKITEM"
	log.Action = "UPDATE"
	log.Log = "Update work item status to: " + workItem.Status
	log.CreatedBy = userId
	if err := w.DB.Create(&log).Error; err != nil {
		return nil, err
	}
	return &workItem, nil
}

func (w WorkItemDao) UpdateSprintWithTasks(workItem model.WorkItem, userId uuid.UUID) (*model.WorkItem, error) {
	// update work item to database
	if err := w.DB.Model(&model.WorkItem{}).Where("id = ?", workItem.Id).Update("sprint_id", workItem.SprintId).Error; err != nil {
		return nil, err
	}
	var log model.ItemLog
	log.Id = uuid.New()
	log.ItemId = workItem.Id
	log.ItemType = "WORKITEM"
	log.Action = "UPDATE"
	log.Log = "Update work item sprint, move to sprint: " + workItem.SprintId.String()
	log.CreatedBy = userId
	if err := w.DB.Create(&log).Error; err != nil {
		return nil, err
	}
	// update tasks to database
	if err := w.DB.Model(&model.Task{}).Where("work_item_id = ?", workItem.Id).Update("sprint_id", workItem.SprintId).Error; err != nil {
		return nil, err
	}
	// TODO: add task logs
	return &workItem, nil
}

func (w WorkItemDao) Delete(workItem model.WorkItem, userId uuid.UUID) (bool, error) {
	// delete work item from database
	if err := w.DB.Where("id = ?", workItem.Id).Delete(&model.WorkItem{}).Error; err != nil {
		return false, err
	}
	var log model.ItemLog
	log.Id = uuid.New()
	log.ItemId = workItem.Id
	log.ItemType = "WORKITEM"
	log.Action = "DELETE"
	log.Log = "Delete work item"
	log.CreatedBy = userId
	if err := w.DB.Create(&log).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (w WorkItemDao) ReOrder(workItemIds []int64) ([]int64, error) {
	// re-order work items
	for i, id := range workItemIds {
		if err := w.DB.Model(&model.WorkItem{}).Where("id = ?", id).Update("order", i+1).Error; err != nil {
			return nil, err
		}
	}
	return workItemIds, nil
}

func NewWorkItemDao(d *gorm.DB) *WorkItemDao {
	return &WorkItemDao{d}
}
