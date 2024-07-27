package model

type User struct {
	BaseUUID

	Name         string `gorm:"column:name;size:255;index:idx_name,unique" json:"name,omitempty"`
	Email        string `gorm:"column:email;size:255;index:idx_name,unique" json:"email,omitempty"`
	Password     string `gorm:"column:password;size:255" json:"password,omitempty"`
	IsSuperAdmin bool   `gorm:"column:is_super_admin" json:"is_super_admin,omitempty"`
}
