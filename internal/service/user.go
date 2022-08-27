package service

//
//func Create2(name string, email string, password string, isSuperAdmin bool) (int64, error) {
//	id, err := model.CreateUser(name, email, password, isSuperAdmin)
//	if err != nil {
//		return 0, err
//	}
//	return id, nil
//}
//
//func List2(keyword string, page int32, pageSize int32) ([]*userV1.User, int64, error) {
//	users, err := model.ListUser(keyword, page, pageSize)
//	if err != nil {
//		return nil, 0, err
//	}
//	var list []*userV1.User
//	for _, u := range users {
//		list = append(list, &userV1.User{
//			Id:        u.Id,
//			Name:      u.Name,
//			Email:     u.Email,
//			CreatedAt: u.CreatedAt.Unix(),
//			UpdatedAt: u.UpdatedAt.Unix(),
//		})
//	}
//	count := model.CountUser(keyword)
//	return list, count, nil
//}
//
//func Update2(user *userV1.User) error {
//	if err := model.UpdateUser(user.Name, user.Email, user.Password); err != nil {
//		return err
//	}
//	return nil
//}
//
//func Delete2(name string) error {
//	if err := model.DeleteUser(name); err != nil {
//		return err
//	}
//	return nil
//}
//
//func Get2(id int64) (*userV1.User, error) {
//	u, err := model.GetUser(id)
//	if err != nil {
//		return nil, err
//	}
//	user := &userV1.User{
//		Id:        u.Id,
//		Name:      u.Name,
//		Email:     u.Email,
//		CreatedAt: u.CreatedAt.Unix(),
//		UpdatedAt: u.UpdatedAt.Unix(),
//	}
//	return user, nil
//}
//
//func SignIn2(name string, password string) (*userV1.User, error) {
//	u, err := model.GetUserWithPassword(name, password)
//	if err != nil {
//		return nil, err
//	}
//	accessToken, refreshToken, errGenTokens := utils.GenerateTokens(strconv.FormatInt(u.Id, 10))
//	if errGenTokens != nil {
//		return nil, errGenTokens
//	}
//	user := &userV1.User{
//		Id:           u.Id,
//		Name:         u.Name,
//		Email:        u.Email,
//		IsSuperAdmin: u.IsSuperAdmin,
//		CreatedAt:    u.CreatedAt.Unix(),
//		UpdatedAt:    u.UpdatedAt.Unix(),
//		Token: &userV1.Token{
//			AccessToken:  accessToken,
//			RefreshToken: refreshToken,
//		},
//	}
//	return user, nil
//}
