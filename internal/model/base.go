package model

import (
	"time"

	"github.com/google/uuid"
)

type BaseUUID struct {
	Id        uuid.UUID `gorm:"column:id;primaryKey;varchar(36)" json:"id,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at,omitempty"`
}

type BaseId struct {
	Id        int64     `gorm:"column:id;primaryKey" json:"id,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at,omitempty"`
}
