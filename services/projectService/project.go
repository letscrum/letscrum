package projectService

import (
	projectV1 "github.com/letscrum/letscrum/apis/project/v1"
	userV1 "github.com/letscrum/letscrum/apis/user/v1"
	"github.com/letscrum/letscrum/models"
)

func Create(name string, displayName string, createdUserId int64) (int64, error) {
	projectId, err := models.CreateProject(name, displayName, createdUserId)
	if err != nil {
		return 0, err
	}
	_, errCreateMember := models.CreateProjectMember(projectId, createdUserId, true)
	if errCreateMember != nil {
		errDeleteProject := models.DeleteProject(projectId)
		if errDeleteProject != nil {
			return 0, errDeleteProject
		}
		return 0, errCreateMember
	}
	return projectId, nil
}

func List(page int32, pageSize int32) ([]*projectV1.Project, int64, error) {
	projects, err := models.ListProject(page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	var list []*projectV1.Project
	for _, p := range projects {
		list = append(list, &projectV1.Project{
			Id:          p.Id,
			Name:        p.Name,
			DisplayName: p.DisplayName,
			CreatedUser: &userV1.User{
				Id:   p.CreatedUser.Id,
				Name: p.CreatedUser.Name,
			},
			CreatedAt: p.CreatedAt.Unix(),
			UpdatedAt: p.UpdatedAt.Unix(),
		})
	}
	count := models.CountProject()
	return list, count, nil
}

func Update(id int64, displayName string) error {
	if err := models.UpdateProject(id, displayName); err != nil {
		return err
	}
	return nil
}

func Delete(id int64) error {
	if err := models.DeleteProject(id); err != nil {
		return err
	}
	return nil
}

func HardDelete(id int64) error {
	if err := models.HardDeleteProjectMemberByProject(id); err != nil {
		return err
	}
	if err := models.HardDeleteProject(id); err != nil {
		return err
	}
	return nil
}

func Get(id int64) (*projectV1.Project, error) {
	p, err := models.GetProject(id)
	if err != nil {
		return nil, err
	}
	project := &projectV1.Project{
		Id:          p.Id,
		Name:        p.Name,
		DisplayName: p.DisplayName,
		CreatedUser: &userV1.User{
			Id:   p.CreatedUser.Id,
			Name: p.CreatedUser.Name,
		},
		CreatedAt: p.CreatedAt.Unix(),
		UpdatedAt: p.UpdatedAt.Unix(),
	}
	return project, nil
}
