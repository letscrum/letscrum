package service

//
//func Create(projectId int64, name string, startDate time.Time, endDate time.Time) (int64, error) {
//	sprintId, err := sprintmodel.CreateSprint(projectId, name, startDate, endDate)
//	if err != nil {
//		return 0, err
//	}
//	projectMembers, errGetMembers := models2.ListProjectMemberByProject(projectId, 1, 999)
//	if errGetMembers != nil {
//		errDeleteSprint := sprintmodel.HardDeleteSprint(sprintId)
//		if errDeleteSprint != nil {
//			return 0, errDeleteSprint
//		}
//		return 0, err
//	}
//	var userIds []int64
//	for _, projectMember := range projectMembers {
//		userIds = append(userIds, projectMember.UserId)
//	}
//	_, errCreateSprintMembers := sprintmembermodel.CreateSprintMembers(sprintId, userIds)
//	if errCreateSprintMembers != nil {
//		errDeleteSprint := sprintmodel.HardDeleteSprint(sprintId)
//		if errDeleteSprint != nil {
//			return 0, errDeleteSprint
//		}
//		return 0, err
//	}
//	return sprintId, nil
//}
//
//func List(projectId int64, page int32, pageSize int32) ([]*projectV1.Sprint, int64, error) {
//	sprints, err := sprintmodel.ListSprintByProject(projectId, page, pageSize)
//	if err != nil {
//		return nil, 0, err
//	}
//	var list []*projectV1.Sprint
//	for _, p := range sprints {
//		list = append(list, &projectV1.Sprint{
//			Id:        p.Id,
//			Name:      p.Name,
//			StartDate: p.StartDate.Unix(),
//			EndDate:   p.EndDate.Unix(),
//			CreatedAt: p.CreatedAt.Unix(),
//			UpdatedAt: p.UpdatedAt.Unix(),
//		})
//	}
//	count := sprintmodel.CountSprintByProject(projectId)
//	return list, count, nil
//}
//
//func Update(id int64, name string, startDate time.Time, endDate time.Time) error {
//	if err := sprintmodel.UpdateSprint(id, name, startDate, endDate); err != nil {
//		return err
//	}
//	return nil
//}
//
//func Delete(id int64) error {
//	if err := sprintmodel.DeleteSprint(id); err != nil {
//		return err
//	}
//	return nil
//}
//
//func HardDelete(id int64) error {
//	if err := sprintmembermodel.HardDeleteSprintMemberBySprint(id); err != nil {
//		return err
//	}
//	if err := sprintmodel.HardDeleteSprint(id); err != nil {
//		return err
//	}
//	return nil
//}
//
//func Get(id int64) (*projectV1.Sprint, error) {
//	s, err := sprintmodel.GetSprint(id)
//	if err != nil {
//		return nil, err
//	}
//	sprint := &projectV1.Sprint{
//		Id:        s.Id,
//		Name:      s.Name,
//		StartDate: s.StartDate.Unix(),
//		EndDate:   s.EndDate.Unix(),
//		CreatedAt: s.CreatedAt.Unix(),
//		UpdatedAt: s.UpdatedAt.Unix(),
//	}
//	return sprint, nil
//}
