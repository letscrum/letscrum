package model

type User struct {
	Model

	Name         string `gorm:"column:name;size:255;index:idx_name,unique" json:"name,omitempty"`
	Email        string `gorm:"column:email;size:255;index:idx_name,unique" json:"email,omitempty"`
	Password     string `gorm:"column:password;size:255" json:"password,omitempty"`
	IsSuperAdmin bool   `gorm:"column:is_super_admin" json:"is_super_admin,omitempty"`
}

//
//func CreateUser(name string, email string, password string, isSuperAdmin bool) (int64, error) {
//    u := User{
//        Name:         name,
//        Email:        email,
//        Password:     password,
//        IsSuperAdmin: isSuperAdmin,
//    }
//
//    //var pInDB *Project
//    //errPName := db.Where("name = ?", p.Name).First(&pInDB).Error
//    //if errPName != nil && errPName != gorm.ErrRecordNotFound {
//    //    return errPName
//    //}
//    //if pInDB != nil && pInDB.Name == p.Name {
//    //    return fmt.Errorf("duplicate project name")
//    //}
//
//    if err := DB.Create(&u).Error; err != nil {
//        return 0, err
//    }
//    return u.Id, nil
//}
//
//func ListUser(keyword string, page int32, pageSize int32) ([]*User, error) {
//    var users []*User
//    err := DB.Where("name LIKE ?", "%"+keyword+"%").Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Find(&users).Error
//    if err != nil {
//        return nil, err
//    }
//    return users, nil
//}
//
//func CountUser(keyword string) int64 {
//    count := int64(0)
//    DB.Model(&User{}).Where("name LIKE ?", "%"+keyword+"%").Count(&count)
//    return count
//}
//
//func UpdateUser(name string, email string, password string) error {
//    u := User{
//        Email:    email,
//        Password: password,
//    }
//    if err := DB.Model(&User{}).Where("name = ?", name).Update("email", u.Email).Update("password", u.Password).Error; err != nil {
//        return err
//    }
//    return nil
//}
//
//func DeleteUser(name string) error {
//    if err := DB.Where("name = ?", name).Delete(&User{}).Error; err != nil {
//        return err
//    }
//    return nil
//}
//
//func GetUser(id int64) (*User, error) {
//    var u *User
//    if err := DB.Where("id = ?", id).Find(&u).Error; err != nil {
//        return nil, err
//    }
//    return u, nil
//}
//
//func GetUserWithPassword(name string, password string) (*User, error) {
//    var u *User
//    if err := DB.Where("name = ?", name).Where("password = ?", password).Find(&u).Error; err != nil {
//        return nil, err
//    }
//    return u, nil
//}
