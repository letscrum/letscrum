package mysql

import (
	"gorm.io/gorm"
)

type LetscrumDao struct {
	Db *gorm.DB
}

func NewLetscrumDao(d *gorm.DB) *LetscrumDao {
	return &LetscrumDao{d}
}
