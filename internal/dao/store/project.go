package store

import (
	"github.com/letscrum/letscrum/internal/model"
	"gorm.io/gorm"
)

type ProjectDao struct {
	DB *gorm.DB
}

func (d ProjectDao) Get(project model.Project) (*model.Project, error) {
	var p *model.Project
	if err := d.DB.Where("org_id = ? AND id = ?", project.OrgId, project.Id).Preload("CreatedUser").Find(&p).Error; err != nil {
		return nil, err
	}
	return p, nil
}

func (d ProjectDao) Create(project model.Project) (*model.Project, error) {
	if err := d.DB.Create(&project).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (d ProjectDao) Update(project model.Project) (*model.Project, error) {
	if err := d.DB.Model(&model.Project{}).Where("org_id = ? AND id = ?", project.OrgId, project.Id).Updates(project).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (d ProjectDao) Delete(project model.Project) (bool, error) {
	if err := d.DB.Where("org_id = ? AND id = ?", project.OrgId, project.Id).Delete(&model.Project{}).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (d ProjectDao) Count(org model.Org, keyword string) int64 {
	count := int64(0)
	d.DB.Where("org_id = ?", org.Id).Where("name LIKE ? Or display_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%").Model(&model.Project{}).Count(&count)
	return count
}

func (d ProjectDao) List(org model.Org, page, size int32, keyword string) ([]*model.Project, error) {
	var projects []*model.Project
	err := d.DB.Where("org_id = ?", org.Id).Where("name LIKE ? Or display_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%").Limit(int(size)).Offset(int((page - 1) * size)).Preload("CreatedUser").Order("updated_at desc").Find(&projects).Error
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (d ProjectDao) ListVisibleProject(org model.Org, page, size int32, keyword string, user model.User) ([]*model.Project, error) {
	var projects []*model.Project
	if user.IsSuperAdmin {
		err := d.DB.Where("org_id = ?", org.Id).Where("name LIKE ? OR display_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%").Limit(int(size)).Offset(int((page - 1) * size)).Preload("CreatedUser").Order("updated_at desc").Find(&projects).Error
		if err != nil {
			return nil, err
		}
	} else {
		// select project where members include user.Id
		err := d.DB.Where("org_id = ? AND (members LIKE ? OR created_by = ?)", org.Id, "%"+`{"user_id":"`+user.Id.String()+`"`+",%", user.Id.String()).Where("(name LIKE ? Or display_name LIKE ?)", "%"+keyword+"%", "%"+keyword+"%").Limit(int(size)).Offset(int((page - 1) * size)).Preload("CreatedUser").Order("updated_at desc").Find(&projects).Error
		if err != nil {
			return nil, err
		}
	}
	return projects, nil
}

func (d ProjectDao) CountVisibleProject(org model.Org, keyword string, user model.User) int64 {
	count := int64(0)
	if user.IsSuperAdmin {
		d.DB.Where("org_id = ?", org.Id).Where("name LIKE ? OR display_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%").Model(&model.Project{}).Count(&count)
	} else {
		// select project where members include user.Id
		d.DB.Where("org_id = ? AND (members LIKE ? OR created_by = ?)", org.Id, "%"+`{"user_id":"`+user.Id.String()+`"`+",%", user.Id.String()).Where("(name LIKE ? Or display_name LIKE ?)", "%"+keyword+"%", "%"+keyword+"%").Model(&model.Project{}).Count(&count)
	}
	return count
}

func (d ProjectDao) CountByOrg(org model.Org) int64 {
	count := int64(0)
	d.DB.Where("org_id = ?", org.Id).Model(&model.Project{}).Count(&count)
	return count
}

func NewProjectDao(d *gorm.DB) *ProjectDao {
	return &ProjectDao{d}
}
