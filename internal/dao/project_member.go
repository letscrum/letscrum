package dao

import "github.com/letscrum/letscrum/internal/model"

type ProjectMemberDao interface {
	List(projectID int64, page, size int32) ([]*model.ProjectMember, error)
	Count() int64
	Add(projectID int64, userIDs []int64) (bool, error)
	Update(projectID, userID int64, isAdmin bool) (bool, error)
	Remove(projectID, userID int64) (bool, error)
}
