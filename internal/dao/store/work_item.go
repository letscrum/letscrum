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

func (w WorkItemDao) Get(workItem model.WorkItem) (*model.WorkItem, error) {
	// get workitem from database
	if err := w.DB.Where("id = ?", workItem.Id).Find(&workItem).Error; err != nil {
		return nil, err
	}
	return &workItem, nil
}

func (w WorkItemDao) ListByProject(projectId uuid.UUID, page, size int32, keyword string) ([]*model.WorkItem, error) {
	// get workitems from database
	var workItems []*model.WorkItem
	err := w.DB.Where("project_id = ?", projectId).Where("title LIKE ?", "%"+keyword+"%").Limit(int(size)).Offset(int((page - 1) * size)).Preload("CreatedUser").Preload("AssignUser").Order("created_at").Find(&workItems).Error
	if err != nil {
		return nil, err
	}
	return workItems, nil
}

func (w WorkItemDao) ListBySprint(sprintId uuid.UUID, page, size int32, keyword string) ([]*model.WorkItem, error) {
	// get workitems from database
	var workItems []*model.WorkItem
	err := w.DB.Where("sprint_id = ?", sprintId).Where("title LIKE ?", "%"+keyword+"%").Limit(int(size)).Offset(int((page - 1) * size)).Preload("CreatedUser").Preload("AssignUser").Order("created_at").Find(&workItems).Error
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
	return &workItem, nil
}

func (w WorkItemDao) Update(workItem model.WorkItem) (*model.WorkItem, error) {
	if err := w.DB.Model(&model.WorkItem{}).Where("id = ?", workItem.Id).Update("assign_to", workItem.AssignTo).Error; err != nil {
		return nil, err
	}
	return &workItem, nil
}

func (w WorkItemDao) UpdateAssignUser(workItem model.WorkItem) (*model.WorkItem, error) {
	if workItem.AssignTo == uuid.Nil {
		if err := w.DB.Model(&model.WorkItem{}).Where("id = ?", workItem.Id).Update("assign_to", sql.NullString{}).Error; err != nil {
			return nil, err
		}
		return &workItem, nil
	}
	if err := w.DB.Model(&model.WorkItem{}).Where("id = ?", workItem.Id).Update("assign_to", workItem.AssignTo).Error; err != nil {
		return nil, err
	}
	return &workItem, nil
}

func (w WorkItemDao) Delete(workItem model.WorkItem) (bool, error) {
	// delete work item from database
	if err := w.DB.Where("id = ?", workItem.Id).Delete(&model.WorkItem{}).Error; err != nil {
		return false, err
	}
	return true, nil
}

func NewWorkItemDao(d *gorm.DB) *WorkItemDao {
	return &WorkItemDao{d}
}
