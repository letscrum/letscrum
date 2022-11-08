package mysql

import (
	"github.com/letscrum/letscrum/internal/model"
	"gorm.io/gorm"
)

type ProjectDao struct {
	Db *gorm.DB
}

func (d *ProjectDao) Create(project *model.Project) (bool, error) {
	if err := d.Db.Create(&project).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (d *ProjectDao) Update(project *model.Project) (bool, error) {
	if err := d.Db.Model(&model.Project{}).Where("id = ?", project.Id).Update("display_name", project.DisplayName).Error; err != nil {
		return false, nil
	}
	return true, nil
}

func (d *ProjectDao) Delete(id int64) (bool, error) {
	if err := d.Db.Where("id = ?", id).Delete(&model.Project{}).Error; err != nil {
		return false, nil
	}
	return true, nil
}

func (d *ProjectDao) Count(keyword string) int64 {
	count := int64(0)
	d.Db.Where("name LIKE ?", "%"+keyword+"%").Or("display_name LIKE ?", "%"+keyword+"%").Model(&model.Project{}).Count(&count)
	return count
}

func (d *ProjectDao) List(page, size int32, keyword string) ([]*model.Project, error) {
	var projects []*model.Project
	err := d.Db.Where("name LIKE ?", "%"+keyword+"%").Or("display_name LIKE ?", "%"+keyword+"%").Limit(int(size)).Offset(int((page - 1) * size)).Preload("CreatedUser").Find(&projects).Error
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (d *ProjectDao) Get(id int64) (*model.Project, error) {
	var p *model.Project
	if err := d.Db.Where("id = ?", id).Preload("CreatedUser").Find(&p).Error; err != nil {
		return nil, err
	}
	return p, nil
}

func NewProjectDao(d *gorm.DB) *ProjectDao {
	return &ProjectDao{d}
}
