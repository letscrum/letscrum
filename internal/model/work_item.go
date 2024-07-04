package model

type WorkItem struct {
	Model

	ProjectId   int64   `gorm:"column:project_id" json:"project_id,omitempty"`
	SprintId    int64   `gorm:"column:sprint_id;default:null" json:"sprint_id,omitempty"`
	FeatureId   int64   `gorm:"column:feature_id;default:null" json:"feature_id,omitempty"`
	Title       string  `gorm:"column:title;size:1000" json:"title,omitempty"`
	Type        string  `gorm:"column:type;size:255" json:"type,omitempty"`
	Description string  `gorm:"column:description;size:5000" json:"description,omitempty"`
	Status      string  `gorm:"column:status;size:255" json:"status,omitempty"`
	AssignTo    int64   `gorm:"column:assign_to;default:null" json:"assign_to,omitempty"`
	CreatedBy   int64   `gorm:"column:created_by" json:"created_by,omitempty"`
	AssignUser  User    `gorm:"foreignKey:AssignTo" json:"assign_user,omitempty"`
	CreatedUser User    `gorm:"foreignKey:CreatedBy" json:"created_user,omitempty"`
	FromProject Project `gorm:"foreignKey:ProjectId" json:"from_project,omitempty"`
	FromSprint  Sprint  `gorm:"foreignKey:SprintId" json:"from_sprint,omitempty"`
	FromFeature Feature `gorm:"foreignKey:FeatureId" json:"from_feature,omitempty"`
}
