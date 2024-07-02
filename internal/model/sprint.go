package model

import (
	"time"
)

type Sprint struct {
	Model

	ProjectID   int64     `gorm:"column:project_id"`
	Name        string    `gorm:"column:name;size:255"`
	Members     string    `gorm:"column:members;size:5000"`
	StartDate   time.Time `gorm:"column:start_date"`
	EndDate     time.Time `gorm:"column:end_date"`
	FromProject Project   `gorm:"foreignKey:ProjectID"`
}

//
//func CreateSprint(projectId int64, name string, startDate time.Time, endDate time.Time) (int64, error) {
//    s := Sprint{
//        ProjectId: projectId,
//        Name:      name,
//        StartDate: startDate,
//        EndDate:   endDate,
//    }
//
//    if err := model.DB.Create(&s).Error; err != nil {
//        return 0, err
//    }
//    return s.Id, nil
//}
//
//func ListSprintByProject(projectId int64, page int32, pageSize int32) ([]*Sprint, error) {
//    var sprints []*Sprint
//    err := model.DB.Where("project_id = ?", projectId).Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Preload("Project").Find(&sprints).Error
//    if err != nil {
//        return nil, err
//    }
//    return sprints, nil
//}
//
//func CountSprintByProject(projectId int64) int64 {
//    count := int64(0)
//    model.DB.Model(&Sprint{}).Where("project_id = ?", projectId).Count(&count)
//    return count
//}
//
//func UpdateSprint(id int64, name string, startDate time.Time, endDate time.Time) error {
//    if err := model.DB.Model(&Sprint{}).Where("id = ?", id).Update("name", name).Update("start_date", startDate).Update("end_date", endDate).Error; err != nil {
//        return err
//    }
//    return nil
//}
//
//func DeleteSprint(id int64) error {
//    if err := model.DB.Where("id = ?", id).Delete(&Sprint{}).Error; err != nil {
//        return err
//    }
//    return nil
//}
//
//func HardDeleteSprint(id int64) error {
//    if err := model.DB.Unscoped().Where("id = ?", id).Delete(&Sprint{}).Error; err != nil {
//        return err
//    }
//    return nil
//}
//
//func GetSprint(id int64) (*Sprint, error) {
//    var p *Sprint
//    if err := model.DB.Where("id = ?", id).Preload("Project").Find(&p).Error; err != nil {
//        return nil, err
//    }
//    return p, nil
//}
