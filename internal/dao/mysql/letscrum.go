package mysql

import (
	"gorm.io/gorm"
)

type LetscrumDao struct {
	DB *gorm.DB
}

func NewLetscrumDao(d *gorm.DB) *LetscrumDao {
	return &LetscrumDao{d}
}
