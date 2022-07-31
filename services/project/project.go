package project

import (
	generalV1 "github.com/letscrum/letscrum/apis/general/v1"
	projectV1 "github.com/letscrum/letscrum/apis/project/v1"
	"github.com/letscrum/letscrum/models"
)

func Create(project *projectV1.Project) error {
	if err := models.CreateProject(project); err != nil {
		return err
	}
	return nil
}

func List(pagination *generalV1.Pagination) ([]*projectV1.Project, int64, error) {
	projects, err := models.ListProject(pagination)
	if err != nil {
		return nil, 0, err
	}
	var list []*projectV1.Project
	for _, p := range projects {
		list = append(list, &projectV1.Project{
			Id:          p.Id,
			Name:        p.Name,
			DisplayName: p.DisplayName,
			CreatedAt:   p.CreatedAt.Unix(),
			UpdatedAt:   p.UpdatedAt.Unix(),
		})
	}
	count := models.CountProject()
	return list, count, nil
}

func Update(project *projectV1.Project) error {
	if err := models.UpdateProject(project.Name, project); err != nil {
		return err
	}
	return nil
}

func Delete(name string) error {
	if err := models.DeleteProject(name); err != nil {
		return err
	}
	return nil
}

func Get(name string) (*projectV1.Project, error) {
	p, err := models.GetProject(name)
	if err != nil {
		return nil, err
	}
	project := &projectV1.Project{
		Id:          p.Id,
		Name:        p.Name,
		DisplayName: p.DisplayName,
		CreatedAt:   p.CreatedAt.Unix(),
		UpdatedAt:   p.UpdatedAt.Unix(),
	}
	return project, nil
}
