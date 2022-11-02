package dao

import "github.com/letscrum/letscrum/internal/model"

type ProjectMemberDao interface {
	List(projectId int64, page, size int32) ([]*model.ProjectMember, error)
	Count() int64
	Add(projectId, userId int64) (bool, error)
	Update(projectId, userId int64, isAdmin bool) (bool, error)
	Remove(projectId, userId int64) (bool, error)
}
