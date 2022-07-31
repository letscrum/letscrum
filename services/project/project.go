package projectService

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
	for _, project := range projects {
		list = append(list, &projectV1.Project{
			Id:          project.ID,
			Name:        project.Name,
			DisplayName: project.DisplayName,
		})
	}
	count := models.CountProject()
	return list, count, nil
}
