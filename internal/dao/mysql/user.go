package mysql

import (
	"github.com/letscrum/letscrum/internal/model"
	"gorm.io/gorm"
)

type UserDao struct {
	Db *gorm.DB
}

func (d *UserDao) Get(id int64) (*model.User, error) {
	var u *model.User
	if err := d.Db.Where("id = ?", id).Find(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func NewUserDao(d *gorm.DB) *UserDao {
	return &UserDao{d}
}
