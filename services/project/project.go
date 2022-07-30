package projectService

import (
	v1 "github.com/letscrum/letscrum/apis/project/v1"
	"github.com/letscrum/letscrum/models"
)

func CreateProject(project *v1.Project) error {
	if err := models.CreateProject(project); err != nil {
		return err
	}

	return nil
}
