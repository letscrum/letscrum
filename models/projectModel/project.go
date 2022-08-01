package projectModel

import (
	"github.com/letscrum/letscrum/models"
	userModel "github.com/letscrum/letscrum/models/user"
)

type Project struct {
	models.Model

	Name        string         `json:"name"`
	DisplayName string         `json:"display_name"`
	CreatedBy   int64          `json:"created_by"`
	CreatedUser userModel.User `gorm:"foreignKey:CreatedBy"`
}

type ProjectMember struct {
	models.Model

	ProjectId int64          `json:"project_id"`
	UserId    int64          `json:"user_id"`
	IsAdmin   bool           `json:"is_admin"`
	User      userModel.User `gorm:"foreignKey:UserId"`
	Project   userModel.User `gorm:"foreignKey:ProjectId"`
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

	if err := models.DB.Create(&p).Error; err != nil {
		return 0, err
	}
	return p.Id, nil
}

func ListProject(page int32, pageSize int32) ([]*Project, error) {
	var projects []*Project
	err := models.DB.Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Preload("CreatedUser").Find(&projects).Error
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func CountProject() int64 {
	count := int64(0)
	models.DB.Model(&Project{}).Count(&count)
	return count
}

func UpdateProject(id int64, displayName string) error {
	if err := models.DB.Model(&Project{}).Where("id = ?", id).Update("display_name", displayName).Error; err != nil {
		return err
	}
	return nil
}

func DeleteProject(id int64) error {
	if err := models.DB.Where("id = ?", id).Delete(&Project{}).Error; err != nil {
		return err
	}
	return nil
}

func GetProject(id int64) (*Project, error) {
	var p *Project
	if err := models.DB.Where("id = ?", id).Preload("CreatedUser").Find(&p).Error; err != nil {
		return nil, err
	}
	return p, nil
}

func CreateProjectMember(projectId int64, userId int64, isAdmin bool) error {
	pm := ProjectMember{
		ProjectId: projectId,
		UserId:    userId,
		IsAdmin:   isAdmin,
	}
	if err := models.DB.Create(&pm).Error; err != nil {
		return err
	}
	return nil
}

func UpdateProjectMember(projectId int64, userId int64, isAdmin bool) error {
	if err := models.DB.Model(&ProjectMember{}).Where("project_id = ?", projectId).Where("user_id = ?", userId).Update("is_admin", isAdmin).Error; err != nil {
		return err
	}
	return nil
}

func DeleteProjectMember(projectId int64, userId int64) error {
	if err := models.DB.Where("project_id = ?", projectId).Where("user_id = ?", userId).Delete(&Project{}).Error; err != nil {
		return err
	}
	return nil
}

func ListProjectMemberByProject(projectId int64, page int32, pageSize int32) ([]*ProjectMember, error) {
	var projectMembers []*ProjectMember
	err := models.DB.Where("project_id", projectId).Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Preload("User").Find(&projectMembers).Error
	if err != nil {
		return nil, err
	}
	return projectMembers, nil
}

func ListProjectMemberByUser(userId int64, page int32, pageSize int32) ([]*ProjectMember, error) {
	var projectMembers []*ProjectMember
	err := models.DB.Where("user_id", userId).Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Preload("Project").Find(&projectMembers).Error
	if err != nil {
		return nil, err
	}
	return projectMembers, nil
}
