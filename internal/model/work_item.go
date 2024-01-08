package model

type WorkItem struct {
	Model

	ProjectID   int64  `gorm:"column:project_id"`
	SprintID    int64  `gorm:"column:sprint_id"`
	FeatureID   int64  `gorm:"column:feature_id"`
	Title       string `gorm:"column:title;size:1000"`
	Type        string `gorm:"column:type;size:255"`
	Description string `gorm:"column:description;size:5000"`
	Status      string `gorm:"column:status;size:255"`
	AssignTo    int64  `gorm:"column:assign_to"`
	CreatedBy   int64  `gorm:"column:created_by"`
	AssignUser  User   `gorm:"foreignKey:AssignTo"`
	CreatedUser User   `gorm:"foreignKey:CreatedBy"`
}
