package model

import "github.com/google/uuid"

type Project struct {
	BaseUUID

	Name        string    `gorm:"column:name;size:255" json:"name,omitempty"`
	DisplayName string    `gorm:"column:display_name;size:255" json:"display_name,omitempty"`
	Description string    `gorm:"column:description;size:5000" json:"description,omitempty"`
	Members     string    `gorm:"column:members;size:5000" json:"members,omitempty"`
	CreatedBy   uuid.UUID `gorm:"column:create_by;varchar(191);default:null" json:"created_by,omitempty"`
	CreatedUser User      `gorm:"foreignKey:CreatedBy" json:"created_user,omitempty"`
}
