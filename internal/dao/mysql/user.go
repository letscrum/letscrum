package mysql

import "gorm.io/gorm"

type UserDao struct {
	Db *gorm.DB
}

func NewUserDao(d *gorm.DB) *UserDao {
	return &UserDao{d}
}
