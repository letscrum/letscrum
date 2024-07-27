package model

type TaskLog struct {
	BaseUUID

	TaskId   int64    `gorm:"column:task_id" json:"task_id,omitempty"`
	Log      string   `gorm:"column:log;size:1000" json:"log,omitempty"`
	FromTask WorkItem `gorm:"foreignKey:TaskId" json:"from_task,omitempty"`
}
