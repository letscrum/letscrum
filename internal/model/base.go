package model

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        int64          `gorm:"column:id;primary_key"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;default:null"`
	CreatedAt time.Time      `gorm:"column:created_at;default:null"`
	UpdatedAt time.Time      `gorm:"column:updated_at;default:null"`
}
