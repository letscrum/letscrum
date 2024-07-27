package model

type Org struct {
	BaseUUID

	Name        string `gorm:"column:name;size:255" json:"name,omitempty"`
	Description string `gorm:"column:description;size:5000" json:"description,omitempty"`
	CreatedBy   int64  `gorm:"column:create_by" json:"created_by,omitempty"`
	CreatedUser User   `gorm:"foreignKey:CreatedBy" json:"created_user,omitempty"`
}
