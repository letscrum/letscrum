package project

import (
	generalV1 "github.com/letscrum/letscrum/apis/general/v1"
	projectV1 "github.com/letscrum/letscrum/apis/project/v1"
	userV1 "github.com/letscrum/letscrum/apis/user/v1"
	projectModel "github.com/letscrum/letscrum/models/project"
)

func Create(project *projectV1.Project) error {
	if err := projectModel.CreateProject(project.Name, project.DisplayName, project.CreatedUser.Id); err != nil {
		return err
	}
	return nil
}

func List(pagination *generalV1.Pagination) ([]*projectV1.Project, int64, error) {
	projects, err := projectModel.ListProject(pagination.Page, pagination.PageSize)
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
	count := projectModel.CountProject()
	return list, count, nil
}

func Update(project *projectV1.Project) error {
	if err := projectModel.UpdateProject(project.Name, project.DisplayName); err != nil {
		return err
	}
	return nil
}

func Delete(name string) error {
	if err := projectModel.DeleteProject(name); err != nil {
		return err
	}
	return nil
}

func Get(name string) (*projectV1.Project, error) {
	p, err := projectModel.GetProject(name)
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
