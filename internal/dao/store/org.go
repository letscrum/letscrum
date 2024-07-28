package store

import (
	"github.com/letscrum/letscrum/internal/model"
	"gorm.io/gorm"
)

type OrgDao struct {
	DB *gorm.DB
}

func NewOrgDao(d *gorm.DB) *OrgDao {
	return &OrgDao{d}
}

func (d OrgDao) Get(org model.Org) (model.Org, error) {
	var o model.Org
	rows := d.DB.Where("id = ?", org.Id).Find(&o).RowsAffected
	if rows == 0 {
		return model.Org{}, gorm.ErrRecordNotFound
	}
	return o, nil
}

func (d OrgDao) List(page, size int32, keyword string) ([]model.Org, error) {
	var orgs []model.Org
	if err := d.DB.Where("name LIKE ?", "%"+keyword+"%").Or("display_name LIKE ?", "%"+keyword+"%").Limit(int(size)).Offset(int((page - 1) * size)).Preload("CreatedUser").Order("updated_at desc").Find(&orgs).Error; err != nil {
		return nil, err
	}
	return orgs, nil
}

func (d OrgDao) Count(keyword string) int64 {
	count := int64(0)
	d.DB.Where("name LIKE ?", "%"+keyword+"%").Or("display_name LIKE ?", "%"+keyword+"%").Model(&model.Org{}).Count(&count)
	return count
}

func (d OrgDao) ListVisibleOrg(page, size int32, keyword string, user model.User) ([]model.Org, error) {
	// get orgs that user is member or createdby is user
	var orgs []model.Org
	if err := d.DB.Where("id IN (SELECT org_id FROM org_user WHERE user_id = ?) OR created_by = ?", user.Id, user.Id).Where("name LIKE ? Or display_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%").Limit(int(size)).Offset(int((page - 1) * size)).Preload("CreatedUser").Order("updated_at desc").Find(&orgs).Error; err != nil {
		return nil, err
	}
	return orgs, nil
}

func (d OrgDao) CountVisibleOrg(keyword string, user model.User) int64 {
	count := int64(0)
	d.DB.Where("id IN (SELECT org_id FROM org_user WHERE user_id = ?) OR created_by = ?", user.Id, user.Id).Where("name LIKE ? Or display_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%").Model(&model.Org{}).Count(&count)
	return count
}

func (d OrgDao) CountByUser(user model.User) int64 {
	count := int64(0)
	d.DB.Model(&model.Org{}).Where("created_by = ?", user.Id).Count(&count)
	return count
}

func (d OrgDao) Create(org model.Org) (model.Org, error) {
	if err := d.DB.Create(&org).Error; err != nil {
		return model.Org{}, err
	}
	return org, nil
}

func (d OrgDao) Update(org model.Org) (model.Org, error) {
	if err := d.DB.Model(&model.Org{}).Where("id = ?", org.Id).Updates(org).Error; err != nil {
		return model.Org{}, err
	}
	return org, nil
}

func (d OrgDao) Delete(org model.Org) (bool, error) {
	if err := d.DB.Where("id = ?", org.Id).Delete(&model.Org{}).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (d OrgDao) AddMember(orgUser model.OrgUser) (model.OrgUser, error) {
	// add orgUser
	if err := d.DB.Create(&orgUser).Error; err != nil {
		return model.OrgUser{}, err
	}
	return orgUser, nil
}

func (d OrgDao) AddMembers(orgUsers []model.OrgUser) ([]model.OrgUser, error) {
	// add orgUsers
	if err := d.DB.Create(&orgUsers).Error; err != nil {
		return nil, err
	}
	return orgUsers, nil
}

func (d OrgDao) RemoveMember(orgUser model.OrgUser) (bool, error) {
	// remove orgUser
	if err := d.DB.Where("org_id = ? AND user_id = ?", orgUser.OrgId, orgUser.UserId).Delete(&model.OrgUser{}).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (d OrgDao) ListMember(org model.Org) ([]model.OrgUser, error) {
	// get orguser by org id and preload user
	var orgUsers []model.OrgUser
	if err := d.DB.Where("org_id = ?", org.Id).Preload("ForOrg").Preload("Member").Find(&orgUsers).Error; err != nil {
		return nil, err
	}
	return orgUsers, nil
}

func (d OrgDao) SetAdmin(orgUser model.OrgUser, isAdmin bool) (model.OrgUser, error) {
	// set orgUser isAdmin
	if err := d.DB.Model(&model.OrgUser{}).Where("org_id = ? AND user_id = ?", orgUser.OrgId, orgUser.UserId).Update("is_admin", isAdmin).Error; err != nil {
		return model.OrgUser{}, err
	}
	return orgUser, nil
}
