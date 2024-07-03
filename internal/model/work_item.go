package model

type WorkItem struct {
	Model

	ProjectId   int64   `gorm:"column:project_id"`
	SprintId    int64   `gorm:"column:sprint_id"`
	FeatureId   int64   `gorm:"column:feature_id"`
	Title       string  `gorm:"column:title;size:1000"`
	Type        string  `gorm:"column:type;size:255"`
	Description string  `gorm:"column:description;size:5000"`
	Status      string  `gorm:"column:status;size:255"`
	AssignTo    int64   `gorm:"column:assign_to"`
	CreatedBy   int64   `gorm:"column:created_by"`
	AssignUser  User    `gorm:"foreignKey:AssignTo"`
	CreatedUser User    `gorm:"foreignKey:CreatedBy"`
	FromProject Project `gorm:"foreignKey:ProjectId"`
	FromSprint  Sprint  `gorm:"foreignKey:SprintId"`
	FromFeature Feature `gorm:"foreignKey:FeatureId"`
}
