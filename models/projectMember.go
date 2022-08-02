package models

func CreateProjectMember(projectId int64, userId int64, isAdmin bool) (int64, error) {
	pm := ProjectMember{
		ProjectId: projectId,
		UserId:    userId,
		IsAdmin:   isAdmin,
	}
	if err := DB.Create(&pm).Error; err != nil {
		return 0, err
	}
	return pm.Id, nil
}

func UpdateProjectMember(projectId int64, userId int64, isAdmin bool) error {
	if err := DB.Model(&ProjectMember{}).Where("project_id = ?", projectId).Where("user_id = ?", userId).Update("is_admin", isAdmin).Error; err != nil {
		return err
	}
	return nil
}

func DeleteProjectMember(projectId int64, userId int64) error {
	if err := DB.Where("project_id = ?", projectId).Where("user_id = ?", userId).Delete(&ProjectMember{}).Error; err != nil {
		return err
	}
	return nil
}

func DeleteProjectMemberByProject(projectId int64) error {
	if err := DB.Where("project_id = ?", projectId).Delete(&ProjectMember{}).Error; err != nil {
		return err
	}
	return nil
}

func HardDeleteProjectMemberByProject(projectId int64) error {
	if err := DB.Unscoped().Where("project_id = ?", projectId).Delete(&ProjectMember{}).Error; err != nil {
		return err
	}
	return nil
}

func ListProjectMemberByProject(projectId int64, page int32, pageSize int32) ([]*ProjectMember, error) {
	var projectMembers []*ProjectMember
	err := DB.Where("project_id", projectId).Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Preload("User").Find(&projectMembers).Error
	if err != nil {
		return nil, err
	}
	return projectMembers, nil
}

func ListProjectMemberByUser(userId int64, page int32, pageSize int32) ([]*ProjectMember, error) {
	var projectMembers []*ProjectMember
	err := DB.Where("user_id", userId).Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Preload("Project").Find(&projectMembers).Error
	if err != nil {
		return nil, err
	}
	return projectMembers, nil
}
