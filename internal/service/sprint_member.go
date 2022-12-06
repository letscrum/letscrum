package service

import (
	v1 "github.com/letscrum/letscrum/api/letscrum/v1"
	"github.com/letscrum/letscrum/internal/dao"
)

type SprintMemberService struct {
	v1.UnimplementedSprintMemberServer
	sprintMemberDao dao.SprintMemberDao
}

func NewSprintMemberService(dao dao.Interface) *SprintMemberService {
	return &SprintMemberService{
		sprintMemberDao: dao.SprintMemberDao(),
	}
}

//
//func Create(sprintId int64, userId int64) (int64, error) {
//	sprintMemberId, err := sprintmembermodel.CreateSprintMember(sprintId, userId)
//	if err != nil {
//		return 0, err
//	}
//	return sprintMemberId, nil
//}
//
//func ListSprintMemberBySprint(sprintId int64, page int32, pageSize int32) ([]*userV1.User, int64, error) {
//	sprintMembers, err := sprintmembermodel.ListSprintMemberBySprint(sprintId, page, pageSize)
//	if err != nil {
//		return nil, 0, err
//	}
//	var members []*userV1.User
//	for _, m := range sprintMembers {
//		members = append(members, &userV1.User{
//			Id:           m.User.Id,
//			Name:         m.User.Name,
//			Email:        m.User.Email,
//			IsSuperAdmin: m.User.IsSuperAdmin,
//			RoleName:     m.Role.Name,
//			Capacity:     m.Capacity,
//		})
//	}
//	count := sprintmembermodel.CountSprintMemberBySprint(sprintId)
//	return members, count, nil
//}
//
//func ListSprintByUser(userId int64, page int32, pageSize int32) ([]*projectV1.Sprint, int64, error) {
//	sprintMembers, err := sprintmembermodel.ListSprintMemberByUser(userId, page, pageSize)
//	if err != nil {
//		return nil, 0, err
//	}
//	var sprints []*projectV1.Sprint
//	for _, s := range sprintMembers {
//		sprints = append(sprints, &projectV1.Sprint{
//			Id:        s.Sprint.Id,
//			Name:      s.Sprint.Name,
//			StartDate: s.Sprint.StartDate.Unix(),
//			EndDate:   s.Sprint.EndDate.Unix(),
//		})
//	}
//	count := sprintmembermodel.CountSprintMemberByUser(userId)
//	return sprints, count, nil
//}
