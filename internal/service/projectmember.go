package service

//
//func Create(projectId int64, userId int64, isAdmin bool) (int64, error) {
//	projectMemberId, err := projectmembermodel.CreateProjectMember(projectId, userId, isAdmin)
//	if err != nil {
//		return 0, err
//	}
//	return projectMemberId, nil
//}
//
//func ListProjectMemberByProject(projectId int64, page int32, pageSize int32) ([]*userV1.User, int64, error) {
//	projectMembers, err := projectmembermodel.ListProjectMemberByProject(projectId, page, pageSize)
//	if err != nil {
//		return nil, 0, err
//	}
//	var members []*userV1.User
//	for _, m := range projectMembers {
//		members = append(members, &userV1.User{
//			Id:           m.User.Id,
//			Name:         m.User.Name,
//			Email:        m.User.Email,
//			IsSuperAdmin: m.User.IsSuperAdmin,
//			IsAdmin:      m.IsAdmin,
//		})
//	}
//	count := projectmembermodel.CountProjectMemberByProject(projectId)
//	return members, count, nil
//}
//
//func ListProjectByUser(userId int64, page int32, pageSize int32) ([]*projectV1.Project, int64, error) {
//	projectMembers, err := projectmembermodel.ListProjectMemberByUser(userId, page, pageSize)
//	if err != nil {
//		return nil, 0, err
//	}
//	var projects []*projectV1.Project
//	for _, p := range projectMembers {
//		projects = append(projects, &projectV1.Project{
//			Id:          p.Project.Id,
//			Name:        p.Project.Name,
//			DisplayName: p.Project.DisplayName,
//		})
//	}
//	count := projectmembermodel.CountProjectMemberByUser(userId)
//	return projects, count, nil
//}
