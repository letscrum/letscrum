package service

import "github.com/letscrum/letscrum/internal/dao"

type Service struct {
	dao dao.Interface
}

func (s *Service) ProjectService() ProjectService {
	return *NewProjectService(s.dao)
}
