package model

import (
	"time"

	"github.com/google/uuid"
)

type SprintBurndown struct {
	SprintId   uuid.UUID `gorm:"column:sprint_id;type:uuid;primaryKey" json:"sprint_id,omitempty"`
	SprintDate time.Time `gorm:"column:sprint_date;type:date;primaryKey" json:"sprint_date,omitempty"`
	TaskCount  int32     `gorm:"column:task_count" json:"task_count,omitempty"`
	WorkHours  float32   `gorm:"column:work_hours" json:"work_hours,omitempty"`
	ForSprint  Sprint    `gorm:"foreignKey:SprintId" json:"for_sprint,omitempty"`
}
