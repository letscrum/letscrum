package db

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

type Options struct {
	Host                  string
	Port                  string
	Username              string
	Password              string
	Database              string
	MaxIdleConnections    int
	MaxOpenConnections    int
	MaxConnectionLifeTime time.Duration
	Logger                logger.Interface
}

func NewGORM(opts *Options) (*gorm.DB, error) {
	dsn := GetConn(opts)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: opts.Logger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(opts.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(opts.MaxConnectionLifeTime)
	sqlDB.SetMaxIdleConns(opts.MaxIdleConnections)

	return db, nil
}

func NewDB(opts *Options) (*sql.DB, error) {
	dsn := GetConn(opts)
	return sql.Open(viper.GetString("data.database.driver"), dsn)
}

func GetConn(opts *Options) string {
	dsn := fmt.Sprintf(`%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=%t&loc=%s`,
		opts.Username,
		opts.Password,
		opts.Host,
		opts.Port,
		opts.Database,
		true,
		"Local")
	return dsn
}
