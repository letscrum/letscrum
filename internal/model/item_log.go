package model

import "github.com/google/uuid"

type ItemLog struct {
	BaseUUID

	ItemType    string    `gorm:"column:item_type;size:255" json:"item_type,omitempty"`
	ItemId      int64     `gorm:"column:item_id" json:"item_id,omitempty"`
	Action      string    `gorm:"column:action;size:255" json:"action,omitempty"`
	Log         string    `gorm:"column:log;size:1000" json:"log,omitempty"`
	CreatedBy   uuid.UUID `gorm:"column:created_by;varchar(191);default:null" json:"created_by,omitempty"`
	CreatedUser User      `gorm:"foreignKey:CreatedBy" json:"created_user,omitempty"`
}
