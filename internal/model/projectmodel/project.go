package projectmodel

import (
	"github.com/letscrum/letscrum/internal/model"
	"github.com/letscrum/letscrum/internal/model/usermodel"
)

type Project struct {
	model.Model

	Name        string         `json:"name"`
	DisplayName string         `json:"display_name"`
	CreatedBy   int64          `json:"created_by"`
	CreatedUser usermodel.User `gorm:"foreignKey:CreatedBy"`
}

func CreateProject(name string, displayName string, createdUserId int64) (int64, error) {
	p := Project{
		Name:        name,
		DisplayName: displayName,
		CreatedBy:   createdUserId,
	}

	//var pInDB *Project
	//errPName := db.Where("name = ?", p.Name).First(&pInDB).Error
	//if errPName != nil && errPName != gorm.ErrRecordNotFound {
	//	return errPName
	//}
	//if pInDB != nil && pInDB.Name == p.Name {
	//	return fmt.Errorf("duplicate project name")
	//}

	if err := model.DB.Create(&p).Error; err != nil {
		return 0, err
	}
	return p.Id, nil
}

func ListProject(page int32, pageSize int32) ([]*Project, error) {
	var projects []*Project
	err := model.DB.Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Preload("CreatedUser").Find(&projects).Error
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func CountProject() int64 {
	count := int64(0)
	model.DB.Model(&Project{}).Count(&count)
	return count
}

func UpdateProject(id int64, displayName string) error {
	if err := model.DB.Model(&Project{}).Where("id = ?", id).Update("display_name", displayName).Error; err != nil {
		return err
	}
	return nil
}

func DeleteProject(id int64) error {
	if err := model.DB.Where("id = ?", id).Delete(&Project{}).Error; err != nil {
		return err
	}
	return nil
}

func HardDeleteProject(id int64) error {
	if err := model.DB.Unscoped().Where("id = ?", id).Delete(&Project{}).Error; err != nil {
		return err
	}
	return nil
}

func GetProject(id int64) (*Project, error) {
	var p *Project
	if err := model.DB.Where("id = ?", id).Preload("CreatedUser").Find(&p).Error; err != nil {
		return nil, err
	}
	return p, nil
}
