package grpc

import (
	"context"
	golog "log"
	"net"
	"os"
	"time"

	v1 "github.com/letscrum/letscrum/api/letscrum/v1"

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

	//opts := grpc.UnaryInterceptor(
	//    grpc_auth.UnaryServerInterceptor(mid.Auth),
	//)

	s := grpc.NewServer()

	var daoInterface dao.Interface
	if daoInterface, err = initDao(); err != nil {
		return err
	}

	letscrumService := service.NewLetscrumService(daoInterface)
	v1.RegisterLetscrumServer(s, letscrumService)
	userService := service.NewUserService(daoInterface)
	v1.RegisterUserServer(s, userService)
	projectService := service.NewProjectService(daoInterface)
	v1.RegisterProjectServer(s, projectService)
	projectMemberService := service.NewProjectMemberService(daoInterface)
	v1.RegisterProjectMemberServer(s, projectMemberService)
	sprintService := service.NewSprintService(daoInterface)
	v1.RegisterSprintServer(s, sprintService)
	sprintMemberService := service.NewSprintMemberService(daoInterface)
	v1.RegisterSprintMemberServer(s, sprintMemberService)
	workItemService := service.NewWorkItemService(daoInterface)
	v1.RegisterWorkItemServer(s, workItemService)
	taskService := service.NewTaskService(daoInterface)
	v1.RegisterTaskServer(s, taskService)

	go func() {
		defer s.GracefulStop()
		<-ctx.Done()
	}()

	go func() error {
		log.L(ctx).Infof("grpc listen on: %s\n", address)
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
			LogLevel:                  logger.LogLevel(viper.GetInt("data.database.log-level")),
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
		MaxIdleConnections:    viper.GetInt("data.database.max-idle-connections"),
		MaxOpenConnections:    viper.GetInt("data.database.max-open-connections"),
		MaxConnectionLifeTime: time.Duration(viper.GetInt("data.database.max-connection-lifetime")) * time.Second,
		Logger:                newLogger,
	}
	letscrumDao, err := mysql.GetDao(&options)
	if err != nil {
		return nil, err
	}
	return letscrumDao, nil
}
