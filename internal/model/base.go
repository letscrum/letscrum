package model

import "time"

type Model struct {
	Id        int64     `gorm:"column:id;primary_key" json:"id,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at,omitempty"`
}
