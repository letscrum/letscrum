package v1

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	generalv1 "github.com/letscrum/letscrum/api/general/v1"
	v1 "github.com/letscrum/letscrum/api/letscrum/v1"
	projectv1 "github.com/letscrum/letscrum/api/project/v1"
	"github.com/letscrum/letscrum/internal/dao"
	"github.com/letscrum/letscrum/internal/model"
	"github.com/letscrum/letscrum/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SprintService struct {
	v1.UnimplementedSprintServer
	sprintDao  dao.SprintDao
	projectDao dao.ProjectDao
}

func NewSprintService(dao dao.Interface) *SprintService {
	return &SprintService{
		sprintDao:  dao.SprintDao(),
		projectDao: dao.ProjectDao(),
	}
}

func (s *SprintService) Create(ctx context.Context, req *projectv1.CreateSprintRequest) (*projectv1.CreateSprintResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var user model.User
	user.Id = claims.Id
	user.IsSuperAdmin = claims.IsSuperAdmin
	var reqProject model.Project
	oId, err := uuid.Parse(req.OrgId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	pId, err := uuid.Parse(req.ProjectId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	reqProject.OrgId = oId
	reqProject.Id = pId
	project, err := s.projectDao.Get(reqProject)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if utils.IsProjectAdmin(*project, user) == false {
		return nil, status.Error(codes.PermissionDenied, utils.ErrNotProjectAdmin)
	}
	var projectMembers []*projectv1.ProjectMember
	err = json.Unmarshal([]byte(project.Members), &projectMembers)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	// add sprint members from project members
	var sprintMembers []*projectv1.SprintMember
	for _, m := range projectMembers {
		var member = &projectv1.SprintMember{
			UserId:   m.UserId,
			UserName: m.UserName,
			Capacity: 0,
			Role:     projectv1.SprintMember_Unassigned,
		}
		sprintMembers = append(sprintMembers, member)
	}
	members, err := json.Marshal(sprintMembers)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var newSprint model.Sprint
	newSprint.Id = uuid.New()
	newSprint.ProjectId = project.Id
	newSprint.Name = req.Name
	newSprint.Members = string(members)
	newSprint.StartDate = time.Unix(req.StartDate, 0)
	newSprint.EndDate = time.Unix(req.EndDate, 0)

	sprint, err := s.sprintDao.Create(newSprint)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &projectv1.CreateSprintResponse{
		Success: sprint.Id != uuid.Nil,
		Item: &projectv1.Sprint{
			Id:        sprint.Id.String(),
			ProjectId: sprint.ProjectId.String(),
			Name:      sprint.Name,
			StartDate: sprint.StartDate.Unix(),
			EndDate:   sprint.EndDate.Unix(),
			CreatedAt: sprint.CreatedAt.Unix(),
			UpdatedAt: sprint.UpdatedAt.Unix(),
			Members:   sprintMembers,
		},
	}, nil
}

func (s *SprintService) Get(ctx context.Context, req *projectv1.GetSprintRequest) (*projectv1.GetSprintResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var user model.User
	user.Id = claims.Id
	user.IsSuperAdmin = claims.IsSuperAdmin
	var reqProject model.Project
	oId, err := uuid.Parse(req.OrgId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	pId, err := uuid.Parse(req.ProjectId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	reqProject.OrgId = oId
	reqProject.Id = pId
	project, err := s.projectDao.Get(reqProject)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if utils.IsProjectMember(*project, user) == false {
		return nil, status.Error(codes.PermissionDenied, utils.ErrNotProjectMember)
	}
	var reqSprint model.Sprint
	reqSprint.ProjectId = project.Id
	sId, err := uuid.Parse(req.SprintId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	reqSprint.Id = sId
	sprint, err := s.sprintDao.Get(reqSprint)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var sprintMembers []*projectv1.SprintMember
	err = json.Unmarshal([]byte(sprint.Members), &sprintMembers)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	var sprintStatus projectv1.Sprint_SprintStatus
	switch {
	case time.Now().After(sprint.StartDate) && time.Now().Before(sprint.EndDate):
		sprintStatus = projectv1.Sprint_Current
		break
	case time.Now().After(sprint.EndDate):
		sprintStatus = projectv1.Sprint_Past
		break
	case time.Now().Before(sprint.StartDate):
		sprintStatus = projectv1.Sprint_Future
		break
	}
	return &projectv1.GetSprintResponse{
		Item: &projectv1.Sprint{
			Id:        sprint.Id.String(),
			ProjectId: sprint.ProjectId.String(),
			Name:      sprint.Name,
			StartDate: sprint.StartDate.Unix(),
			EndDate:   sprint.EndDate.Unix(),
			Status:    sprintStatus,
			CreatedAt: sprint.CreatedAt.Unix(),
			UpdatedAt: sprint.UpdatedAt.Unix(),
			Members:   sprintMembers,
		},
	}, nil
}

func (s *SprintService) List(ctx context.Context, req *projectv1.ListSprintRequest) (*projectv1.ListSprintResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var user model.User
	user.Id = claims.Id
	user.IsSuperAdmin = claims.IsSuperAdmin
	var reqProject model.Project
	oId, err := uuid.Parse(req.OrgId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	pId, err := uuid.Parse(req.ProjectId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	reqProject.OrgId = oId
	reqProject.Id = pId
	project, err := s.projectDao.Get(reqProject)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if utils.IsProjectMember(*project, user) == false {
		return nil, status.Error(codes.PermissionDenied, utils.ErrNotProjectMember)
	}
	sprints, err := s.sprintDao.ListByProject(reqProject, req.Page, req.Size, req.Keyword)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var list []*projectv1.Sprint
	hasCurrent := false
	for _, sprint := range sprints {
		var sprintStatus projectv1.Sprint_SprintStatus
		switch {
		case time.Now().After(sprint.StartDate) && time.Now().Before(sprint.EndDate) && !hasCurrent:
			sprintStatus = projectv1.Sprint_Current
			hasCurrent = true
			break
		case time.Now().After(sprint.StartDate) && time.Now().Before(sprint.EndDate) && hasCurrent:
			sprintStatus = projectv1.Sprint_Future
			break
		case time.Now().After(sprint.EndDate):
			sprintStatus = projectv1.Sprint_Past
			break
		case time.Now().Before(sprint.StartDate):
			sprintStatus = projectv1.Sprint_Future
			break
		}
		var sprintMembers []*projectv1.SprintMember
		err = json.Unmarshal([]byte(sprint.Members), &sprintMembers)
		if err != nil {
			println(sprint.Members)
			return nil, status.Error(codes.Unknown, err.Error())
		}
		var currentSprint = &projectv1.Sprint{
			Id:        sprint.Id.String(),
			ProjectId: sprint.ProjectId.String(),
			Name:      sprint.Name,
			StartDate: sprint.StartDate.Unix(),
			EndDate:   sprint.EndDate.Unix(),
			Status:    sprintStatus,
			CreatedAt: sprint.CreatedAt.Unix(),
			UpdatedAt: sprint.UpdatedAt.Unix(),
			Members:   sprintMembers,
		}
		list = append(list, currentSprint)
	}
	count := s.sprintDao.CountByProject(reqProject, req.Keyword)
	return &projectv1.ListSprintResponse{
		Items: list,
		Pagination: &generalv1.Pagination{
			Page:  req.Page,
			Size:  req.Size,
			Total: int32(count),
		},
	}, nil
}

func (s *SprintService) Update(ctx context.Context, req *projectv1.UpdateSprintRequest) (*projectv1.UpdateSprintResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var user model.User
	user.Id = claims.Id
	user.IsSuperAdmin = claims.IsSuperAdmin
	var reqProject model.Project
	oId, err := uuid.Parse(req.OrgId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	pId, err := uuid.Parse(req.ProjectId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	reqProject.OrgId = oId
	reqProject.Id = pId
	project, err := s.projectDao.Get(reqProject)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if utils.IsProjectAdmin(*project, user) == false {
		return nil, status.Error(codes.PermissionDenied, utils.ErrNotProjectAdmin)
	}
	var sprint model.Sprint
	sId, err := uuid.Parse(req.SprintId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	sprint.Id = sId
	sprint.ProjectId = project.Id
	sprint.Name = req.Name
	sprint.StartDate = time.Unix(req.StartDate, 0)
	sprint.EndDate = time.Unix(req.EndDate, 0)
	updateSprint, err := s.sprintDao.Update(sprint)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var sprintMembers []*projectv1.SprintMember
	err = json.Unmarshal([]byte(updateSprint.Members), &sprintMembers)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &projectv1.UpdateSprintResponse{
		Success: updateSprint.Id != uuid.Nil,
		Item: &projectv1.Sprint{
			Id:        updateSprint.Id.String(),
			ProjectId: updateSprint.ProjectId.String(),
			Name:      updateSprint.Name,
			StartDate: updateSprint.StartDate.Unix(),
			EndDate:   updateSprint.EndDate.Unix(),
			CreatedAt: updateSprint.CreatedAt.Unix(),
			UpdatedAt: updateSprint.UpdatedAt.Unix(),
			Members:   sprintMembers,
		},
	}, nil
}

func (s *SprintService) UpdateMembers(ctx context.Context, req *projectv1.UpdateSprintMembersRequest) (*projectv1.UpdateSprintResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var user model.User
	user.Id = claims.Id
	user.IsSuperAdmin = claims.IsSuperAdmin
	var reqProject model.Project
	oId, err := uuid.Parse(req.OrgId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	pId, err := uuid.Parse(req.ProjectId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	reqProject.OrgId = oId
	reqProject.Id = pId
	project, err := s.projectDao.Get(reqProject)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if utils.IsProjectAdmin(*project, user) == false {
		return nil, status.Error(codes.PermissionDenied, utils.ErrNotProjectAdmin)
	}
	var sprint model.Sprint
	sId, err := uuid.Parse(req.SprintId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	sprint.Id = sId
	var sprintMembers []*projectv1.SprintMember
	for _, m := range req.Members {
		var member = &projectv1.SprintMember{
			UserId:   m.UserId,
			UserName: m.UserName,
			Capacity: m.Capacity,
			Role:     m.Role,
		}
		sprintMembers = append(sprintMembers, member)
	}
	members, err := json.Marshal(sprintMembers)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	sprint.Members = string(members)
	updateSprint, err := s.sprintDao.UpdateMembers(sprint)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &projectv1.UpdateSprintResponse{
		Success: updateSprint.Id != uuid.Nil,
		Item: &projectv1.Sprint{
			Id:        updateSprint.Id.String(),
			ProjectId: updateSprint.ProjectId.String(),
			Name:      updateSprint.Name,
			StartDate: updateSprint.StartDate.Unix(),
			EndDate:   updateSprint.EndDate.Unix(),
			CreatedAt: updateSprint.CreatedAt.Unix(),
			UpdatedAt: updateSprint.UpdatedAt.Unix(),
			Members:   sprintMembers,
		},
	}, nil
}

func (s *SprintService) Delete(ctx context.Context, req *projectv1.DeleteSprintRequest) (*projectv1.DeleteSprintResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var user model.User
	user.Id = claims.Id
	user.IsSuperAdmin = claims.IsSuperAdmin
	var reqProject model.Project
	oId, err := uuid.Parse(req.OrgId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	pId, err := uuid.Parse(req.ProjectId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	reqProject.OrgId = oId
	reqProject.Id = pId
	project, err := s.projectDao.Get(reqProject)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if utils.IsProjectAdmin(*project, user) == false {
		return nil, status.Error(codes.PermissionDenied, utils.ErrNotProjectAdmin)
	}
	var sprint model.Sprint
	sId, err := uuid.Parse(req.SprintId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	sprint.Id = sId
	deleteSprint, err := s.sprintDao.Delete(sprint)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &projectv1.DeleteSprintResponse{
		Success: deleteSprint,
	}, nil
}

func (s SprintService) AddMember(ctx context.Context, req *projectv1.AddSprintMemberRequest) (*projectv1.UpdateSprintResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var user model.User
	user.Id = claims.Id
	user.IsSuperAdmin = claims.IsSuperAdmin
	var reqProject model.Project
	oId, err := uuid.Parse(req.OrgId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	pId, err := uuid.Parse(req.ProjectId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	reqProject.OrgId = oId
	reqProject.Id = pId
	project, err := s.projectDao.Get(reqProject)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if utils.IsProjectAdmin(*project, user) == false {
		return nil, status.Error(codes.PermissionDenied, utils.ErrNotProjectAdmin)
	}
	var sprint model.Sprint
	sId, err := uuid.Parse(req.SprintId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	sprint.Id = sId
	sprint.ProjectId = pId
	getSprint, err := s.sprintDao.Get(sprint)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var sprintMembers []*projectv1.SprintMember
	err = json.Unmarshal([]byte(getSprint.Members), &sprintMembers)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	success := false
	var member model.User
	member.Id = uuid.MustParse(req.Member.UserId)
	member.Name = req.Member.UserName
	if !utils.IsSprintMember(sprintMembers, member) {
		sprintMembers = append(sprintMembers, &projectv1.SprintMember{
			UserId:   req.Member.UserId,
			UserName: req.Member.UserName,
			Capacity: req.Member.Capacity,
			Role:     req.Member.Role,
		})
		members, err := json.Marshal(sprintMembers)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		getSprint.Members = string(members)
		updateSprint, err := s.sprintDao.UpdateMembers(*getSprint)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		if updateSprint.Id != uuid.Nil {
			success = true
		}
	}
	return &projectv1.UpdateSprintResponse{
		Success: success,
		Item: &projectv1.Sprint{
			Id:        getSprint.Id.String(),
			ProjectId: getSprint.ProjectId.String(),
			Name:      getSprint.Name,
			StartDate: getSprint.StartDate.Unix(),
			EndDate:   getSprint.EndDate.Unix(),
			CreatedAt: getSprint.CreatedAt.Unix(),
			UpdatedAt: getSprint.UpdatedAt.Unix(),
			Members:   sprintMembers,
		},
	}, nil
}

func (s SprintService) RemoveMember(ctx context.Context, req *projectv1.RemoveSprintMemberRequest) (*projectv1.UpdateSprintResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var user model.User
	user.Id = claims.Id
	user.IsSuperAdmin = claims.IsSuperAdmin
	var reqProject model.Project
	oId, err := uuid.Parse(req.OrgId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	pId, err := uuid.Parse(req.ProjectId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	reqProject.OrgId = oId
	reqProject.Id = pId
	project, err := s.projectDao.Get(reqProject)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if utils.IsProjectAdmin(*project, user) == false {
		return nil, status.Error(codes.PermissionDenied, utils.ErrNotProjectAdmin)
	}
	var sprint model.Sprint
	sId, err := uuid.Parse(req.SprintId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	sprint.Id = sId
	sprint.ProjectId = pId
	getSprint, err := s.sprintDao.Get(sprint)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var sprintMembers []*projectv1.SprintMember
	err = json.Unmarshal([]byte(getSprint.Members), &sprintMembers)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	success := false
	var member model.User
	member.Id = uuid.MustParse(req.UserId)
	if utils.IsSprintMember(sprintMembers, member) {
		// from sprintMembers remove the item if userId is req.UserId don't use additional memory
		for i, m := range sprintMembers {
			if m.UserId == req.UserId {
				sprintMembers = append(sprintMembers[:i], sprintMembers[i+1:]...)
				break
			}
		}

		members, err := json.Marshal(sprintMembers)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		getSprint.Members = string(members)
		updateSprint, err := s.sprintDao.UpdateMembers(*getSprint)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		if updateSprint.Id != uuid.Nil {
			success = true
		}
	}
	return &projectv1.UpdateSprintResponse{
		Success: success,
		Item: &projectv1.Sprint{
			Id:        getSprint.Id.String(),
			ProjectId: getSprint.ProjectId.String(),
			Name:      getSprint.Name,
			StartDate: getSprint.StartDate.Unix(),
			EndDate:   getSprint.EndDate.Unix(),
			CreatedAt: getSprint.CreatedAt.Unix(),
			UpdatedAt: getSprint.UpdatedAt.Unix(),
			Members:   sprintMembers,
		},
	}, nil
}
