package models

import (
	generalV1 "github.com/letscrum/letscrum/apis/general/v1"
	projectV1 "github.com/letscrum/letscrum/apis/project/v1"
	"gorm.io/gorm"
)

type Project struct {
	Model

	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

func CreateProject(project *projectV1.Project) error {
	p := Project{
		Name:        project.Name,
		DisplayName: project.DisplayName,
	}
	if err := db.Create(&p).Error; err != nil {
		return err
	}
	return nil
}

func ListProject(pagination *generalV1.Pagination) ([]*Project, error) {
	var projects []*Project
	err := db.Limit(int(pagination.PageSize)).Offset(int((pagination.Page - 1) * pagination.PageSize)).Find(&projects).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return projects, nil
}
