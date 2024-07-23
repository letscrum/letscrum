package db

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Options struct {
	Driver                string
	Host                  string
	Port                  string
	Username              string
	Password              string
	Database              string
	MaxIdleConnections    int
	MaxOpenConnections    int
	MaxConnectionLifeTime time.Duration
	AutoCreateAdmin       bool
	Logger                logger.Interface
}

func NewGORM(opts *Options) (*gorm.DB, error) {
	db, err := GetConn(opts)
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
	// test database connection and return *sql.DB
	db, err := GetConn(opts)
	if err != nil {
		return nil, err
	}
	return db.DB()
}

func GetDSN(opts *Options) string {
	var dsn string
	switch opts.Driver {
	case "mysql":
		dsn = fmt.Sprintf(`%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=%t&loc=%s`,
			opts.Username,
			opts.Password,
			opts.Host,
			opts.Port,
			opts.Database,
			true,
			"Local")
	case "postgres":
		dsn = fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`,
			opts.Host,
			opts.Port,
			opts.Username,
			opts.Password,
			opts.Database)
	}
	return dsn
}

func GetConn(opts *Options) (*gorm.DB, error) {
	dsn := GetDSN(opts)
	switch opts.Driver {
	case "mysql":
		return gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: opts.Logger,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
	case "postgres":
		return gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: opts.Logger,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
	}
	panic("never happen")
}
