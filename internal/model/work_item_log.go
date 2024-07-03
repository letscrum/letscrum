package model

type WorkItemLog struct {
	Model

	WorkItemId   int64    `gorm:"column:work_item_id"`
	Log          string   `gorm:"column:log;size:1000"`
	FromWorkItem WorkItem `gorm:"foreignKey:WorkItemId"`
}
