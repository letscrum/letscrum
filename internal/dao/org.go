package dao

import (
	"github.com/letscrum/letscrum/internal/model"
)

type OrgDao interface {
	Get(org model.Org) (*model.Org, error)
	List(page, size int32, keyword string) ([]*model.Org, error)
	Count(keyword string) int64
	Create(org model.Org) (*model.Org, error)
	Update(org model.Org) (*model.Org, error)
	Delete(org model.Org) (bool, error)
}
