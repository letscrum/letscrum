package model

import "github.com/google/uuid"

type Task struct {
	BaseId

	ProjectId    uuid.UUID `gorm:"column:project_id;varchar(191);default:null" json:"project_id,omitempty"`
	SprintId     uuid.UUID `gorm:"column:sprint_id;varchar(191);default:null" json:"sprint_id,omitempty"`
	WorkItemId   int64     `gorm:"column:work_item_id" json:"work_item_id,omitempty"`
	Title        string    `gorm:"column:title;size:1000" json:"title,omitempty"`
	Description  string    `gorm:"column:description;size:5000" json:"description,omitempty"`
	Status       string    `gorm:"column:status;size:255" json:"status,omitempty"`
	AssignTo     uuid.UUID `gorm:"column:assign_to;varchar(191);default:null" json:"assign_to,omitempty"`
	Remaining    float32   `gorm:"column:remaining" json:"remaining,omitempty"`
	CreatedBy    uuid.UUID `gorm:"column:created_by;varchar(191);default:null" json:"created_by,omitempty"`
	Order        int32     `gorm:"column:order;default:0" json:"order,omitempty"`
	AssignUser   User      `gorm:"foreignKey:AssignTo" json:"assign_user,omitempty"`
	CreatedUser  User      `gorm:"foreignKey:CreatedBy" json:"created_user,omitempty"`
	FromProject  Project   `gorm:"foreignKey:ProjectId" json:"from_project,omitempty"`
	FromSprint   Sprint    `gorm:"foreignKey:SprintId" json:"from_sprint,omitempty"`
	FromWorkItem WorkItem  `gorm:"foreignKey:WorkItemId" json:"from_work_item,omitempty"`
}
