package model

type Feature struct {
	Model

	ProjectId   int64   `gorm:"column:project_id"`
	SprintId    int64   `gorm:"column:sprint_id"`
	EpicId      int64   `gorm:"column:epic_id"`
	Title       string  `gorm:"column:title;size:1000"`
	Description string  `gorm:"column:description;size:5000"`
	AssignTo    int64   `gorm:"column:assign_to"`
	CreatedBy   int64   `gorm:"column:created_by"`
	AssignUser  User    `gorm:"foreignKey:AssignTo"`
	CreatedUser User    `gorm:"foreignKey:CreatedBy"`
	FromProject Project `gorm:"foreignKey:ProjectId"`
	FromSprint  Sprint  `gorm:"foreignKey:SprintId"`
	FromEpic    Epic    `gorm:"foreignKey:EpicId"`
}
