package v1

import (
	"context"

	generalv1 "github.com/letscrum/letscrum/api/general/v1"
	v1 "github.com/letscrum/letscrum/api/letscrum/v1"
	projectv1 "github.com/letscrum/letscrum/api/project/v1"
	"github.com/letscrum/letscrum/internal/dao"
	"github.com/letscrum/letscrum/internal/model"
	"github.com/letscrum/letscrum/pkg/utils"
	"github.com/spf13/cast"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SprintMemberService struct {
	v1.UnimplementedSprintMemberServer
	projectMemberDao dao.ProjectMemberDao
	sprintMemberDao  dao.SprintMemberDao
}

func NewSprintMemberService(dao dao.Interface) *SprintMemberService {
	return &SprintMemberService{
		projectMemberDao: dao.ProjectMemberDao(),
		sprintMemberDao:  dao.SprintMemberDao(),
	}
}

func (s *SprintMemberService) List(ctx context.Context, req *projectv1.ListSprintMemberRequest) (*projectv1.ListSprintMemberResponse, error) {
	_, err := utils.AuthJWT(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	members, err := s.sprintMemberDao.List(req.SprintId, 1, 999)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	var memberList []*projectv1.SprintMember
	for _, m := range members {
		var member = &projectv1.SprintMember{
			Id:        m.ID,
			UserId:    m.UserID,
			SprintId:  m.SprintID,
			UserName:  m.MemberUser.Name,
			UserEmail: m.MemberUser.Email,
			Role:      m.Role,
			Capacity:  m.Capacity,
		}
		memberList = append(memberList, member)
	}
	count := s.sprintMemberDao.Count(req.SprintId)
	return &projectv1.ListSprintMemberResponse{
		Items: memberList,
		Pagination: &generalv1.Pagination{
			Page:  req.Page,
			Size:  req.Size,
			Total: int32(count),
		},
	}, nil
}

func (s *SprintMemberService) Update(ctx context.Context, req *projectv1.UpdateSprintMemberRequest) (*projectv1.UpdateSprintMemberResponse, error) {
	jwt, err := utils.AuthJWT(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var reqProject model.Project
	reqProject.ID = req.ProjectId
	reqProject.CreatedUser.ID = cast.ToInt64(jwt.Id)
	myMember, err := s.projectMemberDao.GetByProject(reqProject)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	if !myMember.IsAdmin || !jwt.IsSuperAdmin {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	var memberList []*model.SprintMember
	for _, m := range req.Members {
		var member = &model.SprintMember{
			Model: model.Model{
				ID: m.Id,
			},
			Role:     m.Role,
			Capacity: m.Capacity,
		}
		memberList = append(memberList, member)
	}
	success, err := s.sprintMemberDao.Update(memberList)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &projectv1.UpdateSprintMemberResponse{
		Success: success,
	}, nil
}

//
//func Create(sprintId int64, userId int64) (int64, error) {
//    sprintMemberId, err := sprintmembermodel.CreateSprintMember(sprintId, userId)
//    if err != nil {
//        return 0, err
//    }
//    return sprintMemberId, nil
//}
//
//func ListSprintMemberBySprint(sprintId int64, page int32, pageSize int32) ([]*userV1.User, int64, error) {
//    sprintMembers, err := sprintmembermodel.ListSprintMemberBySprint(sprintId, page, pageSize)
//    if err != nil {
//        return nil, 0, err
//    }
//    var members []*userV1.User
//    for _, m := range sprintMembers {
//        members = append(members, &userV1.User{
//            Id:           m.User.Id,
//            Name:         m.User.Name,
//            Email:        m.User.Email,
//            IsSuperAdmin: m.User.IsSuperAdmin,
//            RoleName:     m.Role.Name,
//            Capacity:     m.Capacity,
//        })
//    }
//    count := sprintmembermodel.CountSprintMemberBySprint(sprintId)
//    return members, count, nil
//}
//
//func ListSprintByUser(userId int64, page int32, pageSize int32) ([]*projectV1.Sprint, int64, error) {
//    sprintMembers, err := sprintmembermodel.ListSprintMemberByUser(userId, page, pageSize)
//    if err != nil {
//        return nil, 0, err
//    }
//    var sprints []*projectV1.Sprint
//    for _, s := range sprintMembers {
//        sprints = append(sprints, &projectV1.Sprint{
//            Id:        s.Sprint.Id,
//            Name:      s.Sprint.Name,
//            StartDate: s.Sprint.StartDate.Unix(),
//            EndDate:   s.Sprint.EndDate.Unix(),
//        })
//    }
//    count := sprintmembermodel.CountSprintMemberByUser(userId)
//    return sprints, count, nil
//}
