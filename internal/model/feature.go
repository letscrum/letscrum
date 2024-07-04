package model

type Feature struct {
	Model

	ProjectId   int64   `gorm:"column:project_id" json:"is_super_admin,omitempty"`
	SprintId    int64   `gorm:"column:sprint_id" json:"sprint_id,omitempty"`
	EpicId      int64   `gorm:"column:epic_id" json:"epic_id,omitempty"`
	Title       string  `gorm:"column:title;size:1000" json:"title,omitempty"`
	Description string  `gorm:"column:description;size:5000" json:"description,omitempty"`
	AssignTo    int64   `gorm:"column:assign_to" json:"assign_to,omitempty"`
	CreatedBy   int64   `gorm:"column:created_by" json:"created_by,omitempty"`
	AssignUser  User    `gorm:"foreignKey:AssignTo;constraint:false" json:"assign_user,omitempty"`
	CreatedUser User    `gorm:"foreignKey:CreatedBy" json:"created_user,omitempty"`
	FromProject Project `gorm:"foreignKey:ProjectId" json:"from_project,omitempty"`
	FromSprint  Sprint  `gorm:"foreignKey:SprintId;constraint:false" json:"from_sprint,omitempty"`
	FromEpic    Epic    `gorm:"foreignKey:EpicId;constraint:false" json:"from_epic,omitempty"`
}
