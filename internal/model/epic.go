package model

import "github.com/google/uuid"

type Epic struct {
	BaseId

	ProjectId   uuid.UUID `gorm:"column:project_id;varchar(191);default:null" json:"project_id,omitempty"`
	SprintId    uuid.UUID `gorm:"column:sprint_id;varchar(191);default:null" json:"sprint_id,omitempty"`
	Title       string    `gorm:"column:title;size:1000" json:"title,omitempty"`
	Description string    `gorm:"column:description;size:5000" json:"description,omitempty"`
	AssignTo    uuid.UUID `gorm:"column:assign_to;varchar(191);default:null" json:"assign_to,omitempty"`
	CreatedBy   uuid.UUID `gorm:"column:created_by;varchar(191);default:null" json:"created_by,omitempty"`
	AssignUser  User      `gorm:"foreignKey:AssignTo;constraint:false" json:"assign_user,omitempty"`
	CreatedUser User      `gorm:"foreignKey:CreatedBy" json:"created_user,omitempty"`
	FromProject Project   `gorm:"foreignKey:ProjectId" json:"from_project,omitempty"`
	FromSprint  Sprint    `gorm:"foreignKey:SprintId;constraint:false" json:"from_sprint,omitempty"`
}
