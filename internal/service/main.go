package service

import "github.com/letscrum/letscrum/internal/dao"

type Service struct {
	dao dao.Interface
}

func (s *Service) ProjectService() ProjectService {
	return *NewProjectService(s.dao)
}

func (s *Service) LetscrumService() LetscrumService {
	return *NewLetscrumService(s.dao)
}

func (s *Service) DemoService() DemoService {
	return *NewDemoService()
}

func (s *Service) DemoDbService() DemoDbService {
	return *NewDemoDbService(s.dao)
	// Replace to follow line if you don't need database operate
	// return *NewDemoService()
}
