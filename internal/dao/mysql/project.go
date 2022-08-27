package mysql

import (
	"context"
	"github.com/letscrum/letscrum/internal/model"
	"gorm.io/gorm"
)

type ProjectDao struct {
	Db *gorm.DB
}

func (d *ProjectDao) Create(ctx context.Context, project *model.Project) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (d *ProjectDao) Update(ctx context.Context, project *model.Project) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (d *ProjectDao) Delete(ctx context.Context, id int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (d *ProjectDao) Count(context.Context) int64 {
	count := int64(0)
	d.Db.Model(&model.Project{}).Count(&count)
	return count
}

func (d *ProjectDao) List(ctx context.Context, page, pageSize int32) ([]*model.Project, error) {
	var projects []*model.Project
	err := d.Db.Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Preload("CreatedUser").Find(&projects).Error
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (d *ProjectDao) Get(ctx context.Context, id int64) (*model.Project, error) {
	var p *model.Project
	if err := d.Db.Where("id = ?", id).Preload("CreatedUser").Find(&p).Error; err != nil {
		return nil, err
	}
	return p, nil
}

func NewProjectDao(d *gorm.DB) *ProjectDao {
	return &ProjectDao{d}
}
