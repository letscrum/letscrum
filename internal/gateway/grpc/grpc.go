package grpc

import (
	"context"
	v1 "github.com/letscrum/letscrum/api/letscrum/v1"
	golog "log"
	"net"
	"os"
	"time"

	"github.com/letscrum/letscrum/internal/dao"
	"github.com/letscrum/letscrum/internal/dao/mysql"
	"github.com/letscrum/letscrum/internal/service"
	"github.com/letscrum/letscrum/pkg/db"
	"github.com/letscrum/letscrum/pkg/log"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"gorm.io/gorm/logger"
)

func Run(ctx context.Context, network, address string) error {
	//init grpc server and run
	l, err := net.Listen(network, address)
	if err != nil {
		return err
	}
	go func() {
		defer func() error {
			if err := l.Close(); err != nil {
				return err
			}
			return nil
		}()
		<-ctx.Done()
	}()
	s := grpc.NewServer()

	var daoInterface dao.Interface
	if daoInterface, err = initDao(); err != nil {
		return err
	}
	projectService := service.NewProjectService(daoInterface)
	v1.RegisterProjectServer(s, projectService)
	letscrumService := service.NewLetscrumService(daoInterface)
	v1.RegisterLetscrumServer(s, letscrumService)
	demoService := service.NewDemoService()
	v1.RegisterDemoServer(s, demoService)

	go func() {
		defer s.GracefulStop()
		<-ctx.Done()
	}()

	go func() error {
		log.L(ctx).Infof("grpc listen on:%s\n", address)
		if err := s.Serve(l); err != nil {
			return err
		}
		return nil
	}()

	return nil
}

func initDao() (dao.Interface, error) {
	newLogger := logger.New(
		golog.New(os.Stdout, "", golog.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
	options := db.Options{
		Host:                  viper.GetString("data.database.host"),
		Port:                  viper.GetString("data.database.port"),
		Username:              viper.GetString("data.database.user"),
		Password:              viper.GetString("data.database.password"),
		Database:              viper.GetString("data.database.database"),
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifeTime: 10 * time.Second,
		Logger:                newLogger,
	}
	hiveDao, err := mysql.GetDao(&options)
	if err != nil {
		return nil, err
	}
	return hiveDao, nil
}
