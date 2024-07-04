package grpc

import (
	"context"
	golog "log"
	"net"
	"os"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	v1 "github.com/letscrum/letscrum/api/letscrum/v1"
	"github.com/letscrum/letscrum/internal/dao"
	"github.com/letscrum/letscrum/internal/dao/mysql"
	"github.com/letscrum/letscrum/internal/mid"
	servicev1 "github.com/letscrum/letscrum/internal/service/v1"
	"github.com/letscrum/letscrum/pkg/db"
	"github.com/letscrum/letscrum/pkg/log"
	"github.com/letscrum/letscrum/pkg/utils"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"gorm.io/gorm/logger"
)

func Run(ctx context.Context, opts utils.Options) error {
	//init grpc server and run
	l, err := net.Listen(opts.Network, opts.GRPCAddr)
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

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			selector.UnaryServerInterceptor(auth.UnaryServerInterceptor(mid.Auth), selector.MatchFunc(mid.AllButHealthZ)),
		),
		grpc.ChainStreamInterceptor(
			selector.StreamServerInterceptor(auth.StreamServerInterceptor(mid.Auth), selector.MatchFunc(mid.AllButHealthZ)),
		),
	)

	var daoInterface dao.Interface
	if daoInterface, err = initDao(); err != nil {
		return err
	}

	letscrumService := servicev1.NewLetscrumService(daoInterface)
	userService := servicev1.NewUserService(daoInterface)
	projectService := servicev1.NewProjectService(daoInterface)
	sprintService := servicev1.NewSprintService(daoInterface)
	workItemService := servicev1.NewWorkItemService(daoInterface)
	taskService := servicev1.NewTaskService(daoInterface)

	v1.RegisterLetscrumServer(s, letscrumService)
	v1.RegisterUserServer(s, userService)
	v1.RegisterProjectServer(s, projectService)
	v1.RegisterSprintServer(s, sprintService)
	v1.RegisterWorkItemServer(s, workItemService)
	v1.RegisterTaskServer(s, taskService)

	go func() {
		defer s.GracefulStop()
		<-ctx.Done()
	}()

	go func() error {
		log.L(ctx).Infof("grpc listen on: %s\n", opts.GRPCAddr)
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
		AutoCreateAdmin:       viper.GetBool("data.database.auto-create-admin"),
		MaxConnectionLifeTime: time.Duration(viper.GetInt("data.database.max-connection-lifetime")) * time.Second,
		Logger:                newLogger,
	}
	letscrumDao, err := mysql.GetDao(&options)
	if err != nil {
		return nil, err
	}
	return letscrumDao, nil
}
