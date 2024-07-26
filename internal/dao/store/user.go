package store

import (
	"github.com/letscrum/letscrum/internal/model"
	"gorm.io/gorm"
)

type UserDao struct {
	DB *gorm.DB
}

func (d *UserDao) List(page, size int32, keyword string) ([]*model.User, error) {
	var users []*model.User
	err := d.DB.Where("name LIKE ?", "%"+keyword+"%").Limit(int(size)).Offset(int((page - 1) * size)).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (d *UserDao) ListByIds(page, size int32, ids []int64) ([]*model.User, error) {
	var users []*model.User
	// get users by ids from database
	err := d.DB.Where("id IN ?", ids).Limit(int(size)).Offset(int((page - 1) * size)).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (d *UserDao) SignIn(name, password string) (*model.User, error) {
	var u *model.User
	if err := d.DB.Where("name = ?", name).Where("password = ?", password).Find(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (d *UserDao) Get(id int64) (*model.User, error) {
	var u *model.User
	if err := d.DB.Where("id = ?", id).Find(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (d *UserDao) Count(keyword string) int64 {
	count := int64(0)
	d.DB.Where("name LIKE ?", "%"+keyword+"%").Model(&model.User{}).Count(&count)
	return count
}

func (d *UserDao) Create(name, email, password string, isSuperAdmin bool) (*model.User, error) {
	user := &model.User{
		Name:         name,
		Email:        email,
		Password:     password,
		IsSuperAdmin: isSuperAdmin,
	}
	if err := d.DB.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (d *UserDao) SetSuperAdmin(id int64, isAdmin bool) (*model.User, error) {
	if err := d.DB.Model(&model.User{}).Where("id = ?", id).Update("is_super_admin", isAdmin).Error; err != nil {
		return nil, err
	}
	return d.Get(id)
}

func (d *UserDao) ListSuperAdmins() ([]*model.User, error) {
	var admins []*model.User
	if err := d.DB.Where("is_super_admin = ?", true).Find(&admins).Error; err != nil {
		return nil, err
	}
	return admins, nil
}

func (d *UserDao) UpdatePassword(id int64, oldPassword, newPassword string) (*model.User, error) {
	if err := d.DB.Model(&model.User{}).Where("id = ?", id).Where("password = ?", oldPassword).Update("password", newPassword).Error; err != nil {
		return nil, err
	}
	return d.Get(id)
}

func (d *UserDao) ResetPassword(id int64, newPassword string) (*model.User, error) {
	if err := d.DB.Model(&model.User{}).Where("id = ?", id).Update("password", newPassword).Error; err != nil {
		return nil, err
	}
	return d.Get(id)
}

func NewUserDao(d *gorm.DB) *UserDao {
	return &UserDao{d}
}
