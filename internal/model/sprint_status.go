package model

import (
	"time"

	"github.com/google/uuid"
)

type SprintStatus struct {
	SprintId             uuid.UUID `gorm:"column:sprint_id;type:uuid;primaryKey" json:"sprint_id,omitempty"`
	SprintDate           time.Time `gorm:"column:sprint_date;type:date;primaryKey" json:"sprint_date,omitempty"`
	StartWorkItemCount   int32     `gorm:"column:start_work_item_count" json:"start_work_item_count,omitempty"`
	CurrentWorkItemCount int32     `gorm:"column:current_work_item_count" json:"current_work_item_count,omitempty"`
	TotalWorkItemCount   int32     `gorm:"column:total_work_item_count" json:"total_work_item_count,omitempty"`
	StartTaskCount       int32     `gorm:"column:start_task_count" json:"start_task_count,omitempty"`
	CurrentTaskCount     int32     `gorm:"column:current_task_count" json:"current_task_count,omitempty"`
	TotalTaskCount       int32     `gorm:"column:total_task_count" json:"total_task_count,omitempty"`
	StartWorkHours       float64   `gorm:"column:start_work_hours" json:"start_work_hours,omitempty"`
	CurrentWorkHours     float64   `gorm:"column:current_work_hours" json:"current_work_hours,omitempty"`
	TotalWorkHours       float64   `gorm:"column:total_work_hours" json:"total_work_hours,omitempty"`
	ForSprint            Sprint    `gorm:"foreignKey:SprintId" json:"for_sprint,omitempty"`
}
