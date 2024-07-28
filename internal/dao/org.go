package dao

import (
	"github.com/letscrum/letscrum/internal/model"
)

type OrgDao interface {
	Get(org model.Org) (model.Org, error)
	List(page, size int32, keyword string) ([]model.Org, error)
	Count(keyword string) int64
	ListVisibleOrg(page, size int32, keyword string, user model.User) ([]model.Org, error)
	CountVisibleOrg(keyword string, user model.User) int64
	CountByUser(user model.User) int64
	Create(org model.Org) (model.Org, error)
	Update(org model.Org) (model.Org, error)
	Delete(org model.Org) (bool, error)
	AddMember(orgUser model.OrgUser) (model.OrgUser, error)
	AddMembers(orgUsers []model.OrgUser) ([]model.OrgUser, error)
	RemoveMember(orgUser model.OrgUser) (bool, error)
	ListMember(org model.Org) ([]model.OrgUser, error)
	SetAdmin(orgUser model.OrgUser, isAdmin bool) (model.OrgUser, error)
}
