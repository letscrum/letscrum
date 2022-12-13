package model

type Task struct {
	Model

	ProjectID   int64   `json:"project_id"`
	SprintID    int64   `json:"sprint_id"`
	WorkItemID  int64   `json:"work_item_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	AssignTo    int64   `json:"assign_to"`
	Remaining   float32 `json:"remaining"`
	CreatedBy   int64   `json:"created_by"`
	AssignUser  User    `gorm:"foreignKey:AssignTo"`
	CreatedUser User    `gorm:"foreignKey:CreatedBy"`
}
