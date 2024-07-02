package model

type Project struct {
	Model

	Name        string `gorm:"column:name;size:255"`
	DisplayName string `gorm:"column:display_name;size:255"`
	Description string `gorm:"column:description;size:5000"`
	Members     string `gorm:"column:members;size:5000"`
	CreatedBy   int64  `gorm:"column:create_by"`
	CreatedUser User   `gorm:"foreignKey:CreatedBy"`
}

//
//func CreateProject(name string, displayName string, createdUserId int64) (int64, error) {
//    p := Project{
//        Name:        name,
//        DisplayName: displayName,
//        CreatedBy:   createdUserId,
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
//    if err := DB.Create(&p).Error; err != nil {
//        return 0, err
//    }
//    return p.Id, nil
//}
//
//func UpdateProject(id int64, displayName string) error {
//    if err := DB.Model(&Project{}).Where("id = ?", id).Update("display_name", displayName).Error; err != nil {
//        return err
//    }
//    return nil
//}
//
//func DeleteProject(id int64) error {
//    if err := DB.Where("id = ?", id).Delete(&Project{}).Error; err != nil {
//        return err
//    }
//    return nil
//}
//
//func HardDeleteProject(id int64) error {
//    if err := DB.Unscoped().Where("id = ?", id).Delete(&Project{}).Error; err != nil {
//        return err
//    }
//    return nil
//}
