package dao

import "github.com/letscrum/letscrum/internal/model"

type SprintMemberDao interface {
	List(sprintID int64, page, size int32) ([]*model.SprintMember, error)
	Count(sprintID int64) int64
	Add(sprintID int64, userID int64) (bool, error)
	Update(sprintID, userID int64) (bool, error)
	Remove(sprintID, userID int64) (bool, error)
}
