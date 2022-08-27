package mysql

import (
	"context"
	"github.com/letscrum/letscrum/internal/model"
	"gorm.io/gorm"
)

type LetscrumDao struct {
	Db *gorm.DB
}

func (d LetscrumDao) SignIn(ctx context.Context, name, password string) (*model.User, error) {
	var u *model.User
	if err := d.Db.Where("name = ?", name).Where("password = ?", password).Find(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func NewLetscrumDao(d *gorm.DB) *LetscrumDao {
	return &LetscrumDao{d}
}
