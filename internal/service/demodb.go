package service

import (
	"context"
	generalv1 "github.com/letscrum/letscrum/api/general/v1"
	letscrumv1 "github.com/letscrum/letscrum/api/letscrum/v1"
	"github.com/letscrum/letscrum/internal/dao"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DemoDbService struct {
	// This is generated by protoc
	letscrumv1.UnimplementedDemoDbServer
	// Remove follow line if you don't need database operate
	dao dao.DemoDbDao
}

func NewDemoDbService(dao dao.Interface) *DemoDbService {
	return &DemoDbService{dao: dao.DemoDbDao()}
}

func (s *DemoDbService) DemoDb(ctx context.Context, req *generalv1.DemoDbRequest) (*generalv1.DemoDbResponse, error) {
	demoDb, err := s.dao.DemoDb(ctx, req.DemoDb)
	if err != nil {
		result := status.Convert(err)
		if result.Code() == codes.NotFound {
			return nil, status.Errorf(codes.NotFound, "get err: %s not found", req.DemoDb)
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &generalv1.DemoDbResponse{
		DemoDb: &generalv1.DemoDb{
			DemoDb: demoDb.DemoDb,
		},
	}, nil
}
