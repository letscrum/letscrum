package model

type Task struct {
	Model

	ProjectID    int64    `gorm:"column:project_id"`
	SprintID     int64    `gorm:"column:sprint_id"`
	WorkItemID   int64    `gorm:"column:work_item_id"`
	Title        string   `gorm:"column:title;size:1000"`
	Description  string   `gorm:"column:description;size:5000"`
	Status       string   `gorm:"column:status;size:255"`
	AssignTo     int64    `gorm:"column:assign_to"`
	Remaining    float32  `gorm:"column:remaining"`
	CreatedBy    int64    `gorm:"column:created_by"`
	AssignUser   User     `gorm:"foreignKey:AssignTo"`
	CreatedUser  User     `gorm:"foreignKey:CreatedBy"`
	FromProject  Project  `gorm:"foreignKey:ProjectID"`
	FromSprint   Sprint   `gorm:"foreignKey:SprintID"`
	FromWorkItem WorkItem `gorm:"foreignKey:WorkItemID"`
}
