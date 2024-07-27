package store

import (
	"github.com/google/uuid"
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

func (d *UserDao) ListByIds(page, size int32, ids []uuid.UUID) ([]*model.User, error) {
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

func (d *UserDao) Get(user model.User) (*model.User, error) {
	var u *model.User
	if err := d.DB.Where("id = ?", user.Id).Find(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (d *UserDao) GetByName(user model.User) (*model.User, error) {
	var u *model.User
	if err := d.DB.Where("name = ?", user.Name).Find(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (d *UserDao) Count(keyword string) int64 {
	count := int64(0)
	d.DB.Where("name LIKE ?", "%"+keyword+"%").Model(&model.User{}).Count(&count)
	return count
}

func (d *UserDao) Create(user model.User) (*model.User, error) {
	if err := d.DB.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *UserDao) SetSuperAdmin(user model.User, isAdmin bool) (*model.User, error) {
	if err := d.DB.Model(&model.User{}).Where("id = ?", user.Id).Update("is_super_admin", isAdmin).Error; err != nil {
		return nil, err
	}
	return d.Get(user)
}

func (d *UserDao) ListSuperAdmins() ([]*model.User, error) {
	var admins []*model.User
	if err := d.DB.Where("is_super_admin = ?", true).Find(&admins).Error; err != nil {
		return nil, err
	}
	return admins, nil
}

func (d *UserDao) UpdatePassword(user model.User, newPassword string) (*model.User, error) {
	if err := d.DB.Model(&model.User{}).Where("id = ?", user.Id).Where("password = ?", user.Password).Update("password", newPassword).Error; err != nil {
		return nil, err
	}
	return d.Get(user)
}

func (d *UserDao) ResetPassword(user model.User, newPassword string) (*model.User, error) {
	if err := d.DB.Model(&model.User{}).Where("id = ?", user.Id).Update("password", newPassword).Error; err != nil {
		return nil, err
	}
	return d.Get(user)
}

func NewUserDao(d *gorm.DB) *UserDao {
	return &UserDao{d}
}
