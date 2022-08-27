package model

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	Id        int64          `json:"id" gorm:"primary_key"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"default:null"`
	CreatedAt time.Time      `json:"created_at" gorm:"default:null"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"default:null"`
}
