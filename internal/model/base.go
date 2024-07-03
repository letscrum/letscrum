package model

import "time"

type Model struct {
	Id        int64     `gorm:"column:id;primary_key"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}
