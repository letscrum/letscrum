package models

import (
	generalV1 "github.com/letscrum/letscrum/apis/general/v1"
	userV1 "github.com/letscrum/letscrum/apis/user/v1"
)

type User struct {
	Model

	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateUser(user *userV1.User) error {
	u := User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	//var pInDB *Project
	//errPName := db.Where("name = ?", p.Name).First(&pInDB).Error
	//if errPName != nil && errPName != gorm.ErrRecordNotFound {
	//	return errPName
	//}
	//if pInDB != nil && pInDB.Name == p.Name {
	//	return fmt.Errorf("duplicate project name")
	//}

	if err := db.Create(&u).Error; err != nil {
		return err
	}
	return nil
}

func ListUser(pagination *generalV1.Pagination) ([]*User, error) {
	var users []*User
	err := db.Limit(int(pagination.PageSize)).Offset(int((pagination.Page - 1) * pagination.PageSize)).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func CountUser() int64 {
	count := int64(0)
	db.Model(&User{}).Count(&count)
	return count
}

func UpdateUser(name string, user *userV1.User) error {
	u := User{
		Email:    user.Email,
		Password: user.Password,
	}
	if err := db.Model(&User{}).Where("name = ?", name).Update("email", u.Email).Update("password", u.Password).Error; err != nil {
		return err
	}
	return nil
}

func DeleteUser(name string) error {
	if err := db.Where("name = ?", name).Delete(&User{}).Error; err != nil {
		return err
	}
	return nil
}

func GetUser(name string) (*User, error) {
	var u *User
	if err := db.Where("name = ?", name).Find(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func GetUserWithPassword(name string, password string) (*User, error) {
	var u *User
	if err := db.Where("name = ?", name).Where("password = ?", password).Find(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}
