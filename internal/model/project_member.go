package model

type ProjectMember struct {
	Model

	ProjectID int64   `gorm:"column:project_id"`
	UserID    int64   `gorm:"column:user_id"`
	IsAdmin   bool    `gorm:"column:is_admin"`
	User      User    `gorm:"foreignKey:UserID"`
	Project   Project `gorm:"foreignKey:ProjectID"`
}

//
//func CreateProjectMember(projectId int64, userId int64, isAdmin bool) (int64, error) {
//	pm := ProjectMember{
//		ProjectId: projectId,
//		UserId:    userId,
//		IsAdmin:   isAdmin,
//	}
//	if err := model.DB.Create(&pm).Error; err != nil {
//		return 0, err
//	}
//	return pm.Id, nil
//}
//
//func UpdateProjectMember(projectId int64, userId int64, isAdmin bool) error {
//	if err := model.DB.Model(&ProjectMember{}).Where("project_id = ?", projectId).Where("user_id = ?", userId).Update("is_admin", isAdmin).Error; err != nil {
//		return err
//	}
//	return nil
//}
//
//func DeleteProjectMember(projectId int64, userId int64) error {
//	if err := model.DB.Where("project_id = ?", projectId).Where("user_id = ?", userId).Delete(&ProjectMember{}).Error; err != nil {
//		return err
//	}
//	return nil
//}
//
//func HardDeleteProjectMember(projectId int64, userId int64) error {
//	if err := model.DB.Unscoped().Where("project_id = ?", projectId).Where("user_id = ?", userId).Delete(&ProjectMember{}).Error; err != nil {
//		return err
//	}
//	return nil
//}
//
//func DeleteProjectMemberByProject(projectId int64) error {
//	if err := model.DB.Where("project_id = ?", projectId).Delete(&ProjectMember{}).Error; err != nil {
//		return err
//	}
//	return nil
//}
//
//func HardDeleteProjectMemberByProject(projectId int64) error {
//	if err := model.DB.Unscoped().Where("project_id = ?", projectId).Delete(&ProjectMember{}).Error; err != nil {
//		return err
//	}
//	return nil
//}
//
//func ListProjectMemberByProject(projectId int64, page int32, pageSize int32) ([]*ProjectMember, error) {
//	var projectMembers []*ProjectMember
//	err := model.DB.Where("project_id", projectId).Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Preload("User").Find(&projectMembers).Error
//	if err != nil {
//		return nil, err
//	}
//	return projectMembers, nil
//}
//
//func CountProjectMemberByProject(projectId int64) int64 {
//	count := int64(0)
//	model.DB.Model(&ProjectMember{}).Where("project_id = ?", projectId).Count(&count)
//	return count
//}
//
//func ListProjectMemberByUser(userId int64, page int32, pageSize int32) ([]*ProjectMember, error) {
//	var projectMembers []*ProjectMember
//	err := model.DB.Where("user_id", userId).Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Preload("Project").Find(&projectMembers).Error
//	if err != nil {
//		return nil, err
//	}
//	return projectMembers, nil
//}
//
//func CountProjectMemberByUser(userId int64) int64 {
//	count := int64(0)
//	model.DB.Model(&ProjectMember{}).Where("user_id = ?", userId).Count(&count)
//	return count
//}
