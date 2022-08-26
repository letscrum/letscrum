package mysql

import (
	"context"
	"github.com/letscrum/letscrum/internal/model"
	"gorm.io/gorm"
)

type ProjectDao struct {
	db *gorm.DB
}

func NewProjectDao(d *gorm.DB) *ProjectDao {
	return &ProjectDao{d}
}

func (u *ProjectDao) Get(ctx context.Context, id int64) (*model.Project, error) {
	var p *model.Project
	if err := model.DB.Where("id = ?", id).Preload("CreatedUser").Find(&p).Error; err != nil {
		return nil, err
	}
	return p, nil
}
