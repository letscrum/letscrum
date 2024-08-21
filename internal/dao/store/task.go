package store

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/letscrum/letscrum/internal/model"
	"gorm.io/gorm"
)

type TaskDao struct {
	DB *gorm.DB
}

func (t TaskDao) ListByWorkItemIds(workItemIds []int64) ([]*model.Task, error) {
	var tasks []*model.Task
	err := t.DB.Where("work_item_id IN ?", workItemIds).Preload("CreatedUser").Preload("AssignUser").Find(&tasks).Error
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
	err := t.DB.Where("title LIKE ?", "%"+keyword+"%").Limit(int(size)).Offset(int((page - 1) * size)).Preload("CreatedUser").Preload("AssignUser").Order("created_at").Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (t TaskDao) ListByWorkItem(workItemId int64, page, size int32, keyword string) ([]*model.Task, error) {
	var tasks []*model.Task
	err := t.DB.Where("work_item_id = ?", workItemId).Where("title LIKE ?", "%"+keyword+"%").Limit(int(size)).Offset(int((page - 1) * size)).Preload("CreatedUser").Preload("AssignUser").Order("created_at").Find(&tasks).Error
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

func (t TaskDao) Create(task model.Task) (*model.Task, error) {
	if err := t.DB.Create(&task).Error; err != nil {
		return nil, err
	}
	var log model.ItemLog
	log.Id = uuid.New()
	log.ItemId = task.Id
	log.ItemType = "TASK"
	log.Action = "CREATE"
	log.Log = "Create task"
	log.CreatedBy = task.CreatedBy
	if err := t.DB.Create(&log).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (t TaskDao) Update(task model.Task, userId uuid.UUID) (*model.Task, error) {
	if err := t.DB.Model(&model.Task{}).Where("id = ?", task.Id).Updates(task).Error; err != nil {
		return nil, err
	}
	// add update log
	var log model.ItemLog
	log.Id = uuid.New()
	log.ItemId = task.Id
	log.ItemType = "TASK"
	log.Action = "UPDATE"
	log.Log = "Update task"
	log.CreatedBy = userId
	if err := t.DB.Create(&log).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (t TaskDao) UpdateStatus(task model.Task, userId uuid.UUID) (*model.Task, error) {
	if err := t.DB.Model(&model.Task{}).Where("id = ?", task.Id).Update("status", task.Status).Error; err != nil {
		return nil, err
	}
	var log model.ItemLog
	log.Id = uuid.New()
	log.ItemId = task.Id
	log.ItemType = "TASK"
	log.Action = "UPDATE"
	log.Log = "Update task status to: " + task.Status
	log.CreatedBy = userId
	if err := t.DB.Create(&log).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (t TaskDao) UpdateAssignUser(task model.Task, userId uuid.UUID) (*model.Task, error) {
	if task.AssignTo == uuid.Nil {
		if err := t.DB.Model(&model.Task{}).Where("id = ?", task.Id).Update("assign_to", sql.NullString{}).Error; err != nil {
			return nil, err
		}
	} else {
		if err := t.DB.Model(&model.Task{}).Where("id = ?", task.Id).Update("assign_to", task.AssignTo).Error; err != nil {
			return nil, err
		}
	}
	var log model.ItemLog
	log.Id = uuid.New()
	log.ItemId = task.Id
	log.ItemType = "TASK"
	log.Action = "UPDATE"
	log.Log = "Update task, assign task to the user, id is: " + task.AssignTo.String()
	log.CreatedBy = userId
	if err := t.DB.Create(&log).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (t TaskDao) Move(task model.Task, userId uuid.UUID) (*model.Task, error) {
	if err := t.DB.Model(&model.Task{}).Where("id = ?", task.Id).Update("status", task.Status).Update("work_item_id", task.WorkItemId).Error; err != nil {
		return nil, err
	}
	var log model.ItemLog
	log.Id = uuid.New()
	log.ItemId = task.Id
	log.ItemType = "TASK"
	log.Action = "UPDATE"
	log.Log = "Update task status and move task to work item"
	log.CreatedBy = userId
	if err := t.DB.Create(&log).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (t TaskDao) Delete(task model.Task, userId uuid.UUID) (bool, error) {
	if err := t.DB.Where("id = ?", task.Id).Delete(&model.Task{}).Error; err != nil {
		return false, err
	}
	var log model.ItemLog
	log.Id = uuid.New()
	log.ItemId = task.Id
	log.ItemType = "TASK"
	log.Action = "DELETE"
	log.Log = "Delete task"
	log.CreatedBy = userId
	if err := t.DB.Create(&log).Error; err != nil {
		return false, err
	}
	return true, nil
}

func NewTaskDao(d *gorm.DB) *TaskDao {
	return &TaskDao{d}
}
