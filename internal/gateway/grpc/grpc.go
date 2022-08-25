package grpc

import (
	"context"
	golog "log"
	"net"
	"os"
	"time"

	"github.com/daocloud/skoala/api/hive/v1alpha1"
	"github.com/daocloud/skoala/app/hive/internal/dao"
	"github.com/daocloud/skoala/app/hive/internal/dao/mysql"
	svcv1alpha1 "github.com/daocloud/skoala/app/hive/internal/service/v1alpha1"
	"github.com/daocloud/skoala/app/pkg/db"
	"github.com/daocloud/skoala/pkg/log"

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

	// init dao,service
	var hiveDao dao.HiveDao
	if hiveDao, err = initDao(); err != nil {
		return err
	}
	hiveService := svcv1alpha1.NewHiveService(hiveDao)
	registrationService := svcv1alpha1.NewRegistrationService(hiveService)
	bookService := svcv1alpha1.NewBookService(hiveService)
	v1alpha1.RegisterHiveServer(s, hiveService)
	v1alpha1.RegisterRegistrationServer(s, registrationService)
	v1alpha1.RegisterBookServer(s, bookService)

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

func initDao() (dao.HiveDao, error) {
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
	hiveDao, err := mysql.GetHiveDaoOr(&options)
	if err != nil {
		return nil, err
	}
	return hiveDao, nil
}
