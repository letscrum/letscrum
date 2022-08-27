package service

import (
	"context"
	generalv1 "github.com/letscrum/letscrum/api/general/v1"
	letscrumv1 "github.com/letscrum/letscrum/api/letscrum/v1"
	"github.com/letscrum/letscrum/internal/dao"
	"github.com/letscrum/letscrum/pkg/build"
)

type LetscrumServiceInterface interface {
	GetVersion(context.Context) (*generalv1.GetVersionResponse, error)
}

type LetscrumService struct {
	letscrumv1.UnimplementedLetscrumServer
	dao dao.LetscrumDao
}

func (s *LetscrumService) GetVersion(_ context.Context) (*generalv1.GetVersionResponse, error) {
	ver := build.Version()
	response := generalv1.GetVersionResponse{
		Version: &generalv1.Version{
			Version:   ver.Version,
			GitCommit: ver.GitCommit,
			BuildDate: ver.BuildDate,
			GoVersion: ver.GoVersion,
		},
	}
	return &response, nil
}
