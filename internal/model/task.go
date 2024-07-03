package model

type Task struct {
	Model

	ProjectId    int64    `gorm:"column:project_id"`
	SprintId     int64    `gorm:"column:sprint_id"`
	WorkItemId   int64    `gorm:"column:work_item_id"`
	Title        string   `gorm:"column:title;size:1000"`
	Description  string   `gorm:"column:description;size:5000"`
	Status       string   `gorm:"column:status;size:255"`
	AssignTo     int64    `gorm:"column:assign_to"`
	Remaining    float32  `gorm:"column:remaining"`
	CreatedBy    int64    `gorm:"column:created_by"`
	AssignUser   User     `gorm:"foreignKey:AssignTo"`
	CreatedUser  User     `gorm:"foreignKey:CreatedBy"`
	FromProject  Project  `gorm:"foreignKey:ProjectId"`
	FromSprint   Sprint   `gorm:"foreignKey:SprintId"`
	FromWorkItem WorkItem `gorm:"foreignKey:WorkItemId"`
}
