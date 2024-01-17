package mysql

import (
	"github.com/letscrum/letscrum/internal/model"
	"gorm.io/gorm"
)

type ProjectMemberDao struct {
	DB *gorm.DB
}

func (p ProjectMemberDao) GetByProject(project model.Project) (*model.ProjectMember, error) {
	var pMember *model.ProjectMember
	err := p.DB.Where("project_id = ?", project.ID).Where("user_id", project.CreatedUser.ID).Preload("MemberUser").Find(&pMember).Error
	if err != nil {
		return nil, err
	}
	return pMember, nil
}

func (p ProjectMemberDao) Get(projectMember model.ProjectMember) (*model.ProjectMember, error) {
	var pMember *model.ProjectMember
	err := p.DB.Where("id = ?", projectMember.ID).Preload("MemberUser").Find(&pMember).Error
	if err != nil {
		return nil, err
	}
	return pMember, nil
}

func (p ProjectMemberDao) ListByProject(project model.Project, page, size int32) ([]*model.ProjectMember, error) {
	var pMembers []*model.ProjectMember
	err := p.DB.Where("project_id = ?", project.ID).Limit(int(size)).Offset(int((page - 1) * size)).Preload("MemberUser").Find(&pMembers).Error
	if err != nil {
		return nil, err
	}
	return pMembers, nil
}

func (p ProjectMemberDao) BatchAdd(projectMembers []model.ProjectMember) ([]*model.ProjectMember, error) {
	var pMembers []*model.ProjectMember
	for _, pm := range projectMembers {
		projectMember := model.ProjectMember{
			ProjectID: pm.ProjectID,
			UserID:    pm.UserID,
			IsAdmin:   pm.IsAdmin,
		}
		pMembers = append(pMembers, &projectMember)
	}
	if len(pMembers) > 0 {
		if err := p.DB.Create(&pMembers).Error; err != nil {
			return nil, err
		}
	}
	return pMembers, nil
}

func (p ProjectMemberDao) BatchUpdate(projectMembers []model.ProjectMember) ([]*model.ProjectMember, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProjectMemberDao) BatchRemove(projectMembers []model.ProjectMember) ([]*model.ProjectMember, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProjectMemberDao) Count() int64 {
	//TODO implement me
	panic("implement me")
}

func NewProjectMemberDao(d *gorm.DB) *ProjectMemberDao {
	return &ProjectMemberDao{d}
}
