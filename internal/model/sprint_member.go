package model

type SprintMember struct {
	Model

	SprintID   int64   `gorm:"column:sprint_id"`
	UserID     int64   `gorm:"column:user_id"`
	Role       string  `gorm:"column:role;size:255"`
	Capacity   float32 `gorm:"column:capacity"`
	FromSprint Sprint  `gorm:"foreignKey:SprintID"`
	MemberUser User    `gorm:"foreignKey:UserID"`
}

//
//func CreateSprintMember(sprintId int64, userId int64) (int64, error) {
//    sm := SprintMember{
//        SprintId: sprintId,
//        UserId:   userId,
//    }
//    if err := model.DB.Create(&sm).Error; err != nil {
//        return 0, err
//    }
//    return sm.Id, nil
//}
//
//func CreateSprintMembers(sprintId int64, userIds []int64) ([]int64, error) {
//    var sprintMembers []*SprintMember
//    for _, userId := range userIds {
//        sprintMembers = append(sprintMembers, &SprintMember{
//            SprintId: sprintId,
//            UserId:   userId,
//        })
//    }
//    if err := model.DB.Create(&sprintMembers).Error; err != nil {
//        return nil, err
//    }
//    var sprintIds []int64
//    for _, sprintMember := range sprintMembers {
//        sprintIds = append(sprintIds, sprintMember.Id)
//    }
//    return sprintIds, nil
//}
//
//func UpdateSprintMember(sprintId int64, userId int64, roleId int64, capacity int32) error {
//    if err := model.DB.Model(&SprintMember{}).Where("sprint_id = ?", sprintId).Where("user_id = ?", userId).Update("role_id", roleId).Update("capacity", capacity).Error; err != nil {
//        return err
//    }
//    return nil
//}
//
//func DeleteSprintMember(sprintId int64, userId int64) error {
//    if err := model.DB.Where("sprint_id = ?", sprintId).Where("user_id = ?", userId).Delete(&SprintMember{}).Error; err != nil {
//        return err
//    }
//    return nil
//}
//
//func HardDeleteSprintMember(sprintId int64, userId int64) error {
//    if err := model.DB.Unscoped().Where("sprint_id = ?", sprintId).Where("user_id = ?", userId).Delete(&SprintMember{}).Error; err != nil {
//        return err
//    }
//    return nil
//}
//
//func DeleteSprintMemberBySprint(sprintId int64) error {
//    if err := model.DB.Where("sprint_id = ?", sprintId).Delete(&SprintMember{}).Error; err != nil {
//        return err
//    }
//    return nil
//}
//
//func HardDeleteSprintMemberBySprint(sprintId int64) error {
//    if err := model.DB.Unscoped().Where("sprint_id = ?", sprintId).Delete(&SprintMember{}).Error; err != nil {
//        return err
//    }
//    return nil
//}
//
//func ListSprintMemberBySprint(sprintId int64, page int32, pageSize int32) ([]*SprintMember, error) {
//    var sprintMembers []*SprintMember
//    err := model.DB.Where("sprint_id", sprintId).Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Preload("User").Preload("Role").Find(&sprintMembers).Error
//    if err != nil {
//        return nil, err
//    }
//    return sprintMembers, nil
//}
//
//func CountSprintMemberBySprint(sprintId int64) int64 {
//    count := int64(0)
//    model.DB.Model(&SprintMember{}).Where("sprint_id = ?", sprintId).Count(&count)
//    return count
//}
//
//func ListSprintMemberByUser(userId int64, page int32, pageSize int32) ([]*SprintMember, error) {
//    var sprintMembers []*SprintMember
//    err := model.DB.Where("user_id", userId).Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Preload("Sprint").Find(&sprintMembers).Error
//    if err != nil {
//        return nil, err
//    }
//    return sprintMembers, nil
//}
//
//func CountSprintMemberByUser(userId int64) int64 {
//    count := int64(0)
//    model.DB.Model(&SprintMember{}).Where("user_id = ?", userId).Count(&count)
//    return count
//}
