package models

import v1 "github.com/letscrum/letscrum/apis/project/v1"

type Project struct {
	Model

	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

func CreateProject(project *v1.Project) error {
	p := Project{
		Name:        project.Name,
		DisplayName: project.DisplayName,
	}
	if err := db.Create(&p).Error; err != nil {
		return err
	}
	return nil
}

func ListProject() error {

}
