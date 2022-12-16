package service

import "github.com/letscrum/letscrum/internal/dao"

type Service struct {
	dao dao.Interface
}

func (s *Service) LetscrumService() LetscrumService {
	return *NewLetscrumService(s.dao)
}

func (s *Service) ProjectService() ProjectService {
	return *NewProjectService(s.dao)
}

func (s *Service) ProjectMemberService() ProjectMemberService {
	return *NewProjectMemberService(s.dao)
}

func (s *Service) SprintService() SprintService {
	return *NewSprintService(s.dao)
}

func (s *Service) SprintMemberService() SprintMemberService {
	return *NewSprintMemberService(s.dao)
}
