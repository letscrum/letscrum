package model

import "github.com/google/uuid"

type WorkItem struct {
	BaseId

	ProjectId   uuid.UUID `gorm:"column:project_id;varchar(191);default:null" json:"project_id,omitempty"`
	SprintId    uuid.UUID `gorm:"column:sprint_id;varchar(191);default:null" json:"sprint_id,omitempty"`
	FeatureId   int64     `gorm:"column:feature_id;default:null" json:"feature_id,omitempty"`
	Title       string    `gorm:"column:title;size:1000" json:"title,omitempty"`
	Type        string    `gorm:"column:type;size:255" json:"type,omitempty"`
	Description string    `gorm:"column:description;size:5000" json:"description,omitempty"`
	Status      string    `gorm:"column:status;size:255" json:"status,omitempty"`
	AssignTo    uuid.UUID `gorm:"column:assign_to;varchar(191);default:null" json:"assign_to,omitempty"`
	CreatedBy   uuid.UUID `gorm:"column:created_by;varchar(191);default:null" json:"created_by,omitempty"`
	AssignUser  User      `gorm:"foreignKey:AssignTo" json:"assign_user,omitempty"`
	CreatedUser User      `gorm:"foreignKey:CreatedBy" json:"created_user,omitempty"`
	FromProject Project   `gorm:"foreignKey:ProjectId" json:"from_project,omitempty"`
	FromSprint  Sprint    `gorm:"foreignKey:SprintId" json:"from_sprint,omitempty"`
	FromFeature Feature   `gorm:"foreignKey:FeatureId" json:"from_feature,omitempty"`
}
