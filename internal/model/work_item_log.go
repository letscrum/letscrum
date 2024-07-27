package model

type WorkItemLog struct {
	BaseUUID

	WorkItemId   int64    `gorm:"column:work_item_id" json:"work_item_id,omitempty"`
	Log          string   `gorm:"column:log;size:1000" json:"log,omitempty"`
	FromWorkItem WorkItem `gorm:"foreignKey:WorkItemId" json:"from_work_item,omitempty"`
}
