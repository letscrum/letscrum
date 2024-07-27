package store

import (
	"strconv"

	"github.com/letscrum/letscrum/internal/model"
	"gorm.io/gorm"
)

type ProjectDao struct {
	DB *gorm.DB
}

func (d ProjectDao) Get(project model.Project) (*model.Project, error) {
	var p *model.Project
	if err := d.DB.Where("id = ?", project.Id).Preload("CreatedUser").Find(&p).Error; err != nil {
		return nil, err
	}
	return p, nil
}

func (d ProjectDao) Create(project model.Project) (*model.Project, error) {
	if err := d.DB.Create(&project).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (d ProjectDao) Update(project model.Project) (*model.Project, error) {
	if err := d.DB.Model(&model.Project{}).Where("id = ?", project.Id).Updates(project).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (d ProjectDao) Delete(project model.Project) (bool, error) {
	if err := d.DB.Where("id = ?", project.Id).Delete(&model.Project{}).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (d ProjectDao) Count(keyword string) int64 {
	count := int64(0)
	d.DB.Where("name LIKE ?", "%"+keyword+"%").Or("display_name LIKE ?", "%"+keyword+"%").Model(&model.Project{}).Count(&count)
	return count
}

func (d ProjectDao) List(page, size int32, keyword string) ([]*model.Project, error) {
	var projects []*model.Project
	err := d.DB.Where("name LIKE ?", "%"+keyword+"%").Or("display_name LIKE ?", "%"+keyword+"%").Limit(int(size)).Offset(int((page - 1) * size)).Preload("CreatedUser").Order("updated_at desc").Find(&projects).Error
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (d ProjectDao) ListVisibleProject(page, size int32, keyword string, user model.User) ([]*model.Project, error) {
	var projects []*model.Project
	if user.IsSuperAdmin {
		err := d.DB.Where("name LIKE ?", "%"+keyword+"%").Or("display_name LIKE ?", "%"+keyword+"%").Limit(int(size)).Offset(int((page - 1) * size)).Preload("CreatedUser").Order("updated_at desc").Find(&projects).Error
		if err != nil {
			return nil, err
		}
	} else {
		// select project where members include user.Id
		err := d.DB.Where("members LIKE ?", "%"+`{"user_id":`+strconv.FormatInt(user.Id, 10)+",%").Where("name LIKE ?", "%"+keyword+"%").Or("display_name LIKE ?", "%"+keyword+"%").Where("members LIKE ?", "%"+`{"user_id":`+strconv.FormatInt(user.Id, 10)+",%").Limit(int(size)).Offset(int((page - 1) * size)).Preload("CreatedUser").Order("updated_at desc").Find(&projects).Error
		if err != nil {
			return nil, err
		}
	}
	return projects, nil
}

func NewProjectDao(d *gorm.DB) *ProjectDao {
	return &ProjectDao{d}
}
