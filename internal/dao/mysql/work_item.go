package mysql

import (
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

func (w WorkItemDao) ListByProject(projectId int64, page, size int32, keyword string) ([]*model.WorkItem, error) {
	// get workitems from database
	var workItems []*model.WorkItem
	err := w.DB.Where("project_id = ?", projectId).Where("title LIKE ?", "%"+keyword+"%").Limit(int(size)).Offset(int((page - 1) * size)).Preload("CreatedUser").Preload("AssignUser").Order("created_at").Find(&workItems).Error
	if err != nil {
		return nil, err
	}
	return workItems, nil
}

func (w WorkItemDao) ListBySprint(sprintId int64, page, size int32, keyword string) ([]*model.WorkItem, error) {
	// get workitems from database
	var workItems []*model.WorkItem
	err := w.DB.Where("sprint_id = ?", sprintId).Where("title LIKE ?", "%"+keyword+"%").Limit(int(size)).Offset(int((page - 1) * size)).Preload("CreatedUser").Preload("AssignUser").Order("created_at").Find(&workItems).Error
	if err != nil {
		return nil, err
	}
	return workItems, nil
}

func (w WorkItemDao) CountByProject(projectId int64, keyword string) int64 {
	count := int64(0)
	w.DB.Where("project_id = ?", projectId).Where("title LIKE ?", "%"+keyword+"%").Model(&model.WorkItem{}).Count(&count)
	return count
}

func (w WorkItemDao) CountBySprint(sprintId int64, keyword string) int64 {
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
	// update work item to database
	if err := w.DB.Model(&model.WorkItem{}).Where("id = ?", workItem.Id).Updates(workItem).Error; err != nil {
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
