package models

type SprintMember struct {
	Model

	SprintId int64  `json:"sprint_id"`
	UserId   int64  `json:"user_id"`
	RoleId   int64  `json:"role_id"`
	User     User   `gorm:"foreignKey:UserId"`
	Sprint   Sprint `gorm:"foreignKey:SprintId"`
}

func CreateSprintMember(projectId int64, userId int64, isAdmin bool) (int64, error) {
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

func UpdateSprintMember(projectId int64, userId int64, isAdmin bool) error {
	if err := DB.Model(&ProjectMember{}).Where("project_id = ?", projectId).Where("user_id = ?", userId).Update("is_admin", isAdmin).Error; err != nil {
		return err
	}
	return nil
}

func DeleteSprintMember(projectId int64, userId int64) error {
	if err := DB.Where("project_id = ?", projectId).Where("user_id = ?", userId).Delete(&ProjectMember{}).Error; err != nil {
		return err
	}
	return nil
}

func HardDeleteSprintMember(projectId int64, userId int64) error {
	if err := DB.Unscoped().Where("project_id = ?", projectId).Where("user_id = ?", userId).Delete(&ProjectMember{}).Error; err != nil {
		return err
	}
	return nil
}

func DeleteSprintMemberBySprint(projectId int64) error {
	if err := DB.Where("project_id = ?", projectId).Delete(&ProjectMember{}).Error; err != nil {
		return err
	}
	return nil
}

func HardDeleteSprintMemberBySprint(projectId int64) error {
	if err := DB.Unscoped().Where("project_id = ?", projectId).Delete(&ProjectMember{}).Error; err != nil {
		return err
	}
	return nil
}

func ListSprintMemberBySprint(projectId int64, page int32, pageSize int32) ([]*ProjectMember, error) {
	var projectMembers []*ProjectMember
	err := DB.Where("project_id", projectId).Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Preload("User").Find(&projectMembers).Error
	if err != nil {
		return nil, err
	}
	return projectMembers, nil
}

func CountSprintMemberBySprint(projectId int64) int64 {
	count := int64(0)
	DB.Model(&ProjectMember{}).Where("project_id = ?", projectId).Count(&count)
	return count
}
