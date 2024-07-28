package model

import "github.com/google/uuid"

type OrgUser struct {
	OrgId   uuid.UUID `gorm:"column:org_id;type:uuid;primaryKey" json:"org_id,omitempty"`
	UserId  uuid.UUID `gorm:"column:user_id;type:uuid;primaryKey" json:"user_id,omitempty"`
	IsAdmin bool      `gorm:"column:is_admin" json:"is_admin,omitempty"`
	ForOrg  Org       `gorm:"foreignKey:OrgId" json:"for_org,omitempty"`
	Member  User      `gorm:"foreignKey:UserId" json:"member,omitempty"`
}
