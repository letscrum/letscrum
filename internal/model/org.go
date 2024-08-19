package model

import "github.com/google/uuid"

type Org struct {
	BaseUUID

	Name              string    `gorm:"column:name;size:255;index:idx_name,unique" json:"name,omitempty"`
	DisplayName       string    `gorm:"column:display_name;size:255" json:"display_name,omitempty"`
	Description       string    `gorm:"column:description;size:5000" json:"description,omitempty"`
	MemberLimitation  int64     `gorm:"column:member_limitation;default:5" json:"member_limitation,omitempty"`
	ProjectLimitation int64     `gorm:"column:project_limitation;default:2" json:"project_limitation,omitempty"`
	CreatedBy         uuid.UUID `gorm:"column:created_by;varchar(191);default:null" json:"created_by,omitempty"`
	CreatedUser       User      `gorm:"foreignKey:CreatedBy" json:"created_user,omitempty"`
}
