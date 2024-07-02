package v1

import (
	"context"
	"encoding/json"
	"time"

	generalv1 "github.com/letscrum/letscrum/api/general/v1"
	v1 "github.com/letscrum/letscrum/api/letscrum/v1"
	projectv1 "github.com/letscrum/letscrum/api/project/v1"
	userV1 "github.com/letscrum/letscrum/api/user/v1"
	"github.com/letscrum/letscrum/internal/dao"
	"github.com/letscrum/letscrum/internal/model"
	"github.com/letscrum/letscrum/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProjectService struct {
	v1.UnimplementedProjectServer
	projectDao dao.ProjectDao
	userDao    dao.UserDao
	sprintDao  dao.SprintDao
}

func NewProjectService(dao dao.Interface) *ProjectService {
	return &ProjectService{
		projectDao: dao.ProjectDao(),
		userDao:    dao.UserDao(),
		sprintDao:  dao.SprintDao(),
	}
}

func (s ProjectService) Get(ctx context.Context, req *projectv1.GetProjectRequest) (*projectv1.GetProjectResponse, error) {
	var reqProject model.Project
	reqProject.ID = req.ProjectId
	project, err := s.projectDao.Get(reqProject)
	if err != nil {
		result := status.Convert(err)
		if result.Code() == codes.NotFound {
			return nil, status.Errorf(codes.NotFound, "get book err: %d not found", req.ProjectId)
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}
	if project.ID == 0 {
		return nil, status.Error(codes.NotFound, "project not fount.")
	}
	sprints, err := s.sprintDao.ListByProject(reqProject, 1, 999, "")
	if err != nil {
		result := status.Convert(err)
		if result.Code() == codes.NotFound {
			return nil, status.Errorf(codes.NotFound, "not found.")
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}
	sprint := projectv1.Sprint{}
	for i, s := range sprints {
		// Set the 1st Sprint as default
		if i == 0 {
			sprint = projectv1.Sprint{
				Id:        s.ID,
				ProjectId: s.ProjectID,
				Name:      s.Name,
				StartDate: s.StartDate.Unix(),
				EndDate:   s.EndDate.Unix(),
				Status:    projectv1.Sprint_UNKNOWN,
				CreatedAt: s.CreatedAt.Unix(),
				UpdatedAt: s.UpdatedAt.Unix(),
			}
		}
		// Set the real current sprint
		if time.Now().After(s.StartDate) && time.Now().Before(s.EndDate) {
			sprint = projectv1.Sprint{
				Id:        s.ID,
				ProjectId: s.ProjectID,
				Name:      s.Name,
				StartDate: s.StartDate.Unix(),
				EndDate:   s.EndDate.Unix(),
				Status:    projectv1.Sprint_Current,
				CreatedAt: s.CreatedAt.Unix(),
				UpdatedAt: s.UpdatedAt.Unix(),
			}
			break
		}
	}
	var members []*projectv1.ProjectMember
	err = json.Unmarshal([]byte(project.Members), &members)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &projectv1.GetProjectResponse{
		Item: &projectv1.Project{
			Id:          project.ID,
			Name:        project.Name,
			DisplayName: project.DisplayName,
			Description: project.Description,
			CreatedUser: &userV1.User{
				Id:           project.CreatedUser.ID,
				Name:         project.CreatedUser.Name,
				IsSuperAdmin: project.CreatedUser.IsSuperAdmin,
			},
			Members:       members,
			CurrentSprint: &sprint,
			CreatedAt:     project.CreatedAt.Unix(),
			UpdatedAt:     project.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *ProjectService) List(ctx context.Context, req *projectv1.ListProjectRequest) (*projectv1.ListProjectResponse, error) {
	req.Page, req.Size = utils.Pagination(req.Page, req.Size)
	projects, err := s.projectDao.List(req.Page, req.Size, req.Keyword)
	if err != nil {
		result := status.Convert(err)
		if result.Code() == codes.NotFound {
			return nil, status.Errorf(codes.NotFound, "not found.")
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}
	var list []*projectv1.Project
	for _, p := range projects {
		var members []*projectv1.ProjectMember
		if p.Members != "" {
			err = json.Unmarshal([]byte(p.Members), &members)
			if err != nil {
				return nil, status.Error(codes.Unknown, err.Error())
			}
		}
		var project = &projectv1.Project{
			Id:          p.ID,
			Name:        p.Name,
			DisplayName: p.DisplayName,
			Description: p.Description,
			Members:     members,
			CreatedUser: &userV1.User{
				Id:           p.CreatedUser.ID,
				Name:         p.CreatedUser.Name,
				IsSuperAdmin: p.CreatedUser.IsSuperAdmin,
			},
			CreatedAt: p.CreatedAt.Unix(),
			UpdatedAt: p.UpdatedAt.Unix(),
		}
		list = append(list, project)
	}
	count := s.projectDao.Count(req.Keyword)
	return &projectv1.ListProjectResponse{
		Items: list,
		Pagination: &generalv1.Pagination{
			Page:  req.Page,
			Size:  req.Size,
			Total: int32(count),
		},
	}, nil
}

func (s *ProjectService) Create(ctx context.Context, req *projectv1.CreateProjectRequest) (*projectv1.CreateProjectResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	// if claims.IsSuperAdmin is false permission dene
	if !claims.IsSuperAdmin {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	if req.DisplayName == "" {
		return nil, status.Error(codes.InvalidArgument, "project display name can't be empty.")
	}
	// convert req.Members to project Members
	var members []*projectv1.ProjectMember
	members = append(members, &projectv1.ProjectMember{
		UserId:   int64(claims.ID),
		UserName: claims.Name,
		IsAdmin:  true,
	})
	users, err := s.userDao.ListByIds(1, 999, req.Members)
	for _, u := range users {
		isAdmin := false
		if u.IsSuperAdmin == true {
			isAdmin = true
		}
		member := &projectv1.ProjectMember{
			UserId:   u.ID,
			UserName: u.Name,
			IsAdmin:  isAdmin,
		}
		members = append(members, member)
	}
	// convert members to json string
	membersJson, err := json.Marshal(members)
	newProject := model.Project{
		Name:        req.DisplayName,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Members:     string(membersJson),
		CreatedBy:   int64(claims.ID),
	}
	project, err := s.projectDao.Create(newProject)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	success := false
	if project.ID > 0 {
		success = true
	}
	return &projectv1.CreateProjectResponse{
		Success: success,
		Id:      project.ID,
	}, nil
}

func (s *ProjectService) Update(ctx context.Context, req *projectv1.UpdateProjectRequest) (*projectv1.UpdateProjectResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	var reqProject model.Project
	reqProject.ID = req.ProjectId
	reqProject.CreatedUser.ID = int64(claims.ID)
	project, err := s.projectDao.Get(reqProject)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	var projectMembers []*projectv1.ProjectMember
	err = json.Unmarshal([]byte(project.Members), &projectMembers)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	for _, m := range projectMembers {
		if m.UserId == int64(claims.ID) && m.IsAdmin == false {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
	}
	project.DisplayName = req.DisplayName
	project.Description = req.Description
	// convert req.Members to project Members
	var members []*projectv1.ProjectMember
	users, err := s.userDao.ListByIds(1, 999, req.Members)
	for _, u := range users {
		isAdmin := false
		if u.IsSuperAdmin == true {
			isAdmin = true
		}
		member := &projectv1.ProjectMember{
			UserId:   u.ID,
			UserName: u.Name,
			IsAdmin:  isAdmin,
		}
		members = append(members, member)
	}
	membersJson, err := json.Marshal(members)
	project.Members = string(membersJson)
	updatedProject, err := s.projectDao.Update(*project)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &projectv1.UpdateProjectResponse{
		Success: updatedProject != nil,
		Id:      updatedProject.ID,
	}, nil
}

func (s *ProjectService) Delete(ctx context.Context, req *projectv1.DeleteProjectRequest) (*projectv1.DeleteProjectResponse, error) {
	claims, err := utils.GetTokenDetails(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	user, err := s.userDao.Get(int64(claims.ID))
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if !user.IsSuperAdmin {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	var reqProject model.Project
	reqProject.ID = req.ProjectId
	deletedProject, err := s.projectDao.Delete(reqProject)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &projectv1.DeleteProjectResponse{
		Success: deletedProject != nil,
		Id:      deletedProject.ID,
	}, nil
}

//
//func Create1(name string, displayName string, createdUserId int64) (int64, error) {
//    projectId, err := model.CreateProject(name, displayName, createdUserId)
//    if err != nil {
//        return 0, err
//    }
//    _, errCreateMember := models2.CreateProjectMember(projectId, createdUserId, true)
//    if errCreateMember != nil {
//        errDeleteProject := model.DeleteProject(projectId)
//        if errDeleteProject != nil {
//            return 0, errDeleteProject
//        }
//        return 0, errCreateMember
//    }
//    return projectId, nil
//}
//
//func List1(page int32, pageSize int32) ([]*projectv1.Project, int64, error) {
//    projects, err := model.ListProject(page, pageSize)
//    if err != nil {
//        return nil, 0, err
//    }
//    var list []*projectv1.Project
//    for _, p := range projects {
//        list = append(list, &projectv1.Project{
//            Id:          p.Id,
//            Name:        p.Name,
//            DisplayName: p.DisplayName,
//            CreatedUser: &userV1.User{
//                Id:   p.CreatedUser.Id,
//                Name: p.CreatedUser.Name,
//            },
//            CreatedAt: p.CreatedAt.Unix(),
//            UpdatedAt: p.UpdatedAt.Unix(),
//        })
//    }
//    count := model.CountProject()
//    return list, count, nil
//}
//
//func Update1(id int64, displayName string) error {
//    if err := model.UpdateProject(id, displayName); err != nil {
//        return err
//    }
//    return nil
//}
//
//func Delete1(id int64) error {
//    if err := model.DeleteProject(id); err != nil {
//        return err
//    }
//    return nil
//}
//
//func HardDelete1(id int64) error {
//    if err := models2.HardDeleteProjectMemberByProject(id); err != nil {
//        return err
//    }
//    if err := model.HardDeleteProject(id); err != nil {
//        return err
//    }
//    return nil
//}
