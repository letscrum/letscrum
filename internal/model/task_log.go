package model

type TaskLog struct {
	Model

	TaskID   int64    `gorm:"column:task_id"`
	Log      string   `gorm:"column:log;size:1000"`
	FromTask WorkItem `gorm:"foreignKey:TaskID"`
}
