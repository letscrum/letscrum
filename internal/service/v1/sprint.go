package v1

import (
	"context"
	"encoding/json"
	"time"

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
	var reqProject model.Project
	reqProject.Id = req.ProjectId
	reqProject.CreatedUser.Id = int64(claims.Id)
	project, err := s.projectDao.Get(reqProject)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	var projectMembers []*projectv1.ProjectMember
	err = json.Unmarshal([]byte(project.Members), &projectMembers)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	for _, m := range projectMembers {
		if m.UserId == int64(claims.Id) && m.IsAdmin == false {
			return nil, status.Error(codes.PermissionDenied, "No permission.")
		}
	}

	// add sprint members from project members
	var sprintMembers []*projectv1.SprintMember
	for _, m := range projectMembers {
		var member = &projectv1.SprintMember{
			UserId:   m.UserId,
			UserName: m.UserName,
			Capacity: 0,
		}
		sprintMembers = append(sprintMembers, member)
	}
	members, err := json.Marshal(sprintMembers)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	newSprint := model.Sprint{
		ProjectId: req.ProjectId,
		Name:      req.Name,
		Members:   string(members),
		StartDate: time.Unix(req.StartDate, 0),
		EndDate:   time.Unix(req.EndDate, 0),
	}

	sprint, err := s.sprintDao.Create(newSprint)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	success := false
	if sprint.Id > 0 {
		success = true
	}
	return &projectv1.CreateSprintResponse{
		Success: success,
		Id:      sprint.Id,
	}, nil
}

func (s *SprintService) Get(ctx context.Context, req *projectv1.GetSprintRequest) (*projectv1.GetSprintResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var reqProject model.Project
	reqProject.Id = req.ProjectId
	reqProject.CreatedUser.Id = int64(claims.Id)
	project, err := s.projectDao.Get(reqProject)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	var projectMembers []*projectv1.ProjectMember
	err = json.Unmarshal([]byte(project.Members), &projectMembers)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	for _, m := range projectMembers {
		if m.UserId == int64(claims.Id) && m.IsAdmin == false {
			return nil, status.Error(codes.PermissionDenied, "No permission.")
		}
	}
	var reqSprint model.Sprint
	reqSprint.Id = req.SprintId
	sprint, err := s.sprintDao.Get(reqSprint)
	if err != nil {
		result := status.Convert(err)
		if result.Code() == codes.NotFound {
			return nil, status.Errorf(codes.NotFound, "not found.")
		}
		return nil, status.Error(codes.Unknown, err.Error())
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
			Id:        sprint.Id,
			ProjectId: sprint.ProjectId,
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
	var reqProject model.Project
	reqProject.Id = req.ProjectId
	reqProject.CreatedUser.Id = int64(claims.Id)
	project, err := s.projectDao.Get(reqProject)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	var projectMembers []*projectv1.ProjectMember
	err = json.Unmarshal([]byte(project.Members), &projectMembers)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	for _, m := range projectMembers {
		if m.UserId == int64(claims.Id) && m.IsAdmin == false {
			return nil, status.Error(codes.PermissionDenied, "No permission.")
		}
	}
	sprints, err := s.sprintDao.ListByProject(reqProject, req.Page, req.Size, req.Keyword)
	if err != nil {
		result := status.Convert(err)
		if result.Code() == codes.NotFound {
			return nil, status.Errorf(codes.NotFound, "not found.")
		}
		return nil, status.Error(codes.Unknown, err.Error())
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
			Id:        sprint.Id,
			ProjectId: sprint.ProjectId,
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
	var reqProject model.Project
	reqProject.Id = req.ProjectId
	reqProject.CreatedUser.Id = int64(claims.Id)
	project, err := s.projectDao.Get(reqProject)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	var projectMembers []*projectv1.ProjectMember
	err = json.Unmarshal([]byte(project.Members), &projectMembers)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	for _, m := range projectMembers {
		if m.UserId == int64(claims.Id) && m.IsAdmin == false {
			return nil, status.Error(codes.PermissionDenied, "No permission.")
		}
	}
	var sprint model.Sprint
	sprint.Id = req.SprintId
	sprint.ProjectId = req.ProjectId
	sprint.Name = req.Name
	sprint.StartDate = time.Unix(req.StartDate, 0)
	sprint.EndDate = time.Unix(req.EndDate, 0)
	updateSprint, err := s.sprintDao.Update(sprint)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &projectv1.UpdateSprintResponse{
		Success: updateSprint != nil,
		Id:      updateSprint.Id,
	}, nil
}

func (s *SprintService) Delete(ctx context.Context, req *projectv1.DeleteSprintRequest) (*projectv1.DeleteSprintResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var reqProject model.Project
	reqProject.Id = req.ProjectId
	reqProject.CreatedUser.Id = int64(claims.Id)
	project, err := s.projectDao.Get(reqProject)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	var projectMembers []*projectv1.ProjectMember
	err = json.Unmarshal([]byte(project.Members), &projectMembers)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	for _, m := range projectMembers {
		if m.UserId == int64(claims.Id) && m.IsAdmin == false {
			return nil, status.Error(codes.PermissionDenied, "No permission.")
		}
	}
	var sprint model.Sprint
	sprint.Id = req.SprintId
	deleteSprint, err := s.sprintDao.Delete(sprint)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &projectv1.DeleteSprintResponse{
		Success: deleteSprint,
	}, nil

}

//
//func Create(projectId int64, name string, startDate time.Time, endDate time.Time) (int64, error) {
//    sprintId, err := sprintmodel.CreateSprint(projectId, name, startDate, endDate)
//    if err != nil {
//        return 0, err
//    }
//    projectMembers, errGetMembers := models2.ListProjectMemberByProject(projectId, 1, 999)
//    if errGetMembers != nil {
//        errDeleteSprint := sprintmodel.HardDeleteSprint(sprintId)
//        if errDeleteSprint != nil {
//            return 0, errDeleteSprint
//        }
//        return 0, err
//    }
//    var userIds []int64
//    for _, projectMember := range projectMembers {
//        userIds = append(userIds, projectMember.UserId)
//    }
//    _, errCreateSprintMembers := sprintmembermodel.CreateSprintMembers(sprintId, userIds)
//    if errCreateSprintMembers != nil {
//        errDeleteSprint := sprintmodel.HardDeleteSprint(sprintId)
//        if errDeleteSprint != nil {
//            return 0, errDeleteSprint
//        }
//        return 0, err
//    }
//    return sprintId, nil
//}
//
//func List(projectId int64, page int32, pageSize int32) ([]*projectV1.Sprint, int64, error) {
//    sprints, err := sprintmodel.ListSprintByProject(projectId, page, pageSize)
//    if err != nil {
//        return nil, 0, err
//    }
//    var list []*projectV1.Sprint
//    for _, p := range sprints {
//        list = append(list, &projectV1.Sprint{
//            Id:        p.Id,
//            Name:      p.Name,
//            StartDate: p.StartDate.Unix(),
//            EndDate:   p.EndDate.Unix(),
//            CreatedAt: p.CreatedAt.Unix(),
//            UpdatedAt: p.UpdatedAt.Unix(),
//        })
//    }
//    count := sprintmodel.CountSprintByProject(projectId)
//    return list, count, nil
//}
//
//func Update(id int64, name string, startDate time.Time, endDate time.Time) error {
//    if err := sprintmodel.UpdateSprint(id, name, startDate, endDate); err != nil {
//        return err
//    }
//    return nil
//}
//
//func Delete(id int64) error {
//    if err := sprintmodel.DeleteSprint(id); err != nil {
//        return err
//    }
//    return nil
//}
//
//func HardDelete(id int64) error {
//    if err := sprintmembermodel.HardDeleteSprintMemberBySprint(id); err != nil {
//        return err
//    }
//    if err := sprintmodel.HardDeleteSprint(id); err != nil {
//        return err
//    }
//    return nil
//}
//
//func Get(id int64) (*projectV1.Sprint, error) {
//    s, err := sprintmodel.GetSprint(id)
//    if err != nil {
//        return nil, err
//    }
//    sprint := &projectV1.Sprint{
//        Id:        s.Id,
//        Name:      s.Name,
//        StartDate: s.StartDate.Unix(),
//        EndDate:   s.EndDate.Unix(),
//        CreatedAt: s.CreatedAt.Unix(),
//        UpdatedAt: s.UpdatedAt.Unix(),
//    }
//    return sprint, nil
//}
