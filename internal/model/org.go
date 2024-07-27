package model

type Org struct {
	BaseUUID

	Name        string `gorm:"column:name;size:255" json:"name,omitempty"`
	DisplayName string `gorm:"column:display_name;size:255" json:"display_name,omitempty"`
	Description string `gorm:"column:description;size:5000" json:"description,omitempty"`
}
