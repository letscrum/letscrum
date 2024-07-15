package mysql

import (
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

func (t TaskDao) Get(task model.Task) (*model.Task, error) {
	if err := t.DB.Where("id = ?", task.Id).Find(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
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
	return &task, nil
}

func (t TaskDao) Update(task model.Task) (*model.Task, error) {
	if err := t.DB.Model(&model.Task{}).Where("id = ?", task.Id).Updates(task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (t TaskDao) UpdateStatus(task model.Task) (*model.Task, error) {
	if err := t.DB.Model(&model.Task{}).Where("id = ?", task.Id).Update("status", task.Status).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (t TaskDao) UpdateAssignUser(task model.Task) (*model.Task, error) {
	if err := t.DB.Model(&model.Task{}).Where("id = ?", task.Id).Update("assign_to", task.AssignTo).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (t TaskDao) Move(task model.Task) (*model.Task, error) {
	if err := t.DB.Model(&model.Task{}).Where("id = ?", task.Id).Update("status", task.Status).Update("work_item_id", task.WorkItemId).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (t TaskDao) Delete(task model.Task) (bool, error) {
	if err := t.DB.Where("id = ?", task.Id).Delete(&model.Task{}).Error; err != nil {
		return false, err
	}
	return true, nil
}

func NewTaskDao(d *gorm.DB) *TaskDao {
	return &TaskDao{d}
}
