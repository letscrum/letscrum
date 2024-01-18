package dao

import "github.com/letscrum/letscrum/internal/model"

type SprintMemberDao interface {
	ListBySprint(sprint model.Sprint, page, size int32) ([]*model.SprintMember, error)
	CountBySprint(sprint model.Sprint) int64
	BatchAdd(sprintMembers []model.SprintMember) ([]*model.SprintMember, error)
	BatchUpdate(sprintMembers []model.SprintMember) ([]*model.SprintMember, error)
	BatchRemove(sprintMembers []model.SprintMember) ([]*model.SprintMember, error)
}
