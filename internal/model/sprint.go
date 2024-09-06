package model

import (
	"time"

	"github.com/google/uuid"
)

type Sprint struct {
	BaseUUID

	ProjectId    uuid.UUID `gorm:"column:project_id;varchar(191);default:null" json:"project_id,omitempty"`
	Name         string    `gorm:"column:name;size:255" json:"name,omitempty"`
	Members      string    `gorm:"column:members;size:5000" json:"members,omitempty"`
	StartDate    time.Time `gorm:"column:start_date;type:date" json:"start_date,omitempty"`
	EndDate      time.Time `gorm:"column:end_date;type:date" json:"end_date,omitempty"`
	BurndownType string    `gorm:"column:burndown_type;size:255" json:"burndown_type,omitempty"`
	FromProject  Project   `gorm:"foreignKey:ProjectId" json:"from_project,omitempty"`
}
